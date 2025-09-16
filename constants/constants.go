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
	MsgTypeNew      = "D" // New Order
	MsgTypeStatus   = "H" // Status
	MsgTypeCancel   = "F" // Cancel
	MsgTypeLogon    = "A" // Logon
	MsgTypeQuoteReq = "R" // Quote Request
	MsgTypeQuote    = "S" // Quote
	MsgTypeQuoteAck = "b" // Quote Acknowledgment

	FixTimeFormat = "20060102-15:04:05.000"

	DefaultTargetCompID = "COIN"

	OrdTypeLimit  = "LIMIT"
	OrdTypeMarket = "MARKET"
	OrdTypeVwap   = "VWAP"
	OrdTypeRfq    = "RFQ"
	SideBuy       = "BUY"
	SideSell      = "SELL"

	OrdTypeLimitFix         = "2" // Limit order type
	OrdTypeMarketFix        = "1" // Market order type
	OrdTypeVwapFix          = "2" // VWAP order type (uses limit type)
	OrdTypePreviouslyQuoted = "D" // Previously Quoted (for RFQ accept)
	TimeInForceDay          = "1" // Day
	TimeInForceIoc          = "3" // Immediate or Cancel
	TimeInForceGtd          = "6" // Good Till Date
	TimeInForceFok          = "4" // Fill or Kill (for RFQ)
	TargetStrategyLimit     = "L" // Limit strategy
	TargetStrategyMarket    = "M" // Market strategy
	TargetStrategyVwap      = "V" // VWAP strategy
	TargetStrategyRfq       = "R" // RFQ strategy
	SideBuyFix              = "1" // Buy side
	SideSellFix             = "2" // Sell side

	TagAccount           = quickfix.Tag(1)
	TagClOrdId           = quickfix.Tag(11)
	TagOrderId           = quickfix.Tag(37)
	TagOrderQty          = quickfix.Tag(38)
	TagOrdType           = quickfix.Tag(40)
	TagOrigClOrdId       = quickfix.Tag(41)
	TagTargetStrategy    = quickfix.Tag(847)
	TagPx                = quickfix.Tag(44)
	TagExecInst          = quickfix.Tag(18)
	TagSenderCompId      = quickfix.Tag(49)
	TagSendingTime       = quickfix.Tag(52)
	TagSide              = quickfix.Tag(54)
	TagSymbol            = quickfix.Tag(55)
	TagTargetCompId      = quickfix.Tag(56)
	TagText              = quickfix.Tag(58)
	TagTimeInForce       = quickfix.Tag(59)
	TagValidUntilTime    = quickfix.Tag(62)
	TagHmac              = quickfix.Tag(96)
	TagQuoteId           = quickfix.Tag(117)
	TagQuoteReqId        = quickfix.Tag(131)
	TagBidPx             = quickfix.Tag(132)
	TagOfferPx           = quickfix.Tag(133)
	TagBidSize           = quickfix.Tag(134)
	TagOfferSize         = quickfix.Tag(135)
	TagQuoteAckStatus    = quickfix.Tag(297)
	TagQuoteRejectReason = quickfix.Tag(300)
	TagMsgType           = quickfix.Tag(35)
	TagExecType          = quickfix.Tag(150)
	TagCashOrderQty      = quickfix.Tag(152)
	TagPassword          = quickfix.Tag(554)
	TagDropCopyFlag      = quickfix.Tag(9406)
	TagAccessKey         = quickfix.Tag(9407)
	TagStartTime         = quickfix.Tag(168)
	TagExpireTime        = quickfix.Tag(126)
	TagParticipationRate = quickfix.Tag(849)

	QuoteAckStatusRejected = "5"
)
