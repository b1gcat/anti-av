package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/b1gcat/anti-av/utils"
)

type Mail struct {
	Sender   string
	Receiver string
	Interval int
	Body     string
	Attach   string
}

var (
	mail Mail
)

func main() {
	initialize()

	mail.attack()
}

func (m *Mail) attack() {
	vulners := m.getReceiver()

	/*
			swaks --to demo@qq.com --from test@163.com --ehlo 163.com --header-X-Mailer 163.com --header-Message-Id 1000
		--data aaa@sohu.com4.eml.txt

		--attach-type text/html --attach-body ~/Downloads/test.eml
	*/
	ehlo := strings.Split(m.Sender, "@")[1]
	xMail := ehlo
	for _, v := range vulners {
		cmd := make([]string, 0)
		cmd = append(cmd, "swaks",
			"--ehlo", m.Sender,
			"--to", v,
			"--from", m.Sender,
			"--header-X-Mailer", xMail,
			"--header-Message-Id", strconv.FormatInt(rand.Int63(), 10),
			"--attach-type text/html",
			"--attach-body", m.Body,
		)
		if m.Attach != "" {
			cmd = append(cmd, "--attach", m.Attach)
		}
		cmdJoin := strings.Join(cmd, " ")
		fmt.Println("[+] run:", cmdJoin)
		if err := utils.Cmd(cmdJoin); err != nil {
			fmt.Println("[-] Failed to run")
		}
		time.Sleep(time.Second * time.Duration(m.Interval))
	}
}

func initialize() {
	flag.StringVar(&mail.Sender, "s", "demo@qq.com", "发件人")
	flag.StringVar(&mail.Receiver, "r", "demo@163.com", "接收人: xxx@163.com,yyy@163.com 或 vuls.txt")
	flag.IntVar(&mail.Interval, "i", 10, "邮件发送间隔(秒)")
	flag.StringVar(&mail.Body, "f", "body.txt", "plain或html格式的邮件内容")
	flag.StringVar(&mail.Attach, "a", "", "附件")

	flag.Parse()

	mail.validate()
}

func (m *Mail) validate() {
	if _, err := exec.LookPath("swaks"); err != nil {
		fmt.Println("[-] 未安装 swaks: http://www.jetmore.org/john/code/swaks/")
		os.Exit(0)
	}
	if _, err := os.Stat(mail.Body); err != nil {
		fmt.Println("[-] eml文件不存在:", mail.Body)
		os.Exit(0)
	}

	if !verifyEmailFormat(mail.Sender) {
		fmt.Println("[-] 发件人格式错误")
		os.Exit(0)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("[+] 是否已经启用vpn,避免溯源?(回答:yes/no)")
		ans, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println("[-] ", err.Error())
			os.Exit(0)
		}
		if string(ans) == "yes" || string(ans) == "no" {
			break
		}
	}
}

func (m *Mail) getReceiver() []string {
	r := make([]string, 0)

	defer func() {
		fmt.Println("[+] vulners:", r)
	}()

	if _, err := os.Stat(m.Receiver); err != nil {
		rs := strings.Split(m.Receiver, ",")
		for _, v := range rs {
			if verifyEmailFormat(v) {
				r = append(r, v)
			} else {
				fmt.Println("[-] 忽略错误的邮件:", v)
			}
		}
		return r
	}

	f, err := os.Open(m.Receiver)
	if err != nil {
		fmt.Println("[-] 忽略错误的邮件接收者:", err.Error())
		os.Exit(0)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		v, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("[-] 解析接收人失败:", err.Error())
			os.Exit(0)
		}

		if verifyEmailFormat(string(v)) {
			r = append(r, string(v))
		} else {
			fmt.Println("[-] 忽略错误的邮件:", string(v))
		}
	}
	return r
}

func verifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func init() {
	rand.Seed(time.Now().Unix()) // unix 时间戳，秒
}
