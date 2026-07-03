package handlers

import (
	"net/http"
	"svoy-vpn/internal/database"
	"time"
)

type SubscriptionResponse struct {
	Active    bool       `json:"active"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (e *Env) GetSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tgID, ok := ctx.Value("TgID").(int64)
	if !ok {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to get TgID from context")
		return
	}

	subscription, err := database.GetSubscription(ctx, e.Conn, tgID)
	if err != nil {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to get subscription")
		return
	}

	e.RespondWithJSON(w, http.StatusOK, SubscriptionResponse{
		Active:    subscription.Status == "active",
		ExpiresAt: subscription.ExpiresAt,
	})
}
