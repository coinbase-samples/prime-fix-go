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

package model

type OrderInfo struct {
	ClOrdId           string `json:"clOrdId"`
	OrderId           string `json:"orderId"`
	Side              string `json:"side"`
	Symbol            string `json:"symbol"`
	Quantity          string `json:"quantity"`
	LimitPrice        string `json:"limitPrice"`
	StartTime         string `json:"startTime,omitempty"`
	ExpireTime        string `json:"expireTime,omitempty"`
	ParticipationRate string `json:"participationRate,omitempty"`
}

type QuoteRequestInfo struct {
	QuoteReqId string `json:"quoteReqId"`
	Account    string `json:"account"`
	Side       string `json:"side"`
	Symbol     string `json:"symbol"`
	OrderQty   string `json:"orderQty"`
	Price      string `json:"price"`
}

type QuoteInfo struct {
	QuoteId        string `json:"quoteId"`
	QuoteReqId     string `json:"quoteReqId"`
	Account        string `json:"account"`
	Symbol         string `json:"symbol"`
	BidPx          string `json:"bidPx,omitempty"`
	OfferPx        string `json:"offerPx,omitempty"`
	BidSize        string `json:"bidSize,omitempty"`
	OfferSize      string `json:"offerSize,omitempty"`
	ValidUntilTime string `json:"validUntilTime"`
}
