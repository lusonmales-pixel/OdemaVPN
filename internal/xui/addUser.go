package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type XUIClientSettings struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	LimitIP    int    `json:"limitIp"`
	TotalGB    int64  `json:"totalGb"`
	ExpiryTime int64  `json:"expiryTime"`
	Enable     bool   `json:"enable"`
	Flow       string `json:"flow"`
}

type XUIAddClient struct {
	InboundID int64  `json:"id"`
	Settings  string `json:"settings"`
}

type XUIClientsFields struct {
	Clients []XUIClientSettings `json:"clients"`
}

func (x *XUIClient) AddUser(ctx context.Context, inboundID int64, uuid string, TgID int64) error {
	clientSpec := XUIClientSettings{
		ID:         uuid,
		Email:      strconv.FormatInt(TgID, 10),
		LimitIP:    0,
		TotalGB:    0,
		ExpiryTime: 0,
		Enable:     true,
		Flow:       "",
	}

	var wrap XUIClientsFields

	wrap.Clients = []XUIClientSettings{clientSpec}

	settingsBuff, err := json.Marshal(wrap)
	if err != nil {
		return err
	}

	finalReqData := XUIAddClient{
		InboundID: inboundID,
		Settings:  string(settingsBuff),
	}

	finalBody, err := json.Marshal(finalReqData)
	if err != nil {
		return err
	}

	addClientUrl := x.BaseURL + "/panel/api/inbounds/addClient"

	req, err := http.NewRequestWithContext(ctx, "POST", addClientUrl, bytes.NewReader(finalBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session", Value: x.CookieSession})

	resp, err := x.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add client, server returned status: %d", resp.StatusCode)
	}

	return nil
}
