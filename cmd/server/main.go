package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"svoy-vpn/internal/database"
	"svoy-vpn/internal/handlers"
	"svoy-vpn/internal/xui"
	"time"

	"github.com/joho/godotenv"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	xuiURL := os.Getenv("XUI_BASE_URL")
	xuiApiToken := os.Getenv("XUI_API_TOKEN")
	inboundIDStr := os.Getenv("XUI_INBOUND_ID")

	inboundID, err := strconv.ParseInt(inboundIDStr, 10, 64)
	if err != nil {
		log.Fatalln("Invalid XUI_INBOUND_ID in config:", err)
	}

	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}
	defer conn.Close(ctx)

	if err := database.InitTable(ctx, conn); err != nil {
		log.Fatalln("Failed to create tables:", err)
	}

	xuiClient := xui.CreateClient(xuiURL, xuiApiToken)

	env := &handlers.Env{
		Conn:         conn,
		XUIClient:    xuiClient,
		BotToken:     os.Getenv("BOT_TOKEN"),
		JwtSecret:    []byte(os.Getenv("JWT_SECRET")),
		LavaShopID:   os.Getenv("LAVA_SHOP_ID"),
		LavaSecret:   os.Getenv("LAVA_SECRET_KEY"),
		XUIInboundID: inboundID,
		ServerIp:     os.Getenv("ServerIP"),
		ServerPort:   os.Getenv("ServerPort"),
		ServerPBK:    os.Getenv("ServerPBK"),
		ServerSNI:    os.Getenv("ServerSNI"),
		ServerSID:    os.Getenv("ServerSID"),
		SubURL:       os.Getenv("SUB_URL"),
	}

	http.Handle("/api/payment/create", env.ValidateJWT(http.HandlerFunc(env.CreateOrder)))
	http.HandleFunc("/api/v1/payments/lava/webhook", env.LavaWebhook)
	http.Handle("/api/user/config", env.ValidateJWT(http.HandlerFunc(env.CreateKey)))
	http.Handle("/api/user/subscription", env.ValidateJWT(http.HandlerFunc(env.GetSubscription)))
	http.HandleFunc("/api/auth", env.Auth)
	http.Handle("/api/referral/getCode", env.ValidateJWT(http.HandlerFunc(env.GetReferralCode)))

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			ctx := context.Background()
			log.Println("Running subscription expiration worker...")

			expiredUsers, err := database.Expire(ctx, env.Conn)
			if err != nil {
				log.Println("Error in expiration worker database check:", err)
				continue
			}

			for _, user := range expiredUsers {
				err := env.XUIClient.DisableUser(ctx, user.TgId)
				if err != nil {
					log.Printf("Failed to disable user %d on 3X-UI panel: %v", user.TgId, err)
				} else {
					log.Printf("Successfully disabled expired user %d in 3X-UI panel", user.TgId)
				}
			}
		}
	}()

	if err := http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}
