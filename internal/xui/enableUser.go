package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type XUIBulkAction struct {
	Emails []string `json:"emails"`
}

func (x *XUIClient) EnableUser(ctx context.Context, TgID int64) error {

	finalReqData := XUIBulkAction{
		Emails: []string{strconv.FormatInt(TgID, 10)},
	}

	finalBody, err := json.Marshal(finalReqData)
	if err != nil {
		return err
	}

	updateUrl := x.BaseURL + "/panel/api/clients/bulkEnable"

	newReq, err := http.NewRequestWithContext(ctx, "POST", updateUrl, bytes.NewBuffer(finalBody))
	if err != nil {
		return err
	}

	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Authorization", "Bearer "+x.ApiToken)
	resp, err := x.HTTPClient.Do(newReq)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add client, server returned status: %d", resp.StatusCode)
	}

	return nil

}
