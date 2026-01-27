<p align="center">
  <img src="assets/logo.png" />
</p>

The central management service for **The Beyond**. A monorepo containing the Telegram Bot, Open API, and backend logic for user authentication, billing, and traffic accounting.

[![License: MIT](https://img.shields.io/badge/License-MIT-99FF00.svg?style=for-the-badge&labelColor=020617)](https://github.com/thebeyond-net/control-plane/blob/main/LICENSE)
[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Postgresql](https://img.shields.io/badge/PostgreSQL-%23326791.svg?style=for-the-badge&logo=postgresql&logoColor=white)](https://postgresql.org)
[![Wiki](https://img.shields.io/badge/Docs-%23FF0.svg?style=for-the-badge&logo=wikibooks&logoColor=020617)](https://github.com/thebeyond-net/control-plane/wiki)

---

# 🔥 Features
- **VPN Management** – User interface via Telegram Bot and client apps
- **Flexible Billing** – Support for Telegram Invoices and [LAVA](https://lava.ru) payments
- **Open API** – Public HTTP interface for mobile and desktop apps
- **Traffic Telemetry** – gRPC-based usage accounting with 2-hour sync interval
- **Deployment** – Simple setup via `docker compose up -d --build`

# 🚀 Getting Started
To run the project locally:

### 1. Clone the Repository
```sh
git clone https://github.com/thebeyond-net/control-plane.git
cd control-plane
```
### 2. Prepare Environment
```sh
cp .env.example .env
```
### 3. Launch with Docker
```sh
docker compose up -d --build
```

> [!IMPORTANT]
> See the [Wiki](https://github.com/thebeyond-net/control-plane/wiki) for production setup, environment variables, and security.

---
### Credits
Badges by [shields.io](https://shields.io).