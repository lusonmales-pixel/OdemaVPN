package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (x *XUIClient) EnableUser(ctx context.Context, inboundID int64, uuid string, TgID int64) error {
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

	updateUrl := x.BaseURL + "/panel/api/inbounds/updateClient/" + uuid

	newReq, err := http.NewRequestWithContext(ctx, "POST", updateUrl, bytes.NewBuffer(finalBody))
	if err != nil {
		return err
	}

	newReq.Header.Set("Content-Type", "application/json")
	newReq.AddCookie(&http.Cookie{Name: "3x-ui", Value: x.CookieSession})

	resp, err := x.HTTPClient.Do(newReq)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add client, server returned status: %d", resp.StatusCode)
	}

	return nil

}
