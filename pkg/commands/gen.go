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

// Build "hugo env" command.
func buildHugoGenCmd(hugo *core.Hugo) *hugoGenCmd {
	h := &hugoGenCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "gen",
		Short: "A collection of several useful generators.",
		RunE:  nil,
	}

	// This is lame - you need to add a command before you can add
	// its subcommands. This should be fixed.
	h.cmd.AddCommand(buildHugoGenAutocompleteCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenDocCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenManCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenDocsHelperCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenChromaStyles(hugo).cmd)

	return h
}

func buildHugoGenAutocompleteCmd(hugo *core.Hugo) *hugoGenAutocompleteCmd {
	h := &hugoGenAutocompleteCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "autocomplete",
		Short: "Generate shell autocompletion script for Hugo",
		Long: `Generates a shell autocompletion script for Hugo.

NOTE: The current version supports Bash only.
      This should work for *nix systems with Bash installed.

By default, the file is written directly to /etc/bash_completion.d
for convenience, and the command may need superuser rights, e.g.:

	$ sudo hugo gen autocomplete

Add ` + "`--completionfile=/path/to/file`" + ` flag to set alternative
file-path and name.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion`,
		RunE: h.genAutocomplete,
	}

	// I think these are being misused
	h.cmd.PersistentFlags().StringP("completionfile", "", "/etc/bash_completion.d/hugo.sh", "autocompletion file")
	h.cmd.PersistentFlags().StringP("type", "", "bash", "autocompletion type (currently only bash supported)")

	// For bash-completion
	h.cmd.PersistentFlags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})

	return h
}

func buildHugoGenDocCmd(hugo *core.Hugo) *hugoGenDocCmd {
	h := &hugoGenDocCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "doc",
		Short: "Generate Markdown documentation for the Hugo CLI.",
		Long: `Generate Markdown documentation for the Hugo CLI.

This command is, mostly, used to create up-to-date documentation
of Hugo's command-line interface for http://gohugo.io/.

It creates one Markdown file per command with front matter suitable
for rendering in Hugo.`,
		RunE: h.genDoc,
	}

	h.cmd.PersistentFlags().String("dir", "/tmp/hugodoc/", "the directory to write the doc.")

	// For bash-completion
	h.cmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})

	return h
}

func buildHugoGenManCmd(hugo *core.Hugo) *hugoGenManCmd {
	h := &hugoGenManCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "man",
		Short: "Generate man pages for the Hugo CLI",
		Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,
		RunE: h.genMan,
	}

	h.cmd.PersistentFlags().String("dir", "man/", "the directory to write the man pages.")

	// For bash-completion
	h.cmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})

	return h
}

func buildHugoGenDocsHelperCmd(hugo *core.Hugo) *hugoGenDocsHelperCmd {
	h := &hugoGenDocsHelperCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:    "docshelper",
		Short:  "Generate some data files for the Hugo docs.",
		Hidden: true,
		RunE:  h.generate,
	}

	// Add inherited flag, although I think this is a mistake, there's
	// no sub-commands for this to inherit
	h.cmd.PersistentFlags().StringVarP(&h.target, "dir", "", "docs/data", "data dir")

	return h
}

func buildHugoGenChromaStyles(hugo *core.Hugo) *hugoGenChromaStylesCmd {
	h := &hugoGenChromaStylesCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "chromastyles",
		Short: "Generate CSS stylesheet for the Chroma code highlighter",
		Long: `Generate CSS stylesheet for the Chroma code highlighter for a given style.
This stylesheet is needed if pygmentsUseClasses is enabled in config.

See https://help.farbox.com/pygments.html for preview of available styles`,
		RunE: h.generate,
	}

	// Add inherited flag, although I think this is a mistake, there's
	// no sub-commands for this to inherit
	h.cmd.PersistentFlags().StringVar(&h.style, "style", "friendly", "highlighter style (see https://help.farbox.com/pygments.html)")
	h.cmd.PersistentFlags().StringVar(&h.highlightStyle, "highlightStyle", "bg:#ffffcc", "style used for highlighting lines (see https://github.com/alecthomas/chroma)")
	h.cmd.PersistentFlags().StringVar(&h.linesStyle, "linesStyle", "", "style used for line numbers (see https://github.com/alecthomas/chroma)")

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoGenCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

type hugoGenAutocompleteCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoGenAutocompleteCmd) genAutocomplete(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo gen autocomplete - hugo gen autocomplete code goes here")
	return nil
}

type hugoGenDocCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoGenDocCmd) genDoc(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo gen doc - hugo gen doc code goes here")
	return nil
}

type hugoGenManCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoGenManCmd) genMan(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo gen man - hugo gen man code goes here")
	return nil
}

type hugoGenDocsHelperCmd struct {
	*core.Hugo
	cmd    *cobra.Command

	target string
}

func (h *hugoGenDocsHelperCmd) generate(cmd *cobra.Command, args []string) error {
	return nil
}

type hugoGenChromaStylesCmd struct {
	*core.Hugo
	cmd            *cobra.Command

	style          string
	highlightStyle string
	linesStyle     string
}

func (h *hugoGenChromaStylesCmd) generate(cmd *cobra.Command, args []string) error {
	return nil
}
