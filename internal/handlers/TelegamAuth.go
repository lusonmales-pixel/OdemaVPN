package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"svoy-vpn/internal/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	InitData string `json:"init_data"`
	RefCode  string `json:"ref_code"`
}

type tgUserData struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type telegramInitData struct {
	Hash      string
	AuthDate  int64
	DataCheck string
	User      tgUserData
}

type sendJWT struct {
	JWT string `json:"jwt"`
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

func parseTelegramInitData(initData string) (telegramInitData, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return telegramInitData{}, err
	}

	hash := values.Get("hash")
	if hash == "" {
		return telegramInitData{}, errors.New("missing hash")
	}

	authDateRaw := values.Get("auth_date")
	authDate, err := strconv.ParseInt(authDateRaw, 10, 64)
	if err != nil {
		return telegramInitData{}, errors.New("invalid auth_date")
	}
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	dataCheckParts := make([]string, 0, len(keys))
	for _, key := range keys {
		dataCheckParts = append(dataCheckParts, key+"="+values.Get(key))
	}

	userRaw := values.Get("user")
	if userRaw == "" {
		return telegramInitData{}, errors.New("missing user")
	}

	var user tgUserData
	if err := json.Unmarshal([]byte(userRaw), &user); err != nil {
		return telegramInitData{}, err
	}
	if user.ID == 0 {
		return telegramInitData{}, errors.New("invalid user id")
	}

	return telegramInitData{
		Hash:      hash,
		AuthDate:  authDate,
		DataCheck: strings.Join(dataCheckParts, "\n"),
		User:      user,
	}, nil
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

	var req AuthRequest
	err = json.Unmarshal(httpRequestBody, &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Error: "Failed to convert request body"})
		return
	}

	tgInitData, err := parseTelegramInitData(req.InitData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseError{Error: "Invalid Telegram initData"})
		return
	}

	if time.Now().Unix()-tgInitData.AuthDate > 86400 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseError{Error: "Telegram initData is expired"})
		return
	}

	calcedHash := createHash(tgInitData.DataCheck, e.BotToken)

	if tgInitData.Hash != calcedHash {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ResponseError{Error: "HashCode is not equals! Access denied!"})
		return
	}

	_, isNew, err := database.CreateUserIfNotExits(ctx, e.Conn, tgInitData.User.ID, tgInitData.User.Username)
	if req.RefCode != "" && isNew {
		_, err = database.GetUserByRefCode(ctx, e.Conn, req.RefCode, tgInitData.User.ID)
		if err != nil {
			log.Println("Failed to apply referal code:", err)
		}
	}

	JWTToken, err := e.generateJWT(tgInitData.User.ID)
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
