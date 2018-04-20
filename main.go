// Copyright 2015 The Hugo Authors. All rights reserved.
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

package main

import (
	"os"

	"github.com/neurocline/drouet/pkg/commands"
)

// Current design for code - no use of os.Exit anywhere in the code except
// at the top level. Calls to panic are allowed, and if not caught will
// result in a somewhat clean exit (because panic respects defer calls,
// whereas os.Exit does not).

func main() {

	// Run the command processor
	statusCode := commands.Execute()

	// Default for Go is to return 0. If we want to return a non-zero
	// status code, we do it with os.Exit.
	if statusCode != 0 {
		os.Exit(statusCode)
	}
}
