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
	Enable     bool   `json:"enable"`
	TotalGB    int64  `json:"totalGb"`
	ExpiryTime int64  `json:"expiryTime"`
	TgID       int    `json:"tgId"`
	LimitIP    int    `json:"limitIp"`
	Reset      int    `json:"reset"`
	SubId      string `json:"subId"`
	Comment    string `json:"comment"`
	Security   string `json:"security"`
}

type XUIAddClient struct {
	Client     XUIClientSettings `json:"client"`
	InboundIDs []int64           `json:"inboundIds"`
}

type XUIError struct {
	SuccessStatus bool   `json:"success"`
	Message       string `json:"msg"`
}

func (x *XUIClient) AddUser(ctx context.Context, inboundID int64, uuid string, TgID int64, subID string) error {
	var xuiError XUIError
	clientSpec := XUIClientSettings{
		ID:         uuid,
		Email:      strconv.FormatInt(TgID, 10),
		LimitIP:    5,
		TotalGB:    0,
		ExpiryTime: 0,
		Enable:     true,
		SubId:      subID,
	}
	finalReqData := XUIAddClient{
		Client:     clientSpec,
		InboundIDs: []int64{inboundID},
	}

	finalBody, err := json.Marshal(finalReqData)
	if err != nil {
		return err
	}

	addClientUrl := x.BaseURL + "/panel/api/clients/add"

	req, err := http.NewRequestWithContext(ctx, "POST", addClientUrl, bytes.NewReader(finalBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+x.ApiToken)

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

	if !xuiError.SuccessStatus && strings.Contains(xuiError.Message, "email already in use") {
		err = x.EnableUser(ctx, TgID)
		if err != nil {
			return err
		}
	}
	return nil
}
