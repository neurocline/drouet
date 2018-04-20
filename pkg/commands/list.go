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
	//"fmt"
	"path/filepath"

	"github.com/neurocline/drouet/pkg/core"
	"github.com/neurocline/drouet/pkg/hugolib"

	"github.com/neurocline/cobra"

	jww "github.com/spf13/jwalterweatherman"
)

// Build "hugo list" command.
func buildHugoListCmd(hugo *commandeer) *hugoListCmd {
	h := &hugoListCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "list",
		Short: "Listing out various types of content",
		Long: `Listing out various types of content.

List requires a subcommand, e.g. ` + "`hugo list drafts`.",
		RunE: nil,
	}

	// Add hugo list subcommands
	h.cmd.AddCommand(buildHugoListDraftsCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoListExpiredCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoListFutureCmd(hugo).cmd)

	// Add flags used by all list subcommands
	h.cmd.PersistentFlags().StringP("source", "s", "", "filesystem path to read files relative from")

	// Set bash completion
	h.cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})

	return h
}

// Build "hugo list drafts" command.
func buildHugoListDraftsCmd(hugo *commandeer) *hugoListCmd {
	h := &hugoListCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "drafts",
		Short: "List all drafts",
		Long:  `List all of the drafts in your content directory.`,
		RunE:  h.listDrafts,
	}

	return h
}

// Build "hugo list expired" command.
func buildHugoListExpiredCmd(hugo *commandeer) *hugoListCmd {
	h := &hugoListCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "expired",
		Short: "List all posts already expired",
		Long: `List all of the posts in your content directory which has already
expired.`,
		RunE: h.listExpired,
	}

	return h
}

// Build "hugo list future" command.
func buildHugoListFutureCmd(hugo *commandeer) *hugoListCmd {
	h := &hugoListCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "future",
		Short: "List all posts dated in the future",
		Long: `List all of the posts in your content directory which will be
posted in the future.`,
		RunE: h.listFuture,
	}

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoListCmd struct {
	c *commandeer
	cmd *cobra.Command
}

func (h *hugoListCmd) listDrafts(cmd *cobra.Command, args []string) error {
	cfgInit := func(c *commandeer) error {
		c.Set("buildDrafts", true)
		return nil
	}
	return h.listContent(cfgInit, (*hugolib.Page).IsDraft)
}

func (h *hugoListCmd) listExpired(cmd *cobra.Command, args []string) error {
	cfgInit := func(c *commandeer) error {
		c.Set("buildExpired", true)
		return nil
	}
	return h.listContent(cfgInit, (*hugolib.Page).IsExpired)
}

func (h *hugoListCmd) listFuture(cmd *cobra.Command, args []string) error {
	cfgInit := func(c *commandeer) error {
		c.Set("buildFuture", true)
		return nil
	}
	return h.listContent(cfgInit, (*hugolib.Page).IsFuture)
}

func (h *hugoListCmd) listContent(cb func(c *commandeer) error, isType func(p *hugolib.Page) bool) error {

	err := h.c.InitializeConfig(cb, h.cmd)
	if err != nil {
		return err
	}

	sites, err := hugolib.NewHugoSites(*h.c.HugoBaseConfig)

	if err != nil {
		return core.NewSystemError("Error creating sites", err)
	}

	if err := sites.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
		return core.NewSystemError("Error Processing Source Content", err)
	}

	for _, p := range sites.Pages() {
		if isType(p) {
			jww.FEEDBACK.Println(filepath.Join(p.File.Dir(), p.File.LogicalName()))
		}
	}

	return nil
}
