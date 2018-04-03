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

	"github.com/neurocline/drouet/pkg/core"

	"github.com/spf13/cobra"
)

// Build "hugo new" command.
func buildHugoNewCmd(hugo *core.Hugo) *hugoNewCmd {
	h := &hugoNewCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "new [path]",
		Short: "Create new content for your site",
		Long: `Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.

You can also specify the kind with ` + "`-k KIND`" + `.

If archetypes are provided in your theme or site, they will be used.`,
		RunE: h.newContent,
	}

	// Add flags used by all subcommands of new
	h.cmd.PersistentFlags().StringP("source", "s", "", "filesystem path to read files relative from")

	// Add flags used by just new
	h.cmd.Flags().StringP("kind", "k", "", "content type to create")
	h.cmd.Flags().String("editor", "", "edit new content with this editor, if provided")

	// Set BASH expansion
	h.cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})

	// Add subcommands
	h.cmd.AddCommand(buildHugoNewSiteCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoNewThemeCmd(hugo).cmd)

	return h
}

// Build "hugo new site" command.
func buildHugoNewSiteCmd(hugo *core.Hugo) *hugoNewSiteCmd {
	h := &hugoNewSiteCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "site [path]",
		Short: "Create a new site (skeleton)",
		Long: `Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use ` + "`hugo new [contentPath]`" + ` to create new content.`,
		RunE: h.newSite,
	}

	return h
}

// Build "hugo new theme" command.
func buildHugoNewThemeCmd(hugo *core.Hugo) *hugoNewThemeCmd {
	h := &hugoNewThemeCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "theme [name]",
		Short: "Create a new theme",
		Long: `Create a new theme (skeleton) called [name] in the current directory.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
as you see fit.`,
		RunE: h.newTheme,
	}

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoNewCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoNewCmd) newContent(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo new - hugo new code goes here")
	return nil
}

type hugoNewSiteCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoNewSiteCmd) newSite(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo new site - hugo new site code goes here")
	return nil
}

type hugoNewThemeCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoNewThemeCmd) newTheme(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo new theme - hugo new theme code goes here")
	return nil
}
