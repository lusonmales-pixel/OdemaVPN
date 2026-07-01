package xui

import (
	"crypto/tls"
	"net/http"
	"time"
)

type XUIClient struct {
	BaseURL    string
	ApiToken   string
	HTTPClient *http.Client
}

func CreateClient(baseurl string, apiToken string) *XUIClient {
	return &XUIClient{
		BaseURL:  baseurl,
		ApiToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

}
