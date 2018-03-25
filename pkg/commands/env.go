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

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build "hugo env" command.
func buildHugoEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Print Hugo version and environment info",
		Long:  `Print Hugo version and environment info. This is useful in Hugo bug reports.`,
		RunE:  env,
	}

	return cmd
}

// ----------------------------------------------------------------------------------------------

func env(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo env - hugo env code goes here")
	return nil
}
