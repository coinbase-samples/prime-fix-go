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
	"prime-fix-go/utils"
)

func (a *FixApp) handleNew(parts []string) {
	if len(parts) < 5 {
		fmt.Println("usage: new <symbol> <MARKET|LIMIT> <BUY|SELL> <qty> [price]")
		return
	}
	msg := builder.BuildNew(
		parts[1],                    // symbol
		parts[2],                    // ordType
		parts[3],                    // side
		parts[4],                    // qty
		utils.GetOptional(parts, 5), // price (optional)
		a.PortfolioId,
	)
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
