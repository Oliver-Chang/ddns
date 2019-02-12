package cloudfare

import (
	"fmt"

	"github.com/Oliver-Chang/ddns/util"
)

// CreateRecord CreateRecord
func CreateRecord(ip string, zoneID string) {
	if util.IsIPv4(ip) == true {
		fmt.Println("ipv4 address ", ip)
		return
	}
	url = "https://api.cloudflare.com/client/v4/zones/%s/dns_records"
}
