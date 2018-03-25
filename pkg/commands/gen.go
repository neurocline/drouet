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
func buildHugoGenCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "A collection of several useful generators.",
		RunE:  nil,
	}

	// This is lame - you need to add a command before you can add
	// its subcommands. This should be fixed.
	cmd.AddCommand(buildHugoGenAutocompleteCmd(h))
	cmd.AddCommand(buildHugoGenDocCmd(h))
	cmd.AddCommand(buildHugoGenManCmd(h))
	cmd.AddCommand(buildHugoGenDocsHelperCmd(h).cmd)
	cmd.AddCommand(buildHugoGenChromaStyles(h).cmd)

	return cmd
}

func buildHugoGenAutocompleteCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
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
	cmd.PersistentFlags().StringP("completionfile", "", "/etc/bash_completion.d/hugo.sh", "autocompletion file")
	cmd.PersistentFlags().StringP("type", "", "bash", "autocompletion type (currently only bash supported)")

	// For bash-completion
	cmd.PersistentFlags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})

	return cmd
}

func buildHugoGenDocCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Generate Markdown documentation for the Hugo CLI.",
		Long: `Generate Markdown documentation for the Hugo CLI.

This command is, mostly, used to create up-to-date documentation
of Hugo's command-line interface for http://gohugo.io/.

It creates one Markdown file per command with front matter suitable
for rendering in Hugo.`,
		RunE: h.genDoc,
	}

	cmd.PersistentFlags().String("dir", "/tmp/hugodoc/", "the directory to write the doc.")

	// For bash-completion
	cmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})

	return cmd
}

func buildHugoGenManCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "man",
		Short: "Generate man pages for the Hugo CLI",
		Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,
		RunE: h.genMan,
	}

	cmd.PersistentFlags().String("dir", "man/", "the directory to write the man pages.")

	// For bash-completion
	cmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})

	return cmd
}

func buildHugoGenDocsHelperCmd(h *hugoCmd) *genDocsHelper {
	g := &genDocsHelper{}

	g.cmd = &cobra.Command{
		Use:    "docshelper",
		Short:  "Generate some data files for the Hugo docs.",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.generate(h)
		},
	}

	// Add inherited flag, although I think this is a mistake, there's
	// no sub-commands for this to inherit
	g.cmd.PersistentFlags().StringVarP(&g.target, "dir", "", "docs/data", "data dir")

	// I bet this could safely return g.cmd
	return g
}

func buildHugoGenChromaStyles(h *hugoCmd) *genChromaStyles {
	g := &genChromaStyles{}

	g.cmd = &cobra.Command{
		Use:   "chromastyles",
		Short: "Generate CSS stylesheet for the Chroma code highlighter",
		Long: `Generate CSS stylesheet for the Chroma code highlighter for a given style.
This stylesheet is needed if pygmentsUseClasses is enabled in config.

See https://help.farbox.com/pygments.html for preview of available styles`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return g.generate(h)
		},
	}

	// Add inherited flag, although I think this is a mistake, there's
	// no sub-commands for this to inherit
	g.cmd.PersistentFlags().StringVar(&g.style, "style", "friendly", "highlighter style (see https://help.farbox.com/pygments.html)")
	g.cmd.PersistentFlags().StringVar(&g.highlightStyle, "highlightStyle", "bg:#ffffcc", "style used for highlighting lines (see https://github.com/alecthomas/chroma)")
	g.cmd.PersistentFlags().StringVar(&g.linesStyle, "linesStyle", "", "style used for line numbers (see https://github.com/alecthomas/chroma)")

	// I bet this could safely return g.cmd
	return g
}

// ----------------------------------------------------------------------------------------------

func (h *hugoCmd) genAutocomplete(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo gen autocomplete - hugo gen autocomplete code goes here")
	return nil
}

func (h *hugoCmd) genDoc(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo gen doc - hugo gen doc code goes here")
	return nil
}

func (h *hugoCmd) genMan(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo gen man - hugo gen man code goes here")
	return nil
}

type genDocsHelper struct {
	target string
	cmd    *cobra.Command
}

func (g *genDocsHelper) generate(h *hugoCmd) error {
	return nil
}

type genChromaStyles struct {
	style          string
	highlightStyle string
	linesStyle     string
	cmd            *cobra.Command
}

func (g *genChromaStyles) generate(h *hugoCmd) error {
	return nil
}
