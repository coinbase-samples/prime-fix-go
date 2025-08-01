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

package formatter

import (
	"strings"
	"testing"

	"github.com/quickfixgo/quickfix"
	"prime-fix-go/constants"
)

func TestFormatFixMessage(t *testing.T) {
	// Create a sample FIX message
	msg := quickfix.NewMessage()

	// Add header fields
	msg.Header.SetString(constants.TagSenderCompId, "COIN")
	msg.Header.SetString(constants.TagTargetCompId, "TEST_CLIENT")
	msg.Header.SetString(constants.TagMsgType, "3") // Reject

	// Add body fields
	msg.Body.SetString(quickfix.Tag(45), "1")           // RefSeqNum
	msg.Body.SetString(quickfix.Tag(58), "auth failed") // Text
	msg.Body.SetString(quickfix.Tag(371), "10")         // RefTagID
	msg.Body.SetString(quickfix.Tag(372), "A")          // RefMsgType
	msg.Body.SetString(quickfix.Tag(373), "abc")        // SessionRejectReason

	// Add trailer field
	msg.Trailer.SetString(quickfix.Tag(10), "049") // CheckSum

	// Format the message
	result := FormatFixMessage(msg, "INCOMING")

	// Basic checks
	if result == "" {
		t.Error("Expected non-empty formatted message")
	}

	// Check for required elements
	if !strings.Contains(result, "INCOMING") {
		t.Error("Expected direction indicator in output")
	}

	if !strings.Contains(result, "TAG") {
		t.Error("Expected table header in output")
	}

	if !strings.Contains(result, "REJECT") {
		t.Error("Expected message type description in output")
	}

	if !strings.Contains(result, "auth failed") {
		t.Error("Expected field value in output")
	}
}
