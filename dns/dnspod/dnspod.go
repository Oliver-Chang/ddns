package dnspod

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const newRecordAPI = "https://dnsapi.cn/Record.Ddns"

// DNSpod DDNSpod
type DNSpod struct {
	domainID   string
	loginToken string
}

// New New
func New(loginID, domainID, token string) *DNSpod {
	return &DNSpod{
		loginToken: fmt.Sprintf("%s,%s", loginID, token),
		domainID:   domainID,
	}
}

// CreateRecord CreateRecord
func (d *DNSpod) CreateRecord(ipv6, domain string) error {
	var (
		// req  *http.Request
		err  error
		resp *http.Response
		body []byte
	)

	resp, err = http.PostForm(newRecordAPI, url.Values{
		"login_token":    {d.loginToken},
		"domain_id":      {d.domainID},
		"format":         {"json"},
		"sub_domain":     {strings.Split(domain, ".")[0]},
		"value":          {ipv6},
		"record_line_id": {"0"},
		"record_type":    {"AAAA"},
	})
	fmt.Println(domain)
	fmt.Println(d.domainID)
	// client := &http.Client{}
	// resp, err = client.Do(req)
	// if err != nil {
	// 	return err
	// }
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
