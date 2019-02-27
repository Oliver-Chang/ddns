package daemon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Oliver-Chang/ddns/utils/iputil"

	cloudflare "github.com/cloudflare/cloudflare-go"

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

	for {
		select {
		case ipv6, ok := <-d.ipChan:
			if ok {
				logger.Logger.Info().Str("ipv6", ipv6).Msg("will update dns record")
				err := d.updateRecord(ipv6)
				if err != nil {
					logger.Logger.Error().Err(err).Msg("update record failed")
				}
				logger.Logger.Info().Msg("update dns record success")
				d.ip = ipv6
			}
		case <-ticker.C:
			if err := d.FetchIPv6(); err != nil {
				logger.Logger.Error().Err(err).Msg("fetch ipv6 address failed")
			}
		}
	}
}

// FetchIPv6  FetchIPv6
func (d *DDNS) FetchIPv6() error {
	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Minute)
	go func(ctx context.Context) {
		ipv6 := iputil.GetIPv6()
		logger.Logger.Info().Str("ipv6", ipv6).Msg("success fetch ipv6 address")
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
		logger.Logger.Info().Msg("create subdomain record")
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
		logger.Logger.Info().Str("create_record", fmt.Sprintf("%+v", resp)).Msg("create dns record success")
		return nil
		// finally have a record with subdomain update record content
	} else {
		record := records[0]
		if recordContent := record.Content; recordContent == ipv6 {
			logger.Logger.Info().Str("ipv6", ipv6).Str("content", recordContent).Msg("record content have not change, should not update")
			return nil
		}
		recordID := record.ID
		record.Content = ipv6
		err := api.UpdateDNSRecord(zoneID, recordID, record)
		if err != nil {
			return err
		}
		logger.Logger.Info().Msg("update subdomain record")
		return nil
	}

}
