package cloudflare

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Oliver-Chang/ddns/util/logger"

	"github.com/pkg/errors"
)

const apiURL = "https://api.cloudflare.com/client/v4"

// CloudFlare CloudFlare
type CloudFlare struct {
	AuthEmail string
	AuthKey   string
	baseURL   string
	headers   *http.Header
	client    *http.Client
}

// New New
func New(authEmail, authKey string) *CloudFlare {
	headers := make(http.Header)
	headers.Add("X-Auth-Email", authEmail)
	headers.Add("X-Auth-Key", authKey)
	headers.Add("Content-Type", "application/json")
	client := http.DefaultClient
	return &CloudFlare{
		AuthEmail: authEmail,
		AuthKey:   authKey,
		headers:   &headers,
		client:    client,
		baseURL:   apiURL,
	}

}

// Record Record
type Record struct {
	ID         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty"`
	Name       string      `json:"name,omitempty"`
	Content    string      `json:"content,omitempty"`
	Proxiable  bool        `json:"proxiable,omitempty"`
	Proxied    bool        `json:"proxied"`
	TTL        int         `json:"ttl,omitempty"`
	Locked     bool        `json:"locked,omitempty"`
	ZoneID     string      `json:"zone_id,omitempty"`
	ZoneName   string      `json:"zone_name,omitempty"`
	CreatedOn  time.Time   `json:"created_on,omitempty"`
	ModifiedOn time.Time   `json:"modified_on,omitempty"`
	Data       interface{} `json:"data,omitempty"` // data returned by: SRV, LOC
	Meta       interface{} `json:"meta,omitempty"`
	Priority   int         `json:"priority"`
}

// CreateRecord CreateRecord
func (c *CloudFlare) CreateRecord(zoneID string, record Record) error {
	var (
		res []byte
		err error
	)
	url := c.baseURL + "/zones/" + zoneID + "/dns_records"
	res, err = c.makeRequest("POST", url, record)
	if err != nil {
		return err
	}
	logger.Logger.WithField("res", res)
	return nil
}

func (c *CloudFlare) makeRequest(method, url string, params interface{}) ([]byte, error) {
	var (
		err       error
		jsonBtyes []byte
		reqBody   io.Reader
	)
	if params != nil {
		if paramsBytes, ok := params.([]byte); ok {
			jsonBtyes = paramsBytes
		} else {
			jsonBtyes, err = json.Marshal(params)
			if err != nil {
				return nil, errors.Wrap(err, "Can't marshal params to bytes")
			}
		}
	} else {
		jsonBtyes = nil
	}
	reqBody = bytes.NewReader(jsonBtyes)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "Create Request failed")
	}
	req.Header = *c.headers
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read respone body failed")
	}
	return respBody, nil
}
