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
	"prime-fix-go/utils"
	"strings"
	"time"

	"prime-fix-go/constants"
	"prime-fix-go/model"

	"github.com/quickfixgo/quickfix"
)

func BuildNew(
	symbol, ordType, side, qtyType, qty, price, portfolio string, vwapParams ...string,
) *quickfix.Message {
	m := quickfix.NewMessage()
	m.Header.SetField(constants.TagMsgType, quickfix.FIXString(constants.MsgTypeNew))
	m.Header.SetField(constants.TagSenderCompId, quickfix.FIXString(os.Getenv("SVC_ACCOUNT_ID")))
	m.Header.SetField(constants.TagTargetCompId, quickfix.FIXString(os.Getenv("TARGET_COMP_ID")))
	m.Header.SetField(constants.TagSendingTime, quickfix.FIXString(time.Now().UTC().Format(constants.FixTimeFormat)))

	clId := fmt.Sprintf("%d", time.Now().UnixNano())
	m.Body.SetField(constants.TagAccount, quickfix.FIXString(portfolio))
	m.Body.SetField(constants.TagClOrdId, quickfix.FIXString(clId))
	m.Body.SetField(constants.TagSymbol, quickfix.FIXString(symbol))

	// Set quantity based on user preference (BASE or QUOTE)
	if strings.EqualFold(qtyType, "BASE") {
		m.Body.SetField(constants.TagOrderQty, quickfix.FIXString(qty))
	} else { // Default to QUOTE
		m.Body.SetField(constants.TagCashOrderQty, quickfix.FIXString(qty))
	}

	if strings.EqualFold(ordType, constants.OrdTypeLimit) {
		m.Body.SetField(constants.TagOrdType, quickfix.FIXString(constants.OrdTypeLimitFix))
		m.Body.SetField(constants.TagTimeInForce, quickfix.FIXString(constants.TimeInForceDay))
		m.Body.SetField(constants.TagPx, quickfix.FIXString(price))
		m.Body.SetField(constants.TagTargetStrategy, quickfix.FIXString(constants.TargetStrategyLimit))
	} else if strings.EqualFold(ordType, constants.OrdTypeVwap) {
		m.Body.SetField(constants.TagOrdType, quickfix.FIXString(constants.OrdTypeVwapFix))
		m.Body.SetField(constants.TagTimeInForce, quickfix.FIXString(constants.TimeInForceGtd))
		m.Body.SetField(constants.TagPx, quickfix.FIXString(price))
		m.Body.SetField(constants.TagTargetStrategy, quickfix.FIXString(constants.TargetStrategyVwap))

		if len(vwapParams) > 0 && vwapParams[0] != "" {
			effectiveTime, err := time.Parse("2006-01-02T15:04:05Z", vwapParams[0])
			if err == nil {
				m.Body.SetField(constants.TagStartTime, quickfix.FIXString(effectiveTime.Format(constants.FixTimeFormat)))
			} else {
				m.Body.SetField(constants.TagStartTime, quickfix.FIXString(vwapParams[0]))
			}
		}

		hasParticipationRate := len(vwapParams) > 1 && vwapParams[1] != ""
		hasExpireTime := len(vwapParams) > 2 && vwapParams[2] != ""

		if hasParticipationRate {
			m.Body.SetField(constants.TagParticipationRate, quickfix.FIXString(vwapParams[1]))
		}
		if hasExpireTime {
			expireTime, err := time.Parse("2006-01-02T15:04:05Z", vwapParams[2])
			if err == nil {
				m.Body.SetField(constants.TagExpireTime, quickfix.FIXString(expireTime.Format(constants.FixTimeFormat)))
			} else {
				m.Body.SetField(constants.TagExpireTime, quickfix.FIXString(vwapParams[2]))
			}
		} else {
			defaultExpire, _ := time.Parse("2006-01-02T15:04:05Z", "2025-07-26T23:59:59Z")
			m.Body.SetField(constants.TagExpireTime, quickfix.FIXString(defaultExpire.Format(constants.FixTimeFormat)))
		}
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
	m.Header.SetField(constants.TagSenderCompId, quickfix.FIXString(os.Getenv("SVC_ACCOUNT_ID")))
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
	m.Header.SetField(constants.TagSenderCompId, quickfix.FIXString(os.Getenv("SVC_ACCOUNT_ID")))
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

func BuildLogon(
	body *quickfix.Body,
	ts, apiKey, apiSecret, passphrase, targetCompId, portfolioId string,
) {
	sig := utils.Sign(ts, "A", "1", apiKey, targetCompId, passphrase, apiSecret)

	body.SetField(constants.TagAccount, quickfix.FIXString(portfolioId))
	body.SetField(constants.TagHmac, quickfix.FIXString(sig))
	body.SetField(constants.TagPassword, quickfix.FIXString(passphrase))
	body.SetField(constants.TagDropCopyFlag, quickfix.FIXString("Y"))
	body.SetField(constants.TagAccessKey, quickfix.FIXString(apiKey))
}
