package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
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

type XUIError struct {
	SuccessStatus bool   `json:"success"`
	Message       string `json:"msg"`
}

func (x *XUIClient) AddUser(ctx context.Context, inboundID int64, uuid string, TgID int64) error {
	var xuiError XUIError
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
	req.AddCookie(&http.Cookie{Name: "3x-ui", Value: x.CookieSession})

	resp, err := x.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &xuiError)
	if err != nil {
		return err
	}

	if !xuiError.SuccessStatus && strings.Contains(xuiError.Message, "Duplicate") {
		err = x.EnableUser(ctx, inboundID, uuid, TgID)
		if err != nil {
			return err
		}
	}
	return nil
}
