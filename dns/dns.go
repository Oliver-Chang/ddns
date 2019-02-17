package dns

import (
	"errors"
	"fmt"

	"github.com/Oliver-Chang/ddns/dns/cloudflare"
	"github.com/Oliver-Chang/ddns/util"
)

// DNSer DNSer
type DNSer interface {
	CreateRecord(ipv6, domain string) error
}

// DNSConfig DNSConfig
type DDNSConfig struct {
	DNS    string
	UID    string
	ZoneID string
	Token  string
	Domain string
}

type DDNS struct {
	config *DDNSConfig
}

func NewDDNS(conf *DDNSConfig) *DDNS {
	return &DDNS{
		config: conf,
	}
}

// CreateRecord CreateRecord
func (d *DDNS) CreateRecord() error {
	var (
		dns DNSer
		ip  string
		err error
		ok  bool
	)
	switch d.config.DNS {
	case "cloudflare":
		dns = cloudflare.New(d.config.UID, d.config.ZoneID, d.config.Token)
		// case "dnspod":
		// 	dns = dnspod.New(d.config.UID, d.config.ZoneID, d.config.Token)
	}
	for {
		if ip, ok = <-util.GetIP(); ok {
			if !util.IsIPv6(ip) {
				return errors.New("ip is not ipv6")
			}
			fmt.Println(ip)
			err = dns.CreateRecord(ip, d.config.Domain)
			// return errors.New("channel error")
			return err
		}

	}
}
