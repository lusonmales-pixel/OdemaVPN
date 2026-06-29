package xui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (x *XUIClient) Connect(ctx context.Context) error {
	formData := url.Values{}

	formData.Set("username", x.Login)
	formData.Set("password", x.Password)

	loginURL := x.BaseURL + "/login"

	req, err := http.NewRequestWithContext(ctx, "POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := x.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Login Failed with status-code: %d", resp.StatusCode)
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	cookies := resp.Cookies()

	for _, cookie := range cookies {
		if cookie.Name == "session" {
			x.CookieSession = cookie.Value
			log.Println("Успешно авторизовались в 3X-UI, кука получена!")
			return nil
		}
	}

	log.Println("Status code 200, but there is not cookie!")
	return nil
}
