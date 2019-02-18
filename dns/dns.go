package dns

import (
	"github.com/Oliver-Chang/ddns/dns/cloudflare"
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
func (d *DDNS) CreateRecord(ipv6, domain string) error {
	var (
		dns DNSer
		err error
	)
	switch d.config.DNS {
	case "cloudflare":
		dns = cloudflare.New(d.config.UID, d.config.ZoneID, d.config.Token)
		// case "dnspod":
		// 	dns = dnspod.New(d.config.UID, d.config.ZoneID, d.config.Token)
	}

	err = dns.CreateRecord(ipv6, domain)
	if err != nil {
		return err
	}
	return nil
}
