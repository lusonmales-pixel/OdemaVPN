package handlers

import (
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

	subID, err := database.GetSubID(ctx, e.Conn, tgID)
	if err != nil {
		log.Println("Database error getting SubID:", err)
		e.RespondWithError(w, http.StatusInternalServerError, "Internal database error")
		return
	}

	SubURL := e.SubURL + subID

	e.RespondWithJSON(w, http.StatusOK, ResponseKey{SubURL: SubURL})
}
