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

package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/quickfixgo/quickfix"
)

func GetString(msg *quickfix.Message, tag quickfix.Tag) string {
	v, err := msg.Body.GetString(tag)
	if err != nil {
		return ""
	}
	return v
}

func Sign(ts, msgType, seq, key, tgt, pass, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(ts + msgType + seq + key + tgt + pass))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func GetOptional(parts []string, idx int) string {
	if len(parts) > idx {
		return parts[idx]
	}
	return ""
}

func LoadSettings(path string) (*quickfix.Settings, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return quickfix.ParseSettings(f)
}
