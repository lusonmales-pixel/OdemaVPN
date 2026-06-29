package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func CreateSignature(byteStruct []byte, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(byteStruct)
	return hex.EncodeToString(h.Sum(nil))
}

func (e *Env) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req OrderRequest

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		e.RespondWithError(w, http.StatusBadRequest, "Failed to read body")
		return
	}

	err = json.Unmarshal(httpRequestBody, &req)
	if err != nil {
		e.RespondWithError(w, http.StatusBadRequest, "Failed to parse JSON request")
		return
	}

	invoiceReq := LavaInvoiceRequest{
		Sum:          req.Amount,
		OrderId:      uuid.New().String(),
		ShopId:       e.LavaShopID,
		CustomFields: req.TgID,
	}

	invoiceReqByte, err := json.Marshal(invoiceReq)
	if err != nil {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to build invoice")
		return
	}

	signature := CreateSignature(invoiceReqByte, e.LavaSecret)
	lavaUrl := "https://api.lava.ru/invoice/create"

	lavaRequest, err := http.NewRequestWithContext(ctx, "POST", lavaUrl, bytes.NewReader(invoiceReqByte))
	if err != nil {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to create request!")
		return
	}
	lavaRequest.Header.Set("Content-Type", "application/json")
	lavaRequest.Header.Set("Accept", "application/json")
	lavaRequest.Header.Set("Signature", signature)

	client := http.Client{}
	response, err := client.Do(lavaRequest)
	if err != nil {
		e.RespondWithError(w, http.StatusBadGateway, "Failed to send request to Lava")
		return
	}

	defer response.Body.Close()

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to read response body!")
		return
	}

	var lavaResp LavaInvoiceResponse
	err = json.Unmarshal(respBytes, &lavaResp)
	if err != nil {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to convert response!")
		return
	}

	e.RespondWithJSON(w, http.StatusOK, LinkResponse{URL: lavaResp.Data.URL})
}
