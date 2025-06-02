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
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"prime-fix-go/builder"
	"prime-fix-go/constants"
	"prime-fix-go/model"
	"prime-fix-go/utils"

	"github.com/quickfixgo/quickfix"
)

const orderFile = "orders.json"

type FixApp struct {
	ApiKey       string
	ApiSecret    string
	Passphrase   string
	SenderCompId string
	TargetCompId string
	PortfolioId  string

	SessionId quickfix.SessionID
	orders    map[string]model.OrderInfo
}

func NewFixApp(
	apiKey, apiSecret, passphrase,
	senderCompId, targetCompId, portfolioId string,
) *FixApp {
	return &FixApp{
		ApiKey:       apiKey,
		ApiSecret:    apiSecret,
		Passphrase:   passphrase,
		SenderCompId: senderCompId,
		TargetCompId: targetCompId,
		PortfolioId:  portfolioId,
		orders:       make(map[string]model.OrderInfo),
	}
}

func (a *FixApp) OnCreate(sid quickfix.SessionID) {
	a.SessionId = sid
}

func (a *FixApp) OnLogout(sid quickfix.SessionID) {
	log.Println("Logout", sid)
}

func (a *FixApp) FromAdmin(_ *quickfix.Message, _ quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *FixApp) ToApp(_ *quickfix.Message, _ quickfix.SessionID) error {
	return nil
}

func (a *FixApp) OnLogon(sid quickfix.SessionID) {
	a.SessionId = sid
	log.Println("✓ FIX logon", sid)
	if err := a.loadOrders(); err != nil {
		log.Println("order cache load err:", err)
	}
	fmt.Println("Commands: new, status, cancel, list, exit")
}

func (a *FixApp) ToAdmin(msg *quickfix.Message, _ quickfix.SessionID) {
	if t, _ := msg.Header.GetString(constants.TagMsgType); t == "A" {
		ts := time.Now().UTC().Format(constants.FixTimeFormat)
		sig := utils.Sign(ts, "A", "1", a.ApiKey, a.TargetCompId, a.Passphrase, a.ApiSecret)
		msg.Body.SetField(constants.TagAccount, quickfix.FIXString(a.PortfolioId))
		msg.Body.SetField(constants.TagHmac, quickfix.FIXString(sig))
		msg.Body.SetField(constants.TagPassword, quickfix.FIXString(a.Passphrase))
		msg.Body.SetField(constants.TagDropCopyFlag, quickfix.FIXString("Y"))
		msg.Body.SetField(constants.TagAccessKey, quickfix.FIXString(a.ApiKey))
	}
}

func (a *FixApp) FromApp(msg *quickfix.Message, _ quickfix.SessionID) quickfix.MessageRejectError {
	if t, _ := msg.Header.GetString(constants.TagMsgType); t == "8" {
		a.handleExecReport(msg)
	}
	return nil
}

func (a *FixApp) handleExecReport(msg *quickfix.Message) {
	info := model.OrderInfo{
		ClOrdId:    utils.GetString(msg, constants.TagClOrdId),
		OrderId:    utils.GetString(msg, constants.TagOrderId),
		Side:       utils.GetString(msg, constants.TagSide),
		Symbol:     utils.GetString(msg, constants.TagSymbol),
		Quantity:   utils.GetString(msg, constants.TagOrderQty),
		LimitPrice: utils.GetString(msg, constants.TagPx),
	}
	if info.Quantity == "" {
		info.Quantity = utils.GetString(msg, constants.TagCashOrderQty)
	}
	if info.ClOrdId == "" {
		return
	}

	existing, exists := a.orders[info.ClOrdId]
	if !exists || (existing.OrderId == "" && info.OrderId != "") {
		if info.OrderId == "" {
			info.OrderId = existing.OrderId
		}
		a.orders[info.ClOrdId] = info
		_ = a.saveOrders()
		log.Printf("⇡ cached/updated %s (OrderId %s)", info.ClOrdId, info.OrderId)
	}
	if _, exists2 := a.orders[info.ClOrdId]; !exists2 {
		if a.orders == nil {
			a.orders = make(map[string]model.OrderInfo)
		}
		a.orders[info.ClOrdId] = info
		_ = a.saveOrders()
		log.Printf("⇡ cached %s (OrderId %s)", info.ClOrdId, info.OrderId)
	}
}

// Commands: new, status, cancel, list, exit.
func Repl(app *FixApp) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("FIX> ")
		line, _ := reader.ReadString('\n')
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) == 0 {
			continue
		}
		cmd := strings.ToLower(parts[0])
		switch cmd {

		case "new":
			if len(parts) < 5 {
				fmt.Println("usage: new <symbol> <MARKET|LIMIT> <BUY|SELL> <qty> [price]")
				continue
			}
			msg := builder.BuildNew(
				parts[1],                    // symbol
				parts[2],                    // ordType
				parts[3],                    // side
				parts[4],                    // qty
				utils.GetOptional(parts, 5), // price (optional)
				app.PortfolioId,
			)
			quickfix.SendToTarget(msg, app.SessionId)

		case "status":
			if len(parts) < 2 {
				fmt.Println("usage: status <ClOrdId> [OrderId] [Side] [Symbol]")
				continue
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
			if cached, ok := app.orders[cl]; ok {
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
				continue
			}
			quickfix.SendToTarget(builder.BuildStatus(cl, ord, side, sym), app.SessionId)

		case "cancel":
			if len(parts) < 2 {
				fmt.Println("usage: cancel <ClOrdId>")
				continue
			}
			info, ok := app.orders[parts[1]]
			if !ok {
				fmt.Println("unknown ClOrdId (not in cache)")
				continue
			}
			quickfix.SendToTarget(builder.BuildCancel(info, app.PortfolioId), app.SessionId)

		case "list":
			if len(app.orders) == 0 {
				fmt.Println("(no cached orders)")
				continue
			}
			for _, o := range app.orders {
				fmt.Printf("%-20s → %s (%s %s %s)\n",
					o.ClOrdId, o.OrderId, o.Side, o.Symbol, o.Quantity)
			}

		case "exit":
			return

		default:
			fmt.Println("unknown command")
		}
	}
}

func (a *FixApp) saveOrders() error {
	data, _ := json.MarshalIndent(a.orders, "", "  ")
	return os.WriteFile(orderFile, data, 0o644)
}

func (a *FixApp) loadOrders() error {
	data, err := os.ReadFile(orderFile)
	if err != nil {
		if os.IsNotExist(err) {
			a.orders = make(map[string]model.OrderInfo)
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &a.orders)
}
