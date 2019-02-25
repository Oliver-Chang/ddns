package ddns

import (
	"context"
	"math/rand"
	"time"

	"github.com/Oliver-Chang/ddns/util/iputil"
	"github.com/Oliver-Chang/ddns/util/logger"
)

// DDNS ddns struct
type DDNS struct {
	ip     string
	config Config
}

// Deamon Deamon
func (d *DDNS) Deamon(cfg *Config) {
	// ipChan := make(chan string)
	// ctx, ctxCancel := context.WithTimeout(context.Background(), (10 * time.Second))
	for {
		// select {
		// case ipv6, ok := <-ipChan:
		// 	if ok {

		// 	}
		// }
	}
}

// FetchLatestIPv6 FetchLatestIPv6
func (d *DDNS) FetchLatestIPv6(ctx context.Context, ipSender chan string) {
	for i := 0; i < 5; i++ {
		ip, err := iputil.GetIP()
		if err != nil {
			t := time.Duration(1 >> uint32(rand.Intn(4)))
			logger.
				time.Sleep(t * time.Second)
			continue
		}
		if iputil.IsIPv6(ip) {
			break
		}
	}
}
