package main

import (
	"fmt"
	"syscall"
	"time"
)

func showWindow(show bool) {
	getWin := syscall.NewLazyDLL(string([]byte{'k', 'e', 'r', 'n', 'e', 'l', '3', '2'})).NewProc("GetConsoleWindow")
	showWin := syscall.NewLazyDLL(string([]byte{'u', 's', 'e', 'r', '3', '2'})).NewProc("ShowWindow")
	hwnd, _, _ := getWin.Call()
	if getWin == nil {
		return
	}

	if show {
		var SW_RESTORE uintptr = 9
		showWin.Call(hwnd, SW_RESTORE)
		return
	}
	fmt.Println("[√] 完成,即将关闭窗口")
	time.Sleep(time.Second * 2)
	var SW_HIDE uintptr = 0
	showWin.Call(hwnd, SW_HIDE)
}
