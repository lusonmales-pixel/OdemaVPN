package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"svoy-vpn/internal/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TgResponseData struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

type sendJWT struct {
	JWT string `json:"jwt"`
}

type RefCode struct {
	ReferralCode string `json:"ref_code"`
}

func (e *Env) generateJWT(tgID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": tgID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(e.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func createHash(dataCheck string, token string) string {
	mac := hmac.New(sha256.New, []byte("WebAppData"))
	mac.Write([]byte(token))
	secretKey := mac.Sum(nil)

	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheck))

	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return calculatedHash
}
func (e *Env) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to convert request body"})
		return
	}
	var JWTToken string
	var tgResp TgResponseData
	var ReferralCode RefCode

	err = json.Unmarshal(httpRequestBody, &tgResp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to convert request body"})
		return
	}
	err = json.Unmarshal(httpRequestBody, &ReferralCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to convert request body"})
		return
	}

	dataCheckString := fmt.Sprintf("auth_date=%d\nfirst_name=%s\nid=%d\nusername=%s", tgResp.AuthDate, tgResp.FirstName, tgResp.ID, tgResp.Username)

	CalcedHash := createHash(dataCheckString, e.BotToken)

	if tgResp.Hash != CalcedHash {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseError{Error: "HashCode is not equals! Access denied!"})
		return
	}
	_, isNew, err := database.CreateUserIfNotExits(ctx, e.Conn, tgResp.ID, tgResp.Username)
	if ReferralCode.ReferralCode != "" && isNew {
		_, err = database.GetUserByRefCode(ctx, e.Conn, ReferralCode.ReferralCode, tgResp.ID)
		if err != nil {
			log.Println("Failed to apply referal code:", err)
		}
	}

	JWTToken, err = e.generateJWT(tgResp.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Error in token generate"})
		return
	}

	jwtResponse := sendJWT{JWT: JWTToken}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	byteJWT, err := json.Marshal(jwtResponse)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to send JWT"})
		return
	}
	w.Write(byteJWT)
}
