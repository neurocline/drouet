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
	"github.com/neurocline/drouet/pkg/core"
	"github.com/neurocline/drouet/pkg/parser"

	"github.com/spf13/cobra"
)

// Build "hugo convert" command.
func buildHugoConvertCmd(hugo *core.Hugo) *hugoConvertCmd {
	h := &hugoConvertCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert your content to different formats",
		Long: `Convert your content (e.g. front matter) to different formats.

See convert's subcommands toJSON, toTOML and toYAML for more information.`,
		RunE: nil, // will print usage
	}

	h.cmd.AddCommand(buildHugoConvertToJsonCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoConvertToTomlCmd(hugo).cmd)
	h.cmd.AddCommand(buildHugoConvertToYamlCmd(hugo).cmd)

	// Flags shared between all convert commands
	h.cmd.PersistentFlags().StringVarP(&h.outputDir, "output", "o", "", "filesystem path to write files to")
	h.cmd.PersistentFlags().StringVarP(&h.source, "source", "s", "", "filesystem path to read files relative from")
	h.cmd.PersistentFlags().BoolVar(&h.unsafe, "unsafe", false, "enable less safe operations, please backup first")

	h.cmd.PersistentFlags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})

	return h
}

func buildHugoConvertToJsonCmd(hugo *core.Hugo) *hugoConvertCmd {
	h := &hugoConvertCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "toJSON",
		Short: "Convert front matter to JSON",
		Long: `toJSON converts all front matter in the content directory
to use JSON for the front matter.`,
		RunE: h.convertToJson,
	}

	return h
}

func buildHugoConvertToTomlCmd(hugo *core.Hugo) *hugoConvertCmd {
	h := &hugoConvertCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "toTOML",
		Short: "Convert front matter to TOML",
		Long: `toTOML converts all front matter in the content directory
to use TOML for the front matter.`,
		RunE: h.convertToToml,
	}

	return h
}

func buildHugoConvertToYamlCmd(hugo *core.Hugo) *hugoConvertCmd {
	h := &hugoConvertCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "toYAML",
		Short: "Convert front matter to YAML",
		Long: `toYAML converts all front matter in the content directory
to use YAML for the front matter.`,
		RunE: h.convertToYaml,
	}

	return h
}

// ----------------------------------------------------------------------------------------------

// All of the "hugo convert" sub-commands share the same set of flags and so
// share the same command structure
type hugoConvertCmd struct {
	*core.Hugo
	cmd *cobra.Command

	outputDir string
	source string
	unsafe bool
}

func (h *hugoConvertCmd) convertToJson(cmd *cobra.Command, args []string) error {
	return h.convertContents(rune([]byte(parser.JSONLead)[0]))
}

func (h *hugoConvertCmd) convertToToml(cmd *cobra.Command, args []string) error {
	return h.convertContents(rune([]byte(parser.TOMLLead)[0]))
}

func (h *hugoConvertCmd) convertToYaml(cmd *cobra.Command, args []string) error {
	return h.convertContents(rune([]byte(parser.YAMLLead)[0]))
}

func (h *hugoConvertCmd) convertContents(mark rune) error {
	if h.outputDir == "" && !h.unsafe {
		return core.NewUserError("Unsafe operation not allowed, use --unsafe or set a different output path")
	}

	if err := h.Hugo.InitializeConfig(h.cmd); err != nil {
		return err
	}

//	h, err := hugolib.NewHugoSites(*c.DepsCfg)
//	if err != nil {
//		return err
//	}
//
//	if err := h.Build(hugolib.BuildCfg{SkipRender: true}); err != nil {
//		return err
//	}
//
//	site := h.Sites[0]
//
//	site.Log.FEEDBACK.Println("processing", len(site.AllPages), "content files")
//	for _, p := range site.AllPages {
//		if err := convertAndSavePage(p, site, mark); err != nil {
//			return err
//		}
//	}
	return nil
}
