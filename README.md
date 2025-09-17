# Go FIX Client for Coinbase Prime

## Introduction
This repository contains a lightweight Go-based FIX client that connects to Coinbase Prime's FIX gateway. It provides an interactive REPL to:
- Create new orders
- Look up existing orders (using a local `orders.json` cache)  
- Cancel orders
- View tabular formatted FIX messages

Under the hood, [QuickFIX/Go](https://github.com/quickfixgo/quickfix) is used to handle FIX message encoding/decoding and session management.

## Prerequisites
- **Go 1.23+** installed (https://golang.org/dl/)
- A valid **Coinbase Prime service account certificate** (PEM format) with private key
- A CA certificate bundle (e.g., `system-roots.pem`) to validate the TLS connection

---

## 1. Configure `fix.cfg` for Native TLS

Coinbase Prime FIX supports **native TLS**, so no stunnel or proxy is required.

### Quick Setup
Copy the example configuration file and rename it:
```bash
cp fix.cfg.example fix.cfg
```

Then edit `fix.cfg` to replace placeholder values with your actual credentials:
- Replace `YOUR_SVC_ACCOUNT_ID` with your service account ID
- Update the `SSLCAFile` path to point to your system's CA certificate bundle

### Generating CA Certificate Bundle
To generate a local CA certificate bundle from your system trust store, run:

```bash
security find-certificate -a -p /System/Library/Keychains/SystemRootCertificates.keychain > ~/system-roots.pem
```

This configuration enables QuickFIX/Go to connect directly over TLS without relying on external proxies like stunnel.

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
FIX> new <symbol> <MARKET|LIMIT|VWAP> <BUY|SELL> <BASE|QUOTE> <qty> [price] [start_time] [participation_rate] [expire_time]
```

#### Quantity Types
- **BASE**: Quantity specified in base currency (e.g., BTC for BTC-USD)
- **QUOTE**: Quantity specified in quote currency (e.g., USD for BTC-USD)

#### Examples

**Market Orders:**
```bash
# Buy 0.1 BTC (base currency)
FIX> new BTC-USD MARKET BUY BASE 0.1

# Buy $1000 worth of BTC (quote currency)
FIX> new BTC-USD MARKET BUY QUOTE 1000
```

**Limit Orders:**
```bash
# Buy 0.1 BTC at $30000 (base currency)
FIX> new BTC-USD LIMIT BUY BASE 0.1 30000

# Buy $3000 worth of BTC at $30000 (quote currency)
FIX> new BTC-USD LIMIT BUY QUOTE 3000 30000
```

**VWAP/TWAP Orders:**
You can specify VWAP orders with various combinations of optional parameters:

```bash
# Basic VWAP with just price (base currency)
FIX> new BTC-USD VWAP BUY BASE 1.0 50000

# VWAP with start time (quote currency)
FIX> new BTC-USD VWAP BUY QUOTE 50000 50000 2025-08-01T10:00:00Z

# VWAP with start time and participation rate (10%)
FIX> new BTC-USD VWAP BUY BASE 1.0 50000 2025-08-01T10:00:00Z 0.1

# VWAP with all parameters (start, participation rate, and expire time)
FIX> new BTC-USD VWAP BUY BASE 1.0 50000 2025-08-01T10:00:00Z 0.1 2025-08-01T16:00:00Z
```

**VWAP Parameters:**
- `start_time`: When execution should begin (ISO 8601 format)
- `participation_rate`: Execution aggressiveness (0.0-1.0, e.g., 0.1 = 10%)
- `expire_time`: When the order should expire (ISO 8601 format)

The order is sent, and the ExecReport (fill/cancel information) will be stored in `orders.json`.

### Look Up an Existing Order

```bash
FIX> status <ClOrdId> [OrderId] [Side] [Symbol]
```

This application automatically generates a unique `ClOrdId` (Client Order ID) using `UnixNano`. This value can be collected from `orders.json`, or from FIX responses sent by the server. `OrderId`, `Side`, and `Symbol` are required, however this app will automatically import these to the request based on the provided `ClOrdId`. 

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

This command lists out all stored orders from `orders.json`.

### Request for Quote (RFQ)

The client supports RFQ (Request for Quote) functionality for obtaining quotes before executing trades:

```bash
FIX> rfq <symbol> <BUY|SELL> <BASE|QUOTE> <qty> <price>
```

**Example:**
```bash
# Request quote to buy $15 worth of SOL-USD with limit price of $250
FIX> rfq SOL-USD BUY QUOTE 15 250
```

#### ⚠️ **Important RFQ Warning**

**The current implementation automatically accepts any RFQ quote that is received.** This is designed for demonstration purposes only. In a production environment, you would want to implement proper quote evaluation logic.

By providing a required limit price, the system sets a worst-case price and helps control the potential risks associated with auto-accepting quotes.

**Use caution when testing with real funds**, as the system will automatically execute trades upon receiving quotes.
