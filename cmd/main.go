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

package main

import (
	"fmt"
	"log"

	"prime-fix-go/constants"
	"prime-fix-go/fixclient"
	"prime-fix-go/formatter"
	"prime-fix-go/utils"

	"github.com/quickfixgo/quickfix"
)

func main() {
	fmt.Printf("%s\n\n", utils.FullVersion())

	settings, err := utils.LoadSettings("fix.cfg")
	if err != nil {
		log.Fatal(err)
	}

	config := constants.NewConfig()
	app := fixclient.NewFixApp(config)

	initiator, err := quickfix.NewInitiator(app,
		quickfix.NewMemoryStoreFactory(),
		settings,
		formatter.NewTableLogFactory(),
	)
	if err != nil {
		log.Fatal("initiator error:", err)
	}

	if err := initiator.Start(); err != nil {
		log.Fatal("start error:", err)
	}
	defer initiator.Stop()

	fixclient.Repl(app)
}
