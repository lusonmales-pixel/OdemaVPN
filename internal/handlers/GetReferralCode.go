package handlers

import (
	"net/http"
	"svoy-vpn/internal/database"
)

type ReferralResponse struct {
	ReferralCode string `json:"ref_code"`
}

func (e *Env) GetReferralCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	TgID, ok := ctx.Value("TgID").(int64)
	if !ok {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to get TgID from context")
		return
	}

	ReferralCodeFinal, err := database.GetReferralCode(ctx, e.Conn, TgID)
	if err != nil {
		e.RespondWithError(w, http.StatusInternalServerError, "Failed to get referral code")
		return
	}

	e.RespondWithJSON(w, http.StatusOK, ReferralResponse{ReferralCode: ReferralCodeFinal})

}
