package xui

import (
	"net/http"
	"time"
)

type XUIClient struct {
	BaseURL       string
	Login         string
	Password      string
	CookieSession string
	HTTPClient    *http.Client
}

func CreateClient(baseurl string, login string, password string) *XUIClient {
	return &XUIClient{
		BaseURL:  baseurl,
		Login:    login,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

}
