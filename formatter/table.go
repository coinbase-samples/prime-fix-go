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
	"fmt"
	"strconv"
	"strings"

	"github.com/quickfixgo/quickfix"
)

type FieldInfo struct {
	Tag         string
	Name        string
	Value       string
	Description string
}

var fixFieldDescriptions = map[string]string{
	"1":    "PortfolioID",
	"8":    "BeginString",
	"9":    "BodyLength",
	"10":   "CheckSum",
	"11":   "ClOrdID",
	"14":   "CumQty",
	"30":   "LastMkt",
	"31":   "LastPx",
	"32":   "LastShares",
	"34":   "MsgSeqNum",
	"35":   "MsgType",
	"37":   "OrderID",
	"38":   "OrderQty",
	"39":   "OrdStatus",
	"40":   "OrdType",
	"41":   "OrigClOrdID",
	"44":   "Price",
	"45":   "RefSeqNum",
	"49":   "SenderCompID",
	"50":   "SenderSubID",
	"52":   "SendingTime",
	"54":   "Side",
	"55":   "Symbol",
	"56":   "TargetCompID",
	"58":   "Text",
	"59":   "TimeInForce",
	"60":   "TransactTime",
	"79":   "PortfolioID",
	"96":   "SecureData",
	"98":   "EncryptMethod",
	"108":  "HeartBtInt",
	"126":  "ExpireTime",
	"141":  "ResetSeqNumFlag",
	"150":  "ExecType",
	"151":  "LeavesQty",
	"152":  "CashOrderQty",
	"168":  "StartTime",
	"371":  "RefTagID",
	"372":  "RefMsgType",
	"373":  "SessionRejectReason",
	"554":  "Password",
	"847":  "TargetStrategy",
	"849":  "ParticipationRate",
	"8002": "FilledAmount",
	"8006": "NetAvgPrice",
	"9406": "DropCopyFlag",
	"9407": "AccessKey",
}

var msgTypeDescriptions = map[string]string{
	"0": "HEARTBEAT",
	"1": "TEST_REQUEST",
	"2": "RESEND_REQUEST",
	"3": "REJECT",
	"4": "SEQUENCE_RESET",
	"5": "LOGOUT",
	"8": "EXECUTION_REPORT",
	"9": "ORDER_CANCEL_REJECT",
	"A": "LOGON",
	"D": "NEW_ORDER",
	"F": "ORDER_CANCEL_REQUEST",
	"G": "ORDER_CANCEL_REPLACE_REQUEST",
	"H": "ORDER_STATUS_REQUEST",
}

var ordStatusDescriptions = map[string]string{
	"0": "NEW",
	"1": "PARTIALLY_FILLED",
	"2": "FILLED",
	"3": "DONE_FOR_DAY",
	"4": "CANCELED",
	"5": "REPLACED",
	"6": "PENDING_CANCEL",
	"7": "STOPPED",
	"8": "REJECTED",
	"9": "SUSPENDED",
	"A": "PENDING_NEW",
	"B": "CALCULATED",
	"C": "EXPIRED",
	"D": "ACCEPTED_FOR_BIDDING",
	"E": "PENDING_REPLACE",
}

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
)

func FormatFixMessage(msg *quickfix.Message, direction string) string {
	var fields []FieldInfo

	// Process header fields
	for _, tag := range msg.Header.Tags() {
		if value, err := msg.Header.GetString(tag); err == nil {
			tagStr := strconv.Itoa(int(tag))
			name := getFieldName(tagStr)
			desc := getValueDescription(tagStr, value)
			fields = append(fields, FieldInfo{
				Tag:         tagStr,
				Name:        name,
				Value:       value,
				Description: desc,
			})
		}
	}

	// Process body fields
	for _, tag := range msg.Body.Tags() {
		if value, err := msg.Body.GetString(tag); err == nil {
			tagStr := strconv.Itoa(int(tag))
			name := getFieldName(tagStr)
			desc := getValueDescription(tagStr, value)
			fields = append(fields, FieldInfo{
				Tag:         tagStr,
				Name:        name,
				Value:       value,
				Description: desc,
			})
		}
	}

	// Process trailer fields
	for _, tag := range msg.Trailer.Tags() {
		if value, err := msg.Trailer.GetString(tag); err == nil {
			tagStr := strconv.Itoa(int(tag))
			name := getFieldName(tagStr)
			desc := getValueDescription(tagStr, value)
			fields = append(fields, FieldInfo{
				Tag:         tagStr,
				Name:        name,
				Value:       value,
				Description: desc,
			})
		}
	}

	return formatTable(fields, direction)
}

