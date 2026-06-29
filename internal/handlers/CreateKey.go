package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"svoy-vpn/internal/database"
)

func (e *Env) CreateKey(w http.ResponseWriter, r *http.Request) {
	var userid UserID
	ctx := r.Context()

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		e.RespondWithError(w, http.StatusBadRequest, "Failed to read body")
		return
	}

	err = json.Unmarshal(httpRequestBody, &userid)
	if err != nil {
		e.RespondWithError(w, http.StatusBadRequest, "Failed to convert request body")
		return
	}

	status, err := database.CheckStatus(ctx, e.Conn, userid.ID)
	if err != nil {
		log.Println("Database error checking status:", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Internal database error")
		return
	}

	if status != "active" {
		e.RespondWithError(w, http.StatusForbidden, "Subscribe is inactive")
		return
	}

	vlessUUID, err := database.GetUUID(ctx, e.Conn, userid.ID)
	if err != nil {
		log.Println("Database error getting UUID:", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Internal database error")
		return
	}

	vlessURL := fmt.Sprintf(
		"vless://%s@2.26.105.226:443?type=tcp&security=reality&pbk=Z1_vIn2G4v97Oisw7SgC6Qh9rW_wF841XWv265U_I00&fp=chrome&sni=www.nvidia.com&sid=0410427b68634839&spx=%%2F#OdemaVPN",
		vlessUUID,
	)

	e.RespondWithJSON(w, http.StatusOK, ResponseKey{VlessURL: vlessURL})
}
