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
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/quickfixgo/quickfix"
)

// csvColumnTags defines the FIX tags to include as columns, in order.
// Each column will have the field name as header (from fixFieldDescriptions).
var csvColumnTags = []string{
	// Header fields
	"8",  // BeginString
	"9",  // BodyLength
	"35", // MsgType
	"34", // MsgSeqNum
	"49", // SenderCompId
	"50", // SenderSubID
	"52", // SendingTime
	"56", // TargetCompId
	// Order identification
	"11", // ClOrdID
	"37", // OrderID
	"41", // OrigClOrdID
	"17", // ExecID
	// Instrument
	"55", // Symbol
	// Order details
	"54",  // Side
	"38",  // OrderQty
	"152", // CashOrderQty
	"44",  // Price
	"40",  // OrdType
	"59",  // TimeInForce
	// Execution details
	"39",  // OrdStatus
	"150", // ExecType
	"20",  // ExecTransType
	"14",  // CumQty
	"151", // LeavesQty
	"6",   // AvgPx
	"31",  // LastPx
	"32",  // LastShares
	"80",  // AllocQty
	"30",  // LastMkt
	"60",  // TransactTime
	// Commission
	"12",  // Commission
	"13",  // CommType
	"120", // SettlCurrency
	// Account/Allocation
	"1",  // PortfolioId/Account
	"78", // NoAllocs
	"79", // AllocAccount
	// Strategy
	"847", // TargetStrategy
	"849", // ParticipationRate
	"168", // StartTime
	"126", // ExpireTime
	// Custom fields
	"8002", // FilledAmount
	"8006", // NetAvgPrice
	// Text/errors
	"58", // Text
	// Trailer
	"10", // CheckSum
}

type CsvLogger struct {
	mu       sync.Mutex
	file     *os.File
	writer   *csv.Writer
	enabled  bool
	filename string
}

var csvLogger = &CsvLogger{}

func StartCsvLogging() (string, error) {
	csvLogger.mu.Lock()
	defer csvLogger.mu.Unlock()

	if csvLogger.enabled {
		return csvLogger.filename, nil
	}

	filename := fmt.Sprintf("fix_messages_%s.csv", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}

	csvLogger.file = file
	csvLogger.writer = csv.NewWriter(file)
	csvLogger.enabled = true
	csvLogger.filename = filename

	// Build header row
	header := []string{"Direction", "Timestamp"}
	for _, tag := range csvColumnTags {
		name := getFieldName(tag)
		header = append(header, fmt.Sprintf("%s (%s)", name, tag))
	}
	header = append(header, "Other Fields")

	if err := csvLogger.writer.Write(header); err != nil {
		csvLogger.file.Close()
		csvLogger.enabled = false
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}
	csvLogger.writer.Flush()

	return filename, nil
}

func StopCsvLogging() error {
	csvLogger.mu.Lock()
	defer csvLogger.mu.Unlock()

	if !csvLogger.enabled {
		return nil
	}

	csvLogger.writer.Flush()
	err := csvLogger.file.Close()
	csvLogger.enabled = false
	csvLogger.file = nil
	csvLogger.writer = nil

	return err
}

func IsCsvLoggingEnabled() bool {
	csvLogger.mu.Lock()
	defer csvLogger.mu.Unlock()
	return csvLogger.enabled
}

func GetCsvFilename() string {
	csvLogger.mu.Lock()
	defer csvLogger.mu.Unlock()
	return csvLogger.filename
}

func LogMessageToCsv(msg *quickfix.Message, direction string) {
	csvLogger.mu.Lock()
	defer csvLogger.mu.Unlock()

	if !csvLogger.enabled || csvLogger.writer == nil {
		return
	}

	timestamp := time.Now().Format(time.RFC3339Nano)

	// Build a map of all field values from the message
	fieldValues := make(map[string]string)

	// Header fields
	for _, tag := range msg.Header.Tags() {
		if value, err := msg.Header.GetString(tag); err == nil {
			tagStr := strconv.Itoa(int(tag))
			fieldValues[tagStr] = value
		}
	}

	// Body fields
	for _, tag := range msg.Body.Tags() {
		if value, err := msg.Body.GetString(tag); err == nil {
			tagStr := strconv.Itoa(int(tag))
			fieldValues[tagStr] = value
		}
	}

	// Trailer fields
	for _, tag := range msg.Trailer.Tags() {
		if value, err := msg.Trailer.GetString(tag); err == nil {
			tagStr := strconv.Itoa(int(tag))
			fieldValues[tagStr] = value
		}
	}

	// Build a set of known tags for quick lookup
	knownTags := make(map[string]bool)
	for _, tag := range csvColumnTags {
		knownTags[tag] = true
	}

	// Build row with fixed column order
	row := []string{direction, timestamp}
	for _, tag := range csvColumnTags {
		value := fieldValues[tag] // Empty string if not present
		row = append(row, value)
	}

	// Collect any fields not in the predefined columns
	var otherFields []string
	for tag, value := range fieldValues {
		if !knownTags[tag] {
			name := getFieldName(tag)
			otherFields = append(otherFields, fmt.Sprintf("%s(%s)=%s", name, tag, value))
		}
	}
	row = append(row, strings.Join(otherFields, "; "))

	if err := csvLogger.writer.Write(row); err != nil {
		fmt.Printf("Error writing to CSV: %v\n", err)
		return
	}
	csvLogger.writer.Flush()
}
