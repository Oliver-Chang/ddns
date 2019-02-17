package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const newRecordAPI = "https://api.cloudflare.com/client/v4/zones/%s/dns_records"

// CloudFlare CloudFlare
type CloudFlare struct {
	authID  string
	zoneID  string
	authKey string
}

// New New
func New(authID, zoneID, authKey string) *CloudFlare {
	return &CloudFlare{
		authID:  authID,
		zoneID:  zoneID,
		authKey: authKey,
	}
}

// CreateRecord CreateRecord
func (c *CloudFlare) CreateRecord(ipv6, domain string) error {
	var (
		req  *http.Request
		resp *http.Response
		data map[string]string
		err  error
		url  string
		body []byte
	)
	data = map[string]string{
		"type":    "AAAA",
		"name":    domain,
		"content": ipv6,
	}
	bytesJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(bytesJSON))
	url = fmt.Sprintf(newRecordAPI, c.zoneID)
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(bytesJSON))
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("X-Auth-Email", c.authID)
	req.Header.Add("X-Auth-Key", c.authKey)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
