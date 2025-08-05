/**
 * Copyright 2025-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fixclient

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"prime-fix-go/builder"
	"strconv"
	"strings"
)

func (a *FixApp) handleNew(parts []string) {
	if len(parts) < 6 {
		fmt.Println("error: insufficient arguments")
		fmt.Println("usage: new <symbol> <MARKET|LIMIT|VWAP> <BUY|SELL> <BASE|QUOTE> <qty> [price] [start_time] [participation_rate] [expire_time]")
		return
	}

	symbol := parts[1]
	ordType := strings.ToUpper(parts[2])
	side := strings.ToUpper(parts[3])
	qtyType := strings.ToUpper(parts[4])
	qty := parts[5]
	var price string

	if side != "BUY" && side != "SELL" {
		fmt.Println("error: side must be BUY or SELL")
		return
	}

	if qtyType != "BASE" && qtyType != "QUOTE" {
		fmt.Println("error: quantity type must be BASE or QUOTE")
		return
	}

	if _, err := strconv.ParseFloat(qty, 64); err != nil {
		fmt.Println("error: qty must be a valid number")
		return
	}

	switch ordType {
	case "MARKET":
		if len(parts) > 6 {
			fmt.Println("error: MARKET orders should not include a price")
			return
		}
	case "LIMIT":
		if len(parts) < 7 {
			fmt.Println("error: price must be specified for LIMIT orders")
			return
		}
		price = parts[6]
		if _, err := strconv.ParseFloat(price, 64); err != nil {
			fmt.Println("error: price must be a valid number")
			return
		}
	case "VWAP":
		if len(parts) < 7 {
			fmt.Println("error: price must be specified for VWAP orders")
			return
		}
		price = parts[6]
		if _, err := strconv.ParseFloat(price, 64); err != nil {
			fmt.Println("error: price must be a valid number")
			return
		}
	default:
		fmt.Println("error: order type must be MARKET, LIMIT, or VWAP")
		return
	}

	var vwapParams []string
	if ordType == "VWAP" && len(parts) > 7 {
		vwapParams = parts[7:]
	}

	msg, err := builder.BuildNew(symbol, ordType, side, qtyType, qty, price, a.PortfolioId, vwapParams...)
	if err != nil {
		fmt.Printf("Error building order: %v\n", err)
		return
	}
	quickfix.SendToTarget(msg, a.SessionId)
}

func (a *FixApp) handleStatus(parts []string) {
	if len(parts) < 2 {
		fmt.Println("usage: status <ClOrdId> [OrderId] [Side] [Symbol]")
		return
	}
	cl := parts[1]
	var ord, side, sym string
	if len(parts) > 2 {
		ord = parts[2]
	}
	if len(parts) > 3 {
		side = parts[3]
	}
	if len(parts) > 4 {
		sym = parts[4]
	}
	if cached, ok := a.orders[cl]; ok {
		if ord == "" {
			ord = cached.OrderId
		}
		if side == "" {
			side = cached.Side
		}
		if sym == "" {
			sym = cached.Symbol
		}
	}
	if ord == "" || side == "" || sym == "" {
		fmt.Println("need OrderId, Side, and Symbol (not cached)")
		return
	}
	quickfix.SendToTarget(builder.BuildStatus(cl, ord, side, sym), a.SessionId)
}

func (a *FixApp) handleCancel(parts []string) {
	if len(parts) < 2 {
		fmt.Println("usage: cancel <ClOrdId>")
		return
	}
	info, ok := a.orders[parts[1]]
	if !ok {
		fmt.Println("unknown ClOrdId (not in cache)")
		return
	}
	quickfix.SendToTarget(builder.BuildCancel(info, a.PortfolioId), a.SessionId)
}

func (a *FixApp) handleList() {
	if len(a.orders) == 0 {
		fmt.Println("(no cached orders)")
		return
	}
	for _, o := range a.orders {
		fmt.Printf("%-20s → %s (%s %s %s)\n",
			o.ClOrdId, o.OrderId, o.Side, o.Symbol, o.Quantity)
	}
}
