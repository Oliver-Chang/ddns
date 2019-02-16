package util

import (
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
)

// 基于 seeip.org 提供的 API 查询当前公共 IP 地址
const seeip = "https://ip.seeip.org"

// GetIP GetIP
func GetIP() <-chan string {
	var (
		ip chan string
	)
	ip = make(chan string)

	go func() {
		var (
			rsp  *http.Response
			err  error
			body []byte
		)
		rsp, err = http.Get(seeip)
		if err != nil {
			fmt.Println(err)
			close(ip)
			return
		}
		defer rsp.Body.Close()
		body, err = ioutil.ReadAll(rsp.Body)
		if err != nil {
			fmt.Println(err)
			close(ip)
			return
		}
		ip <- string(body)
	}()

	return ip
}

// IsIPv6 IsIPv6
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}
