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

package builder

import (
	"fmt"
	"os"
	"strings"
	"time"

	"prime-fix-go/constants"
	"prime-fix-go/model"

	"github.com/quickfixgo/quickfix"
)

func BuildNew(
	symbol, ordType, side, qty, price, portfolio string,
) *quickfix.Message {
	m := quickfix.NewMessage()
	m.Header.SetField(constants.TagMsgType, quickfix.FIXString(constants.MsgTypeNew))
	m.Header.SetField(constants.TagSenderCompId, quickfix.FIXString(os.Getenv("SVC_ACCOUNTID")))
	m.Header.SetField(constants.TagTargetCompId, quickfix.FIXString(os.Getenv("TARGET_COMP_ID")))
	m.Header.SetField(constants.TagSendingTime, quickfix.FIXString(time.Now().UTC().Format(constants.FixTimeFormat)))

	clId := fmt.Sprintf("%d", time.Now().UnixNano())
	m.Body.SetField(constants.TagAccount, quickfix.FIXString(portfolio))
	m.Body.SetField(constants.TagClOrdId, quickfix.FIXString(clId))
	m.Body.SetField(constants.TagSymbol, quickfix.FIXString(symbol))
	m.Body.SetField(constants.TagOrderQty, quickfix.FIXString(qty))

	if strings.EqualFold(ordType, constants.OrdTypeLimit) {
		m.Body.SetField(constants.TagOrdType, quickfix.FIXString(constants.OrdTypeLimitFix))
		m.Body.SetField(constants.TagTimeInForce, quickfix.FIXString(constants.TimeInForceDay))
		m.Body.SetField(constants.TagPx, quickfix.FIXString(price))
		m.Body.SetField(constants.TagTargetStrategy, quickfix.FIXString(constants.TargetStrategyLimit))
	} else {
		m.Body.SetField(constants.TagOrdType, quickfix.FIXString(constants.OrdTypeMarketFix))
		m.Body.SetField(constants.TagTimeInForce, quickfix.FIXString(constants.TimeInForceIoc))
		m.Body.SetField(constants.TagTargetStrategy, quickfix.FIXString(constants.TargetStrategyMarket))
	}

	if strings.EqualFold(side, constants.SideBuy) {
		m.Body.SetField(constants.TagSide, quickfix.FIXString(constants.SideBuyFix))
	} else {
		m.Body.SetField(constants.TagSide, quickfix.FIXString(constants.SideSellFix))
	}

	return m
}

func BuildStatus(clId, ordId, side, symbol string) *quickfix.Message {
	m := quickfix.NewMessage()
	m.Header.SetField(constants.TagMsgType, quickfix.FIXString(constants.MsgTypeStatus))
	m.Header.SetField(constants.TagSenderCompId, quickfix.FIXString(os.Getenv("SVC_ACCOUNTID")))
	m.Header.SetField(constants.TagTargetCompId, quickfix.FIXString(os.Getenv("TARGET_COMP_ID")))
	m.Header.SetField(constants.TagSendingTime, quickfix.FIXString(time.Now().UTC().Format(constants.FixTimeFormat)))

	m.Body.SetField(constants.TagClOrdId, quickfix.FIXString(clId))
	m.Body.SetField(constants.TagOrderId, quickfix.FIXString(ordId))
	m.Body.SetField(constants.TagSide, quickfix.FIXString(side))
	m.Body.SetField(constants.TagSymbol, quickfix.FIXString(symbol))
	return m
}

func BuildCancel(info model.OrderInfo, portfolio string) *quickfix.Message {
	m := quickfix.NewMessage()
	m.Header.SetField(constants.TagMsgType, quickfix.FIXString(constants.MsgTypeCancel))
	m.Header.SetField(constants.TagSenderCompId, quickfix.FIXString(os.Getenv("SVC_ACCOUNTID")))
	m.Header.SetField(constants.TagTargetCompId, quickfix.FIXString(os.Getenv("TARGET_COMP_ID")))
	m.Header.SetField(constants.TagSendingTime, quickfix.FIXString(time.Now().UTC().Format(constants.FixTimeFormat)))

	cancelClId := fmt.Sprintf("cancel-%d", time.Now().UnixNano())
	m.Body.SetField(constants.TagAccount, quickfix.FIXString(portfolio))
	m.Body.SetField(constants.TagClOrdId, quickfix.FIXString(cancelClId))
	m.Body.SetField(constants.TagOrigClOrdId, quickfix.FIXString(info.ClOrdId))
	m.Body.SetField(constants.TagOrderId, quickfix.FIXString(info.OrderId))
	m.Body.SetField(constants.TagOrderQty, quickfix.FIXString(info.Quantity))
	m.Body.SetField(constants.TagSide, quickfix.FIXString(info.Side))
	m.Body.SetField(constants.TagSymbol, quickfix.FIXString(info.Symbol))
	return m
}
