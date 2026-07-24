<div align="center">

# OdemaVPN

### Lightweight Go backend for automated VPN subscription management via Telegram and 3X-UI.

<br>

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-pgx-4169E1?style=for-the-badge&logo=postgresql)
![REST API](https://img.shields.io/badge/API-REST-009688?style=for-the-badge)
![Telegram](https://img.shields.io/badge/Telegram-Login%20%26%20Bot-26A5E4?style=for-the-badge&logo=telegram)
![VLESS](https://img.shields.io/badge/VPN-VLESS-6A1B9A?style=for-the-badge)
![3X--UI](https://img.shields.io/badge/3X--UI-Integrated-FF9800?style=for-the-badge)
![JWT](https://img.shields.io/badge/Auth-JWT-black?style=for-the-badge)
![Docker](https://img.shields.io/badge/Docker-Not%20Required-2496ED?style=for-the-badge&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

</div>

---

# Overview

OdemaVPN is a lightweight backend service that automates the complete lifecycle of VPN subscriptions for Telegram users.

The project integrates **Telegram Login**, **3X-UI**, **Xray/VLESS**, **Lava Payments**, and **PostgreSQL** into a single stateless REST API responsible for user registration, payment processing, VPN provisioning, subscription renewal, referral rewards, and automatic expiration handling.

The architecture intentionally avoids unnecessary frameworks and ORMs, relying on Go's standard library, direct SQL, and a small dependency footprint to maximize maintainability, performance, and deployment simplicity.

---

# Architecture

```mermaid
graph TD

A[Telegram Mini App] --> B[REST API]

B --> C[JWT Authentication]

B --> D[PostgreSQL]

B --> E[Lava Payment API]

E --> F[Payment Webhook]

F --> B

B --> G[3X-UI API]

G --> H[Xray / VLESS]

I[Daily Expiration Worker] --> D

I --> G
```

---

## Registration Flow

```mermaid
sequenceDiagram

participant User

participant Telegram

participant Backend

participant PostgreSQL

participant 3X-UI

User->>Telegram: Login

Telegram-->>Backend: Login Data

Backend->>Backend: Validate Signature

Backend->>PostgreSQL: Create User

Backend->>Backend: Generate UUID & Subscription ID

Backend->>3X-UI: Add VPN Client

Backend-->>User: JWT Token
```

---

# Key Features

- 🔐 **Telegram Login authentication** with signature verification.
- 🌐 **Automatic VPN provisioning** via the 3X-UI REST API.
- 📦 **Personal VLESS subscription generation** for every user.
- 💳 **Lava payment integration** with secure webhook validation.
- 🔄 **Automatic subscription renewal** without losing remaining paid time.
- ⏰ **Daily expiration worker** that disables expired VPN accounts.
- 🎁 **Referral reward system** with automatic bonus activation.
- ⚡ **Stateless JWT authentication** without server-side sessions.
- 🗄️ **Direct SQL using pgx** instead of ORM for predictable performance.
- 🧩 **Minimal dependency footprint** built almost entirely on Go's standard library.

---

# Data Flow

```text
Telegram Login
      │
      ▼
JWT Authentication
      │
      ▼
Create User
      │
      ▼
Generate UUID + Subscription ID
      │
      ▼
Create VPN Client (3X-UI)
      │
      ▼
Issue Subscription URL
      │
      ▼
User Payment
      │
      ▼
Lava Webhook
      │
      ▼
Activate Subscription
      │
      ▼
Enable VPN Client
      │
      ▼
Daily Worker
      │
      ▼
Disable Expired Clients
```

---

# Technology Stack

| Component | Technology |
|------------|------------|
| Language | Go 1.26+ |
| HTTP | net/http |
| Database | PostgreSQL |
| Driver | pgx/v5 |
| Authentication | JWT |
| VPN | VLESS |
| VPN Panel | 3X-UI |
| Payments | Lava API |
| UUID | google/uuid |
| Configuration | godotenv |

---

# Quick Start

## Clone

```bash
git clone https://github.com/username/OdemaVPN.git

cd OdemaVPN
```

---

## Install dependencies

```bash
go mod download
```

---

## Configure environment

Create a `.env` file.

<details>

<summary><strong>.env example</strong></summary>

```env
PORT=8080

DATABASE_URL=postgres://user:password@localhost:5432/odemavpn

JWT_SECRET=your-secret

BOT_TOKEN=telegram_bot_token

LAVA_SHOP_ID=

LAVA_SECRET_KEY=

XUI_BASE_URL=https://panel.example.com

XUI_API_TOKEN=

XUI_INBOUND_ID=1

ServerIP=

ServerPort=

ServerPBK=

ServerSNI=

ServerSID=

SUB_URL=https://vpn.example.com/sub/
```

</details>

---

## Run locally

```bash
go run ./cmd/server
```

or

```bash
go build -o odemavpn ./cmd/server

./odemavpn
```

---

## Docker

At the moment the repository does **not** include Docker or Docker Compose configuration.

A minimal Dockerfile can be added without changing the application architecture.

---

# REST API Overview

| Endpoint | Description |
|-----------|-------------|
| POST /auth/telegram | Authenticate Telegram user |
| GET /subscription | Get subscription URL |
| POST /payment/create | Create payment invoice |
| POST /payment/webhook | Lava webhook |
| GET /profile | User profile |

---

# Subscription Lifecycle

```mermaid
stateDiagram-v2

[*] --> Registered

Registered --> PaymentPending

PaymentPending --> Active

Active --> Renewed

Renewed --> Active

Active --> Expired

Expired --> Disabled
```

---

# Security

- Telegram Login signature verification
- JWT authentication
- UUID-based VLESS clients
- Hidden Subscription IDs
- SQL transactions during registration
- HMAC validation for Lava webhooks
- Stateless authentication

---

# Internal Design Highlights

### Transactional User Creation

User registration is wrapped in a single SQL transaction, ensuring that user records and subscription records are created atomically.

---

### Idempotent Registration

Repeated Telegram logins do not create duplicate accounts thanks to SQL conflict handling.

---

### Bulk VPN Operations

Instead of enabling or disabling VPN clients individually, the backend utilizes 3X-UI bulk endpoints to reduce network overhead.

---

### Smart Subscription Extension

If an active subscription is renewed, additional days are appended to the current expiration date instead of restarting from the purchase date.

---

### Hidden Subscription URLs

The public subscription URL is generated using a random Subscription ID rather than exposing the actual VLESS UUID.

---

# Limitations

> [!IMPORTANT]
>
> The current implementation intentionally remains lightweight and therefore omits several production-oriented components.

Current limitations include:

- no Docker support
- no CI/CD pipeline
- no automated tests
- no Kubernetes deployment
- no Redis caching
- no connection pool (`pgx.Conn` instead of `pgxpool`)
- no structured logging
- no metrics or tracing

---

# Future Improvements

- Docker & Docker Compose
- pgxpool support
- OpenAPI / Swagger
- Prometheus metrics
- Structured logging
- Health checks
- Graceful shutdown
- Configuration validation
- Unit & integration tests
- GitHub Actions pipeline

---

# License

This project is distributed under the **MIT License**.

See the `LICENSE` file for more information.

---

<div align="center">

**OdemaVPN** — lightweight backend infrastructure for Telegram-powered VPN services.

</div>
