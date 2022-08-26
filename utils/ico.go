package utils

import (
	"os"
	"path/filepath"

	"github.com/josephspurrier/goversioninfo"
	"github.com/sirupsen/logrus"
)

func CreateIcoPropertity(arch string, output string) {
	vi := &goversioninfo.VersionInfo{}

	vi.StringFileInfo.CompanyName = "白帽子集团"
	vi.StringFileInfo.InternalName = "主机基线检查工具箱"
	vi.StringFileInfo.FileDescription = "主机基线检查工具箱"
	vi.StringFileInfo.FileVersion = "17.0.10001.10000"
	vi.StringFileInfo.LegalCopyright = "© security. All rights reserved."
	vi.StringFileInfo.OriginalFilename = "security box"
	vi.FixedFileInfo.ProductVersion.Patch = 10001
	vi.FixedFileInfo.ProductVersion.Major = 16
	vi.FixedFileInfo.ProductVersion.Minor = 0
	vi.StringFileInfo.ProductName = "主机基线检查工具"
	vi.StringFileInfo.ProductVersion = "17.0.10001.10000"
	vi.FixedFileInfo.FileVersion.Major = 16
	vi.FixedFileInfo.FileVersion.Minor = 0
	vi.FixedFileInfo.FileVersion.Patch = 10001
	vi.FixedFileInfo.FileVersion.Build = 10000
	vi.IconPath = filepath.Join(output, "resource", "/", "logo.ico")
	vi.Build()
	vi.Walk()

	fileout := "resource_windows.syso"
	if err := vi.WriteSyso(fileout, arch); err != nil {
		logrus.Error("[-] Error writing syso: ", err.Error())
		os.Exit(3)
	}

	logrus.Info("[+] 生成图标资源文件:", arch)
}
