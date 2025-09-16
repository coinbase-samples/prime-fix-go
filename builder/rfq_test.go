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
	"os"
	"testing"

	"prime-fix-go/constants"
)

func TestBuildQuoteRequest(t *testing.T) {
	os.Setenv("SVC_ACCOUNT_ID", "test-sender")
	os.Setenv("TARGET_COMP_ID", "COIN")

	msg := BuildQuoteRequest("BTC-USD", "BUY", "BASE", "1.0", "50000", "test-portfolio")

	if msg == nil {
		t.Fatal("BuildQuoteRequest returned nil")
	}

	msgType, _ := msg.Header.GetString(constants.TagMsgType)
	if msgType != constants.MsgTypeQuoteReq {
		t.Errorf("Expected message type %s, got %s", constants.MsgTypeQuoteReq, msgType)
	}

	symbol, _ := msg.Body.GetString(constants.TagSymbol)
	if symbol != "BTC-USD" {
		t.Errorf("Expected symbol BTC-USD, got %s", symbol)
	}

	side, _ := msg.Body.GetString(constants.TagSide)
	if side != constants.SideBuyFix {
		t.Errorf("Expected side %s, got %s", constants.SideBuyFix, side)
	}

	qty, _ := msg.Body.GetString(constants.TagOrderQty)
	if qty != "1.0" {
		t.Errorf("Expected quantity 1.0, got %s", qty)
	}

	price, _ := msg.Body.GetString(constants.TagPx)
	if price != "50000" {
		t.Errorf("Expected price 50000, got %s", price)
	}

	ordType, _ := msg.Body.GetString(constants.TagOrdType)
	if ordType != constants.OrdTypeLimitFix {
		t.Errorf("Expected order type %s, got %s", constants.OrdTypeLimitFix, ordType)
	}

	tif, _ := msg.Body.GetString(constants.TagTimeInForce)
	if tif != constants.TimeInForceFok {
		t.Errorf("Expected time in force %s, got %s", constants.TimeInForceFok, tif)
	}
}

func TestBuildAcceptQuote(t *testing.T) {
	os.Setenv("SVC_ACCOUNT_ID", "test-sender")
	os.Setenv("TARGET_COMP_ID", "COIN")

	msg := BuildAcceptQuote("quote123", "BTC-USD", "SELL", "1.0", "49500", "test-portfolio")

	if msg == nil {
		t.Fatal("BuildAcceptQuote returned nil")
	}

	msgType, _ := msg.Header.GetString(constants.TagMsgType)
	if msgType != constants.MsgTypeNew {
		t.Errorf("Expected message type %s, got %s", constants.MsgTypeNew, msgType)
	}

	quoteId, _ := msg.Body.GetString(constants.TagQuoteId)
	if quoteId != "quote123" {
		t.Errorf("Expected quote ID quote123, got %s", quoteId)
	}

	ordType, _ := msg.Body.GetString(constants.TagOrdType)
	if ordType != constants.OrdTypePreviouslyQuoted {
		t.Errorf("Expected order type %s, got %s", constants.OrdTypePreviouslyQuoted, ordType)
	}

	targetStrategy, _ := msg.Body.GetString(constants.TagTargetStrategy)
	if targetStrategy != constants.TargetStrategyRfq {
		t.Errorf("Expected target strategy %s, got %s", constants.TargetStrategyRfq, targetStrategy)
	}

	tif, _ := msg.Body.GetString(constants.TagTimeInForce)
	if tif != constants.TimeInForceFok {
		t.Errorf("Expected time in force %s, got %s", constants.TimeInForceFok, tif)
	}
}
