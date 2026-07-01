package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"svoy-vpn/internal/xui"

	"github.com/jackc/pgx/v5"
)

type Env struct {
	Conn         *pgx.Conn
	XUIClient    *xui.XUIClient
	BotToken     string
	JwtSecret    []byte
	LavaShopID   string
	LavaSecret   string
	XUIInboundID int64
	ServerIp     string
	ServerPort   string
	ServerPBK    string
	ServerSNI    string
	ServerSID    string
}

func (e *Env) RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(ResponseError{Error: msg}); err != nil {
		log.Println("Error encoding error response:", err)
	}
}

func (e *Env) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Println("Error encoding JSON response:", err)
	}
}

type LavaWebhookStruct struct {
	Status       string  `json:"status"`
	Sum          float64 `json:"sum"`
	CustomFields int64   `json:"customFields"`
}

type LavaInvoiceRequest struct {
	Sum          float64 `json:"sum"`
	OrderId      string  `json:"orderId"`
	ShopId       string  `json:"shopId"`
	CustomFields int64   `json:"customFields"`
}

type LavaInvoiceResponse struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
}

type OrderRequest struct {
	TgID   int64   `json:"tg_id"`
	Amount float64 `json:"amount"`
}

type LinkResponse struct {
	URL string `json:"url"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseKey struct {
	VlessURL string `json:"vless_url"`
}
