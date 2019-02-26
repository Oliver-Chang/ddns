package daemon

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"go.uber.org/zap"

	"github.com/Oliver-Chang/ddns/utils/logger"

	"github.com/Oliver-Chang/ddns/utils/iputil"
)

// DDNS ddns struct
type DDNS struct {
	ip     string
	config *Config
}

// New DDNS
func New(cfg *Config) *DDNS {
	return &DDNS{
		config: cfg,
	}
}

// Daemon Daemon
func (d *DDNS) Daemon() {
	ipChan := make(chan string)
	schedTime := 8 * time.Minute
	// will pre 8m fetch ip address
	go d.fetchLatestIPv6(ipChan, schedTime)
	for {
		select {
		case ipv6, ok := <-ipChan:
			if ok {
				logger.Logger.Info("will update dns record", zap.String("ipv6", ipv6))
				err := d.updateRecord(ipv6)
				if err != nil {
					logger.Logger.Error("update record failed", zap.NamedError("update_error", err))
				}
				logger.Logger.Info("update dns record success")
				d.ip = ipv6
			}
		}
	}
}

func (d *DDNS) updateRecord(ipv6 string) error {
	var (
		zoneID    string
		api       *cloudflare.API
		records   []cloudflare.DNSRecord
		subdomain string
		err       error
	)
	subdomain = d.config.SubDomain
	api, err = cloudflare.New(d.config.AuthKey, d.config.AuthEmail)
	if err != nil {
		return err
	}
	// subdomain is aaa.bbb.ccc  domain is bbb.ccc
	domain := strings.SplitN(subdomain, ".", 2)[1]
	zoneID, err = api.ZoneIDByName(domain)
	if err != nil {
		return err
	}
	records, err = api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: subdomain})
	if err != nil {
		return err
	}

	// recordLen > 1 数据有错误，请检查 cloudflare 的配置
	if recordLen := len(records); recordLen > 1 {
		return errors.New("subname record len > 1, please check it")
		// recordLen == 0 没有对应 subdomain record， 创建一个。
	} else if recordLen == 0 {
		logger.Logger.Info("no this subdomain record, it will create it")
		resp, err := api.CreateDNSRecord(zoneID,
			cloudflare.DNSRecord{
				Type:     "AAAA",
				Name:     subdomain,
				Content:  ipv6,
				Proxied:  true,
				Priority: 10,
			})
		if err != nil {
			return err
		}
		logger.Logger.Info("create dns record success", zap.String("create_resp", fmt.Sprintf("%+v", resp)))
		return nil
		// finally update record
	} else {
		recordID := records[0].ID
		logger.Logger.Info("update subdomain record")
		err := api.UpdateDNSRecord(zoneID, recordID,
			cloudflare.DNSRecord{
				Content: ipv6,
			})
		if err != nil {
			return err
		}
		return nil
	}

}

// fetchLatestIPv6 每隔一段时间获取最新且更改过的 ipv6 地址
func (d *DDNS) fetchLatestIPv6(ipSender chan string, sleepTime time.Duration) {
	var (
		ip    string
		oldIP string
		err   error
	)
getip:
	oldIP = d.ip
	for i := 0; i < 3; i++ {
		ip, err = iputil.GetIP()
		if err != nil {
			t := time.Duration(1<<uint32(rand.Intn(4))) * time.Second
			logger.Logger.Error("get ip address error", zap.NamedError("get_ip", err), zap.Duration("sleep_time", t))
			time.Sleep(t)
			continue
		}
	}
	if !iputil.IsIPv6(ip) {
		t := time.Duration(1<<uint32(rand.Intn(5))) * time.Second
		logger.Logger.Error("ip address not ipv6", zap.String("ip", ip), zap.Duration("sleep_time", t))
		time.Sleep(t)
		goto getip
	}
	if ip != oldIP {
		ipSender <- ip

	} else {
		logger.Logger.Info("ip address is not change", zap.String("ip", ip), zap.String("old_ip", oldIP))
	}
	time.Sleep(sleepTime)
	goto getip
}
