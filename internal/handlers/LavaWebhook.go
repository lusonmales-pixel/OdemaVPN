package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"svoy-vpn/internal/database"
)

func (e *Env) LavaWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	xuiClient := e.XUIClient
	var lavaHook LavaWebhookStruct

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to read request body!"})
		return
	}

	err = json.Unmarshal(httpRequestBody, &lavaHook)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to convert body!"})
		return
	}

	lavaSignature := r.Header.Get("Signature")

	signature := CreateSignature(httpRequestBody, e.LavaSecret)

	if signature != lavaSignature {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseError{Error: "Signatures are not equal! Access denied."})
		return
	}

	if lavaHook.Status != "success" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseError{Error: "Operation status not success!"})
		return
	}

	err = database.Activate(ctx, e.Conn, lavaHook.CustomFields, lavaHook.Sum)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to activate subscription!"})
		return
	}
	uuid, err := database.GetUUID(ctx, e.Conn, lavaHook.CustomFields)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to get uuid"})
		return
	}

	err = xuiClient.AddUser(ctx, e.XUIInboundID, uuid, lavaHook.CustomFields)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to add user in xui panel"})
		return
	}

	w.WriteHeader(http.StatusOK)

}
