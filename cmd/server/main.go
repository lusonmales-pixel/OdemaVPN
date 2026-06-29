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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	xuiURL := os.Getenv("XUI_BASE_URL")
	xuiLogin := os.Getenv("XUI_LOGIN")
	xuiPassword := os.Getenv("XUI_PASSWORD")
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

	xuiClient := xui.CreateClient(xuiURL, xuiLogin, xuiPassword)
	if err := xuiClient.Connect(ctx); err != nil {
		log.Println("WARNING: Failed to login into 3X-UI panel:", err)
	}

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
	}

	http.HandleFunc("/api/payment/create", env.CreateOrder)
	http.HandleFunc("/api/v1/payments/lava/webhook", env.LavaWebhook)
	http.Handle("/api/user/config", env.ValidateJWT(http.HandlerFunc(env.CreateKey)))
	http.HandleFunc("/api/auth", env.Auth)
	http.HandleFunc("/test/add-user", env.TestAddUser)

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
				err := env.XUIClient.DisableUser(ctx, env.XUIInboundID, user.Vless_uuid, user.TgId)
				if err != nil {
					log.Printf("Failed to disable user %d on 3X-UI panel: %v", user.TgId, err)
				} else {
					log.Printf("Successfully disabled expired user %d in 3X-UI panel", user.TgId)
				}
			}
		}
	}()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}
