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

// Build "hugo convert" command.
func buildHugoConvertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert your content to different formats",
		Long: `Convert your content (e.g. front matter) to different formats.

See convert's subcommands toJSON, toTOML and toYAML for more information.`,
		RunE: nil, // will print usage
	}

	cmd.AddCommand(buildHugoConvertToJsonCmd())
	cmd.AddCommand(buildHugoConvertToTomlCmd())
	cmd.AddCommand(buildHugoConvertToYamlCmd())

	// Flags shared between all convert commands
	cmd.PersistentFlags().StringP("output", "o", "", "filesystem path to write files to")
	cmd.PersistentFlags().StringP("source", "s", "", "filesystem path to read files relative from")
	cmd.PersistentFlags().Bool("unsafe", false, "enable less safe operations, please backup first")

	cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})

	return cmd
}

func buildHugoConvertToJsonCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "toJSON",
		Short: "Convert front matter to JSON",
		Long: `toJSON converts all front matter in the content directory
to use JSON for the front matter.`,
		RunE: convertToJson,
	}

	return cmd
}

func buildHugoConvertToTomlCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "toTOML",
		Short: "Convert front matter to TOML",
		Long: `toTOML converts all front matter in the content directory
to use TOML for the front matter.`,
		RunE: convertToToml,
	}

	return cmd
}

func buildHugoConvertToYamlCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "toYAML",
		Short: "Convert front matter to YAML",
		Long: `toYAML converts all front matter in the content directory
to use YAML for the front matter.`,
		RunE: convertToYaml,
	}

	return cmd
}

func convert(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo convert - hugo convert goes here")
	return nil
}

func convertToJson(cmd *cobra.Command, args []string) error {
	// return convertContents(rune([]byte(parser.JSONLead)[0]))
	fmt.Println("hugo convert json - hugo convert json goes here")
	return nil
}

func convertToToml(cmd *cobra.Command, args []string) error {
	// return convertContents(rune([]byte(parser.TOMLLead)[0]))
	fmt.Println("hugo convert json - hugo convert json goes here")
	return nil
}

func convertToYaml(cmd *cobra.Command, args []string) error {
	// return convertContents(rune([]byte(parser.YAMLLead)[0]))
	fmt.Println("hugo convert json - hugo convert json goes here")
	return nil
}
