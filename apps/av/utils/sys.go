package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Cmd(cmd string) error {
	var sh *exec.Cmd
	if runtime.GOOS == "windows" {
		sh = exec.Command("cmd", "/C", cmd)
	} else {
		sh = exec.Command("sh", "-c", cmd)
	}
	sh.Stdin = os.Stdin
	sh.Stdout = os.Stdout
	sh.Stderr = os.Stderr
	if err := sh.Run(); err != nil {
		return fmt.Errorf("%v:%v", cmd, err.Error())
	}
	return nil
}
