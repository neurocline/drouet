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

// Build "hugo import" command.
func buildHugoImportCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import your site from others.",
		Long: `Import your site from other web site generators like Jekyll.

Import requires a subcommand, e.g. ` + "`hugo import jekyll <jekyll_root_path> <target_path>`.",
		RunE: nil,
	}

	cmd.AddCommand(buildHugoImportJekyllCmd(h))

	return cmd
}

// Build "hugo import" command.
func buildHugoImportJekyllCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jekyll",
		Short: "hugo import from Jekyll",
		Long: `hugo import from Jekyll.

Import from Jekyll requires two paths, e.g. ` + "`hugo import jekyll <jekyll_root_path> <target_path>`.",
		RunE: h.importJekyll,
	}

	cmd.Flags().Bool("force", false, "allow import into non-empty target directory")

	return cmd
}

// ----------------------------------------------------------------------------------------------

func (h *hugoCmd) importJekyll(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo import jekyll - hugo import jekyll code goes here")
	return nil
}
