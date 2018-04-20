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
	"github.com/neurocline/cobra"
)

// Build "hugo check" command.
func buildHugoCheckCmd(hugo *commandeer) *hugoCheckCmd {
	h := &hugoCheckCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "check",
		Args:  cobra.NoArgs,
		Short: "Contains some verification checks",
		RunE:  h.check, // If we have a nil RunE, we get no error checking
	}

	// Add the check ulimit subcommand if it exists (only for Darwin)
	if cmdUlimit := buildHugoCheckUlimitCmd(hugo); cmdUlimit != nil {
		h.cmd.AddCommand(cmdUlimit.cmd)
	}

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoCheckCmd struct {
	c *commandeer
	cmd *cobra.Command
}

func (h *hugoCheckCmd) check(cmd *cobra.Command, args []string) error {
	return nil
}
