// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package commands handles Hugo command-line processing.
//
// Hugo commands and flags are implemented using Cobra. Each
// top-level subcommand is in a source file named after the command.
package commands

import (
	"fmt"
	"runtime"

	jww "github.com/spf13/jwalterweatherman"
)

// Execute builds a command processor and runs the user command.
func Execute() int {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// do work
	fmt.Println("This is where work would be done")

	// If we have any log messages to leveLError or higher, return this
	// as an error code from the process
	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		return -1
	}

	if Hugo != nil {
		if Hugo.Log.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
			return -1
		}
	}

	return 0
}
