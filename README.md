# Go FIX Client for Coinbase Prime

## Introduction
This repository contains a lightweight Go-based FIX client that connects to Coinbase Prime’s FIX gateway via stunnel. It provides an interactive REPL to:
- Create new orders
- Look up existing orders (using a local `orders.json` cache)
- Cancel orders

Under the hood, [QuickFIX/Go](https://github.com/quickfixgo/quickfix) is used to handle FIX message encoding/decoding and session management, while stunnel establishes a secure TLS tunnel to Prime.

## Prerequisites
- **Go 1.20+** installed (https://golang.org/dl/)
- **stunnel (v5.74+)** installed (see https://www.stunnel.org/downloads.html or install via Homebrew: `brew install stunnel`)
- A valid **Coinbase Prime service account certificate** (PEM format) with private key

---

## 1. Configure & Run stunnel

Coinbase Prime requires a TLS tunnel on `localhost:4198` that forwards to `fix.prime.coinbase.com:4198`.

A sample configuration is included in this repo at `resources/stunnel.conf`. You can use it as-is or modify it as needed.

To install Prime’s TLS certificate into `/resources`:

```bash
mkdir -p resources
openssl s_client -showcerts -connect fix.prime.coinbase.com:4198 < /dev/null \
  | openssl x509 -outform PEM > resources/fix-prime.coinbase.com.pem
```

Then start stunnel:

```bash
stunnel resources/stunnel.conf
```

## 2. Configure `fix.cfg`

Edit `fix.cfg` at the project root. Update the `[SESSION]` block to include your actual service account ID under `SenderCompID`. Example:

```bash
SenderCompID=YOUR_SERVICE_ACCOUNT_ID
```

## 3. API credentials

Your Go FIX client also requires a few environment variables to sign the FIX Logon. Set the following in your shell before running:

```bash
export ACCESS_KEY="your_api_access_key"
export SIGNING_KEY="your_api_secret_key"
export PASSPHRASE="your_api_passphrase"
export TARGET_COMP_ID="COIN"
export PORTFOLIO_ID="your_portfolio_id"
export SVC_ACCOUNT_ID="your_service_account_id"
```

## 4. Build & Run the Go FIX Client

Ensure dependencies:

```bash
go mod tidy
```

Run the client:
```bash
go run cmd/main.go
```

On successful FIX Logon, you'll see:

```bash
FIX logon SessionID[YOUR_SENDER->COIN]
Commands: new, status, cancel, list, exit
```

## 5. REPL Commands

Once the client is running, type one of the following at the `FIX>` prompt:

### Create a New Order

```bash
FIX> new <symbol> <MARKET|LIMIT> <BUY|SELL> <qty> [price]
```

Example (Limit Buy 0.1 BTC at $30000):

```bash
FIX> new BTC-USD LIMIT BUY 0.1 30000
```

The order is sent, and the ExecReport (fill/cancel information) will be cached in `orders.json`.

### Look Up an Existing Order

```bash
FIX> status <ClOrdId> [OrderId] [Side] [Symbol]
```

This application automatically generates a unique `ClOrdId` using `UnixNano`. This value can be collected from `orders.json`. `OrderId`, `Side`, and `Symbol` are required, however this app will automatically import these to the request based on the provided `ClOrdId`. 

Example:

```bash
FIX> status 1685727281712345678
```
If `orders.json` contains that `ClOrdId`, its `OrderId`, `Side`, and `Symbol` are filled in automatically.

### Cancel an order

```bash
FIX> cancel <ClOrdID>
```

This request looks up an order by `ClOrdId` and attempts to cancel it.

### List All Cached Orders

```bash
FIX> list
```

This command lists out all cached orders from `orders.json`. 