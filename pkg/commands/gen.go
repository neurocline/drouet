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

	"github.com/neurocline/cobra"

	jww "github.com/spf13/jwalterweatherman"
)

// Build "hugo env" command.
func buildHugoGenCmd(hugo *commandeer) *hugoGenCmd {
	h := &hugoGenCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "gen",
		Short: "A collection of several useful generators.",
		RunE:  nil,
	}

	h.cmd.AddCommand(buildHugoGenAutocompleteCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenDocCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenManCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenDocsHelperCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoGenChromaStyles(hugo).cmd)

	return h
}

// Note - this only makes sense for Unix-like systems (Linux and Mac)
// and probably should be skipped for Windows
func buildHugoGenAutocompleteCmd(hugo *commandeer) *hugoGenAutocompleteCmd {
	h := &hugoGenAutocompleteCmd{c: hugo}

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
	h.cmd.PersistentFlags().StringVarP(&h.autocompleteTarget, "completionfile", "", "/etc/bash_completion.d/hugo.sh", "autocompletion file")
	h.cmd.PersistentFlags().StringVarP(&h.autocompleteType, "type", "", "bash", "autocompletion type (currently only bash supported)")

	// For bash-completion
	h.cmd.PersistentFlags().SetAnnotation("completionfile", cobra.BashCompFilenameExt, []string{})

	return h
}

func buildHugoGenDocCmd(hugo *commandeer) *hugoGenDocCmd {
	h := &hugoGenDocCmd{c: hugo}

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

	h.cmd.PersistentFlags().StringVar(&h.gendocdir, "dir", "/tmp/hugodoc/", "the directory to write the doc.")

	// For bash-completion
	h.cmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})

	return h
}

func buildHugoGenManCmd(hugo *commandeer) *hugoGenManCmd {
	h := &hugoGenManCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "man",
		Short: "Generate man pages for the Hugo CLI",
		Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,
		RunE: h.genMan,
	}

	h.cmd.PersistentFlags().StringVar(&h.genmandir, "dir", "man/", "the directory to write the man pages.")

	// For bash-completion
	h.cmd.PersistentFlags().SetAnnotation("dir", cobra.BashCompSubdirsInDir, []string{})

	return h
}

func buildHugoGenDocsHelperCmd(hugo *commandeer) *hugoGenDocsHelperCmd {
	h := &hugoGenDocsHelperCmd{c: hugo}

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

func buildHugoGenChromaStyles(hugo *commandeer) *hugoGenChromaStylesCmd {
	h := &hugoGenChromaStylesCmd{c: hugo}

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
	c *commandeer
	cmd *cobra.Command
}

type hugoGenAutocompleteCmd struct {
	c *commandeer
	cmd *cobra.Command

	// Where to write the completion file
	autocompleteTarget string

	// bash for now (zsh and others will come)
	autocompleteType string
}

func (h *hugoGenAutocompleteCmd) genAutocomplete(cmd *cobra.Command, args []string) error {
	if h.autocompleteType != "bash" {
		return core.NewUserError("Only Bash is supported for now")
	}

	err := h.cmd.Root().GenBashCompletionFile(h.autocompleteTarget)

	if err != nil {
		return err
	}

	jww.FEEDBACK.Println("Bash completion file for Hugo saved to", h.autocompleteTarget)

	return nil
}

type hugoGenDocCmd struct {
	c *commandeer
	cmd *cobra.Command

	gendocdir string
}

func (h *hugoGenDocCmd) genDoc(cmd *cobra.Command, args []string) error {
	const gendocFrontmatterTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`
	return nil
}

type hugoGenManCmd struct {
	c *commandeer
	cmd *cobra.Command

	genmandir string
}

func (h *hugoGenManCmd) genMan(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo gen man - hugo gen man code goes here")
	return nil
}

type hugoGenDocsHelperCmd struct {
	c *commandeer
	cmd    *cobra.Command

	target string
}

func (h *hugoGenDocsHelperCmd) generate(cmd *cobra.Command, args []string) error {
	return nil
}

type hugoGenChromaStylesCmd struct {
	c *commandeer
	cmd            *cobra.Command

	style          string
	highlightStyle string
	linesStyle     string
}

func (h *hugoGenChromaStylesCmd) generate(cmd *cobra.Command, args []string) error {
	return nil
}