func getFieldName(tag string) string {
	if name, exists := fixFieldDescriptions[tag]; exists {
		return name
	}
	return fmt.Sprintf("Tag%s", tag)
}

func getValueDescription(tag, value string) string {
	if tag == "35" {
		if desc, exists := msgTypeDescriptions[value]; exists {
			return desc
		}
	} else if tag == "39" {
		if desc, exists := ordStatusDescriptions[value]; exists {
			return desc
		}
	}
	return value
}

func formatTable(fields []FieldInfo, direction string) string {
	if len(fields) == 0 {
		return ""
	}

	// Calculate column widths
	maxTag := 3
	maxName := 11
	maxValue := 5
	maxDesc := 17

	for _, field := range fields {
		if len(field.Tag) > maxTag {
			maxTag = len(field.Tag)
		}
		if len(field.Name) > maxName {
			maxName = len(field.Name)
		}
		if len(field.Value) > maxValue {
			maxValue = len(field.Value)
		}
		if len(field.Description) > maxDesc {
			maxDesc = len(field.Description)
		}
	}

	// Limit column widths for readability
	if maxName > 25 {
		maxName = 25
	}
	if maxValue > 40 {
		maxValue = 40
	}
	if maxDesc > 40 {
		maxDesc = 40
	}

	var sb strings.Builder

	// Direction indicator and table structure color
	directionColor := colorGreen
	tableColor := colorWhite
	arrow := "<---"
	if direction == "OUTGOING" {
		directionColor = colorBlue
		tableColor = colorBlue
		arrow = "--->"
	} else if direction == "INCOMING" {
		directionColor = colorYellow
		tableColor = colorYellow
		arrow = "<---"
	}

	sb.WriteString(fmt.Sprintf("%s%s %s%s%s\n", directionColor, arrow, colorBold, direction, colorReset))

	// Top border
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxTag+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxName+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxValue+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxDesc+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s\n", tableColor, colorReset))

	// Header row
	sb.WriteString(fmt.Sprintf("%s|%s %s%-*s%s %s|%s %s%-*s%s %s|%s %s%-*s%s %s|%s %s%-*s%s %s|%s\n",
		tableColor, colorReset, colorBold+colorCyan, maxTag, "TAG", colorReset,
		tableColor, colorReset, colorBold+colorCyan, maxName, "DESCRIPTION", colorReset,
		tableColor, colorReset, colorBold+colorCyan, maxValue, "VALUE", colorReset,
		tableColor, colorReset, colorBold+colorCyan, maxDesc, "VALUE DESCRIPTION", colorReset,
		tableColor, colorReset))

	// Header separator
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxTag+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxName+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxValue+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxDesc+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s\n", tableColor, colorReset))

	// Data rows
	for _, field := range fields {
		tagColor := colorWhite
		if field.Tag == "35" {
			tagColor = colorMagenta
		} else if field.Tag == "8" || field.Tag == "9" || field.Tag == "10" {
			tagColor = colorCyan
		} else if field.Tag == "39" {
			tagColor = colorYellow // OrdStatus
		} else if field.Tag == "14" || field.Tag == "31" || field.Tag == "32" || field.Tag == "8002" || field.Tag == "8006" {
			tagColor = colorGreen // Execution/fill related fields
		}

		name := field.Name
		if len(name) > maxName {
			name = name[:maxName-3] + "..."
		}

		value := field.Value
		if len(value) > maxValue {
			value = value[:maxValue-3] + "..."
		}

		desc := field.Description
		if len(desc) > maxDesc {
			desc = desc[:maxDesc-3] + "..."
		}

		sb.WriteString(fmt.Sprintf("%s|%s %s%*s%s %s|%s %-*s %s|%s %-*s %s|%s %-*s %s|%s\n",
			tableColor, colorReset, tagColor, maxTag, field.Tag, colorReset,
			tableColor, colorReset, maxName, name,
			tableColor, colorReset, maxValue, value,
			tableColor, colorReset, maxDesc, desc,
			tableColor, colorReset))
	}

	// Bottom border
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxTag+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxName+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxValue+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s", tableColor, colorReset))
	sb.WriteString(fmt.Sprintf("%s%s%s", tableColor, strings.Repeat("-", maxDesc+2), colorReset))
	sb.WriteString(fmt.Sprintf("%s+%s\n", tableColor, colorReset))

	return sb.String()
}
