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

// Build "hugo list" command.
func buildHugoListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listing out various types of content",
		Long: `Listing out various types of content.

List requires a subcommand, e.g. ` + "`hugo list drafts`.",
		RunE: nil,
	}

	// Add hugo list subcommands
	cmd.AddCommand(buildHugoListDraftsCmd())
	cmd.AddCommand(buildHugoListExpiredCmd())
	cmd.AddCommand(buildHugoListFutureCmd())

	// Add flags used by all list subcommands
	cmd.PersistentFlags().StringP("source", "s", "", "filesystem path to read files relative from")

	// Set bash completion
	cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})

	return cmd
}

// Build "hugo list drafts" command.
func buildHugoListDraftsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drafts",
		Short: "List all drafts",
		Long:  `List all of the drafts in your content directory.`,
		RunE:  listDrafts,
	}

	return cmd
}

// Build "hugo list expired" command.
func buildHugoListExpiredCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expired",
		Short: "List all posts already expired",
		Long: `List all of the posts in your content directory which has already
expired.`,
		RunE: listExpired,
	}

	return cmd
}

// Build "hugo list future" command.
func buildHugoListFutureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "future",
		Short: "List all posts dated in the future",
		Long: `List all of the posts in your content directory which will be
posted in the future.`,
		RunE: listFuture,
	}

	return cmd
}

// ----------------------------------------------------------------------------------------------

func listDrafts(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo list drafts - hugo list drafts code goes here")
	return nil
}

func listExpired(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo list expired - hugo list expired code goes here")
	return nil
}

func listFuture(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo list future - hugo list future code goes here")
	return nil
}
