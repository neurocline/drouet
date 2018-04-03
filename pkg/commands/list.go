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

// Build "hugo list" command.
func buildHugoListCmd(hugo *core.Hugo) *hugoListCmd {
	h := &hugoListCmd{Hugo: hugo}

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
func buildHugoListDraftsCmd(hugo *core.Hugo) *hugoListDraftsCmd {
	h := &hugoListDraftsCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "drafts",
		Short: "List all drafts",
		Long:  `List all of the drafts in your content directory.`,
		RunE:  h.listDrafts,
	}

	return h
}

// Build "hugo list expired" command.
func buildHugoListExpiredCmd(hugo *core.Hugo) *hugoListExpiredCmd {
	h := &hugoListExpiredCmd{Hugo: hugo}

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
func buildHugoListFutureCmd(hugo *core.Hugo) *hugoListFutureCmd {
	h := &hugoListFutureCmd{Hugo: hugo}

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
	*core.Hugo
	cmd *cobra.Command
}

type hugoListDraftsCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoListDraftsCmd) listDrafts(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo list drafts - hugo list drafts code goes here")
	return nil
}

type hugoListExpiredCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoListExpiredCmd) listExpired(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo list expired - hugo list expired code goes here")
	return nil
}

type hugoListFutureCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoListFutureCmd) listFuture(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo list future - hugo list future code goes here")
	return nil
}
