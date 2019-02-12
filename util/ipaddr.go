package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
)

// Resp Resp
type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		IP string `json:"ip"`
	} `json:"data"`
}

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
			resp Resp
		)
		rsp, err = http.Get("http://tool.oliverch.com/ip")
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
		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Println(err)
			close(ip)
			return
		}
		ip <- resp.Data.IP
	}()
	return ip
}

// IsIPv4 IsIPv4
func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

// IsIPv6 IsIPv6
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}
