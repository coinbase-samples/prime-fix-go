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

To generate a local CA certificate bundle from your system trust store, run:

```bash
security find-certificate -a -p /System/Library/Keychains/SystemRootCertificates.keychain > ~/system-roots.pem
```

Edit or create the `fix.cfg` file at the project root. Replace the `SenderCompID` with your actual service account ID, and adjust the path to your CA file:

```
SSLCAFile=/Users/yourname/system-roots.pem
SenderCompID=YOUR_SVC_ACCOUNT_ID
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

---

## 🎨 Enhanced Field Support

### Expanded FIX Tag Dictionary
The formatter now includes descriptions for 40+ common FIX fields including:

**Core Message Fields:**
- Tag 35 (MsgType): `"D"` → `"NEW_ORDER"`
- Tag 39 (OrdStatus): `"2"` → `"FILLED"`
- Tag 49 (SenderCompID), Tag 56 (TargetCompID)

**Execution & Fill Data:**
- Tag 14 (CumQty) - Cumulative quantity filled
- Tag 31 (LastPx) - Last execution price
- Tag 32 (LastShares) - Last execution quantity
- Tag 8002 (FilledAmount) - Total filled amount
- Tag 8006 (NetAvgPrice) - Net average price

**Trading Information:**
- Tag 50 (SenderSubID), Tag 60 (TransactTime)
- Tag 79 (PortfolioID), Tag 30 (LastMkt)

### Color Coding
- **Message Type (Tag 35)** - Magenta for easy identification
- **Order Status (Tag 39)** - Yellow highlighting
- **Execution Fields** - Green for fill/trade data (Tags 14, 31, 32, 8002, 8006)
- **Header/Trailer** - Cyan for structural fields (Tags 8, 9, 10)

**Unknown fields are handled gracefully** - any field not in the lookup table displays as `"Tag999"` without breaking the formatter.