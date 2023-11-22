package main

type config struct {
	//生成shellcode
	loader         string //构造shellcode的方式
	paylaod        string //payload文件
	os             string //windows,linux
	valid          string //签名域名
	hostObfuscator string //混淆远程加载shellcode的地址，干扰蓝队告警日志研判
	arch           string //编译支持的架构,默认x86
	mode           string //模式MORNAL,INJECT
	noIcon         bool   //禁止签名
	output         string
	//加密 shellcode
	genRaw bool
}
