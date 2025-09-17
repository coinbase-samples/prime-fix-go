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
	"sync"
	"time"

	"prime-fix-go/builder"
	"prime-fix-go/constants"
	"prime-fix-go/model"
	"prime-fix-go/utils"

	"github.com/quickfixgo/quickfix"
)

const orderFile = "orders.json"

type FixApp struct {
	SessionId quickfix.SessionID
	orders    map[string]model.OrderInfo
	config    *constants.Config
	mu        sync.RWMutex
}

func NewFixApp(config *constants.Config) *FixApp {
	return &FixApp{
		SessionId: quickfix.SessionID{},
		orders:    make(map[string]model.OrderInfo),
		config:    config,
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
	fmt.Println("Commands: new, status, cancel, list, rfq, version, exit")
}

func (a *FixApp) ToAdmin(msg *quickfix.Message, _ quickfix.SessionID) {
	if t, _ := msg.Header.GetString(constants.TagMsgType); t == constants.MsgTypeLogon {
		ts := time.Now().UTC().Format(constants.FixTimeFormat)
		builder.BuildLogon(
			&msg.Body,
			ts,
			a.config.AccessKey,
			a.config.SigningKey,
			a.config.Passphrase,
			a.config.TargetCompId,
			a.config.PortfolioId,
		)
	}
}

func (a *FixApp) FromApp(msg *quickfix.Message, _ quickfix.SessionID) quickfix.MessageRejectError {
	msgType, _ := msg.Header.GetString(constants.TagMsgType)
	switch msgType {
	case "8":
		a.handleExecReport(msg)
	case constants.MsgTypeQuote:
		a.handleQuote(msg)
	case constants.MsgTypeQuoteAck:
		a.handleQuoteAck(msg)
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

	a.mu.Lock()
	defer a.mu.Unlock()

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

func (a *FixApp) handleQuote(msg *quickfix.Message) {
	quote := model.QuoteInfo{
		QuoteId:        utils.GetString(msg, constants.TagQuoteId),
		QuoteReqId:     utils.GetString(msg, constants.TagQuoteReqId),
		Account:        utils.GetString(msg, constants.TagAccount),
		Symbol:         utils.GetString(msg, constants.TagSymbol),
		BidPx:          utils.GetString(msg, constants.TagBidPx),
		OfferPx:        utils.GetString(msg, constants.TagOfferPx),
		BidSize:        utils.GetString(msg, constants.TagBidSize),
		OfferSize:      utils.GetString(msg, constants.TagOfferSize),
		ValidUntilTime: utils.GetString(msg, constants.TagValidUntilTime),
	}

	if quote.QuoteId == "" {
		return
	}

	log.Printf("✓ received quote %s for request %s", quote.QuoteId, quote.QuoteReqId)

	if quote.BidPx != "" {
		fmt.Printf("Quote: Bid %s @ %s (valid until %s)\n", quote.BidSize, quote.BidPx, quote.ValidUntilTime)
	}
	if quote.OfferPx != "" {
		fmt.Printf("Quote: Offer %s @ %s (valid until %s)\n", quote.OfferSize, quote.OfferPx, quote.ValidUntilTime)
	}

	// Auto-accept the quote
	a.autoAcceptQuote(quote)
}

func (a *FixApp) autoAcceptQuote(quote model.QuoteInfo) {
	var price, qty, side string
	if quote.BidPx != "" {
		price = quote.BidPx
		qty = quote.BidSize
		side = "SELL"
	} else if quote.OfferPx != "" {
		price = quote.OfferPx
		qty = quote.OfferSize
		side = "BUY"
	} else {
		log.Printf("✗ cannot auto-accept quote %s: no valid bid or offer price", quote.QuoteId)
		return
	}

	acceptMsg := builder.BuildAcceptQuote(quote.QuoteId, quote.Symbol, side, qty, price, a.config.PortfolioId, a.config)
	err := quickfix.SendToTarget(acceptMsg, a.SessionId)
	if err != nil {
		return
	}
	log.Printf("✓ auto-accepting quote %s: %s %s %s @ %s", quote.QuoteId, side, qty, quote.Symbol, price)
}

func (a *FixApp) handleQuoteAck(msg *quickfix.Message) {
	quoteReqId := utils.GetString(msg, constants.TagQuoteReqId)
	quoteAckStatus := utils.GetString(msg, constants.TagQuoteAckStatus)
	rejectReason := utils.GetString(msg, constants.TagQuoteRejectReason)
	text := utils.GetString(msg, constants.TagText)

	if quoteAckStatus == constants.QuoteAckStatusRejected {
		log.Printf("✗ quote request %s rejected: reason=%s, text=%s", quoteReqId, rejectReason, text)
	} else {
		log.Printf("? quote acknowledgment for %s: status=%s", quoteReqId, quoteAckStatus)
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
			app.handleNew(parts)
		case "status":
			app.handleStatus(parts)
		case "cancel":
			app.handleCancel(parts)
		case "list":
			app.handleList()
		case "rfq":
			app.handleRfq(parts)
		case "version":
			fmt.Println(utils.FullVersion())
		case "exit":
			return
		default:
			fmt.Println("unknown command")
		}
	}
}

func (a *FixApp) saveOrders() error {
	data, err := json.MarshalIndent(a.orders, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal orders: %w", err)
	}
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
