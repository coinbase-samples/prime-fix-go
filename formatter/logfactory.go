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
	"bytes"
	"fmt"
	"strings"

	"github.com/quickfixgo/quickfix"
)

type TableLogFactory struct{}

func NewTableLogFactory() *TableLogFactory {
	return &TableLogFactory{}
}

func (f *TableLogFactory) Create() (quickfix.Log, error) {
	return &TableLog{}, nil
}

func (f *TableLogFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	return &TableLog{SessionID: sessionID}, nil
}

type TableLog struct {
	SessionID quickfix.SessionID
}

func (l *TableLog) OnIncoming(msg []byte) {
	message := quickfix.NewMessage()
	buffer := bytes.NewBuffer(msg)
	if err := quickfix.ParseMessage(message, buffer); err == nil {
		formatted := FormatFixMessage(message, "INCOMING")
		fmt.Print(formatted)
	} else {
		fmt.Printf("Error parsing incoming message: %s\nRaw: %s\n", err, string(msg))
	}
}

func (l *TableLog) OnOutgoing(msg []byte) {
	message := quickfix.NewMessage()
	buffer := bytes.NewBuffer(msg)
	if err := quickfix.ParseMessage(message, buffer); err == nil {
		formatted := FormatFixMessage(message, "OUTGOING")
		fmt.Print(formatted)
	} else {
		fmt.Printf("Error parsing outgoing message: %s\nRaw: %s\n", err, string(msg))
	}
}

func (l *TableLog) OnEvent(msg string) {
	if !strings.Contains(msg, "Sending") && !strings.Contains(msg, "Received") {
		fmt.Printf("Event: %s\n", msg)
	}
}

func (l *TableLog) OnEventf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if !strings.Contains(msg, "Sending") && !strings.Contains(msg, "Received") {
		fmt.Printf("Event: %s\n", msg)
	}
}
