package iputil

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Oliver-Chang/ddns/util/logger"
)

// 基于 seeip.org 提供的 API 查询当前公共 IP 地址
const seeip = "https://ip.seeip.org"

// GetIP GetIP
func GetIP() (string, error) {
	var (
		ip   string
		rsp  *http.Response
		err  error
		body []byte
	)
	rsp, err = http.Get(seeip)
	if err != nil {
		logger.Logger.WithError(err).Error()
		return "", err
	}
	defer rsp.Body.Close()
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		// logger.Logger.WithError(err).Error()
		return "", err
	}
	ip = strings.Replace(string(body), "\n", "", -1)
	return ip, nil
}

// IsIPv6 IsIPv6
func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}
