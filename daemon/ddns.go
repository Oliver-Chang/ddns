package daemon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Oliver-Chang/ddns/utils/iputil"

	cloudflare "github.com/cloudflare/cloudflare-go"

	"go.uber.org/zap"

	"github.com/Oliver-Chang/ddns/utils/logger"
)

// DDNS ddns struct
type DDNS struct {
	ip     string
	config *Config
	ipChan chan string
	// ticker
}

// New DDNS
func New(cfg *Config) *DDNS {
	ipChan := make(chan string)
	return &DDNS{
		config: cfg,
		ipChan: ipChan,
	}
}

// Daemon Daemon
func (d *DDNS) Daemon() {
	// schedTime := 8 * time.Minute
	// will pre 8m fetch ip address
	ticker := time.NewTicker(8 * time.Minute)
	if err := d.FetchIPv6(); err != nil {
		logger.Logger.Error("first fetch ipv6 failed", zap.NamedError("fetch_error", err))
	}
	for {
		select {
		case ipv6, ok := <-d.ipChan:
			if ok {
				logger.Logger.Info("will update dns record", zap.String("ipv6", ipv6))
				err := d.updateRecord(ipv6)
				if err != nil {
					logger.Logger.Error("update record failed", zap.NamedError("update_error", err))
				}
				logger.Logger.Info("update dns record success")
				d.ip = ipv6
			}
		case <-ticker.C:
			if err := d.FetchIPv6(); err != nil {
				logger.Logger.Error("fetch ipv6 address failed", zap.NamedError("fetch_error", err))
			}
		}
	}
}

// FetchIPv6  FetchIPv6
func (d *DDNS) FetchIPv6() error {
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Minute)
	go func(ctx context.Context) {
		ipv6 := iputil.GetIPv6()
		logger.Logger.Info("success get ipv6 address", zap.String("ipv6", ipv6))
		d.ipChan <- ipv6
	}(ctx)
	defer cancle()
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
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
	records, err = api.DNSRecords(zoneID, cloudflare.DNSRecord{Name: subdomain, Type: "AAAA"})
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
		// finally have a record with subdomain update record content
	} else {
		record := records[0]
		if recordContent := record.Content; recordContent == ipv6 {
			logger.Logger.Info("record content have not change, should not update", zap.String("ipv6", ipv6), zap.String("content", recordContent))
			return nil
		}
		recordID := record.ID
		record.Content = ipv6
		logger.Logger.Info("update subdomain record")
		err := api.UpdateDNSRecord(zoneID, recordID, record)
		if err != nil {
			return err
		}
		return nil
	}

}
