package xui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (x *XUIClient) DisableUser(ctx context.Context, tgID int64) error {

	finalReqData := XUIBulkAction{
		Emails: []string{strconv.FormatInt(tgID, 10)},
	}

	finalReqDataByte, err := json.Marshal(finalReqData)
	if err != nil {
		return err
	}

	disableClientURL := x.BaseURL + "/panel/api/clients/bulkDisable"

	req, err := http.NewRequestWithContext(ctx, "POST", disableClientURL, bytes.NewReader(finalReqDataByte))
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add client, server returned status: %d", resp.StatusCode)
	}

	return nil
}
