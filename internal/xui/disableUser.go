package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (x *XUIClient) DisableUser(ctx context.Context, inboundID int64, uuid string, tgID int64) error {
	clientSpec := XUIClientSettings{
		ID:         uuid,
		Email:      strconv.FormatInt(tgID, 10),
		LimitIP:    0,
		TotalGB:    0,
		ExpiryTime: 0,
		Enable:     false,
		Flow:       "",
	}

	var wrap XUIClientsFields

	wrap.Clients = []XUIClientSettings{clientSpec}

	wrapByte, err := json.Marshal(wrap)
	if err != nil {
		return err
	}

	finalReqData := XUIAddClient{
		InboundID: inboundID,
		Settings:  string(wrapByte),
	}

	finalReqDataByte, err := json.Marshal(finalReqData)

	disableClientURL := x.BaseURL + "/panel/api/inbounds/updateClient/" + uuid

	req, err := http.NewRequestWithContext(ctx, "POST", disableClientURL, bytes.NewReader(finalReqDataByte))

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
