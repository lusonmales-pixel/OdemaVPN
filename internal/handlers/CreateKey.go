package handlers

import (
	"fmt"
	"log"
	"net/http"
	"svoy-vpn/internal/database"
)

func (e *Env) CreateKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tgID, ok := r.Context().Value("TgID").(int64)
	if !ok {
		e.RespondWithError(w, 401, "Failed To get jwt")
		return
	}

	status, err := database.CheckStatus(ctx, e.Conn, tgID)
	if err != nil {
		log.Println("Database error checking status:", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Internal database error")
		return
	}

	if status != "active" {
		e.RespondWithError(w, http.StatusForbidden, "Subscribe is inactive")
		return
	}

	vlessUUID, err := database.GetUUID(ctx, e.Conn, tgID)
	if err != nil {
		log.Println("Database error getting UUID:", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Internal database error")
		return
	}

	vlessURL := fmt.Sprintf(
		"vless://%s@%s:%s?type=tcp&security=reality&pbk=%s&fp=chrome&sni=%s&sid=%s&spx=%%2F#OdemaVPN",
		vlessUUID,
		e.ServerIp,
		e.ServerPort,
		e.ServerPBK,
		e.ServerSNI,
		e.ServerSID,
	)

	e.RespondWithJSON(w, http.StatusOK, ResponseKey{VlessURL: vlessURL})
}
