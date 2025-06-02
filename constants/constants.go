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

package constants

import "github.com/quickfixgo/quickfix"

const (
	MsgTypeNew    = "D" // New Order
	MsgTypeStatus = "H" // Status
	MsgTypeCancel = "F" // Cancel

	FixTimeFormat = "20060102-15:04:05.000"

	DefaultTargetCompID = "COIN"

	OrdTypeLimit  = "LIMIT"
	OrdTypeMarket = "MARKET"
	SideBuy       = "BUY"
	SideSell      = "SELL"

	OrdTypeLimitFix      = "2" // Limit order type
	OrdTypeMarketFix     = "1" // Market order type
	TimeInForceDay       = "1" // Day
	TimeInForceIoc       = "3" // Immediate or Cancel
	TargetStrategyLimit  = "L" // Limit strategy
	TargetStrategyMarket = "M" // Market strategy
	SideBuyFix           = "1" // Buy side
	SideSellFix          = "2" // Sell side

	TagAccount        = quickfix.Tag(1)
	TagClOrdId        = quickfix.Tag(11)
	TagOrderId        = quickfix.Tag(37)
	TagOrderQty       = quickfix.Tag(38)
	TagOrdType        = quickfix.Tag(40)
	TagOrigClOrdId    = quickfix.Tag(41)
	TagTargetStrategy = quickfix.Tag(847)
	TagPx             = quickfix.Tag(44)
	TagSenderCompId   = quickfix.Tag(49)
	TagSendingTime    = quickfix.Tag(52)
	TagSide           = quickfix.Tag(54)
	TagSymbol         = quickfix.Tag(55)
	TagTargetCompId   = quickfix.Tag(56)
	TagTimeInForce    = quickfix.Tag(59)
	TagHmac           = quickfix.Tag(96)
	TagMsgType        = quickfix.Tag(35)
	TagExecType       = quickfix.Tag(150)
	TagCashOrderQty   = quickfix.Tag(152)
	TagPassword       = quickfix.Tag(554)
	TagDropCopyFlag   = quickfix.Tag(9406)
	TagAccessKey      = quickfix.Tag(9407)
)
