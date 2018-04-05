// Copyright 2016 The Hugo Authors. All rights reserved.
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

// Package commands handles Hugo command-line processing.
//
// Hugo commands and flags are implemented using Cobra.
package commands

import (
	"fmt"

	"github.com/neurocline/drouet/pkg/core"

	"github.com/spf13/cobra"
	"github.com/spf13/nitro"
	"github.com/spf13/pflag"
)

// Execute builds a command processor and runs the user command.
func Execute() int {

	// Do basic system init
	core.GlobalInit()

	// Build our root object and all command processors
	hugo := core.NewHugo()
	root := buildCommand(hugo)

	// Run the command processor; show usage message if there is an error
	// (TBD clean this a bit, full usage is very long for some commands)
	if c, err := root.cmd.ExecuteC(); err != nil {
		if core.IsUserError(err) {
			c.Println("")
			c.Println(c.UsageString())
		}
		return -1
	}

	// Shut everything down as cleanly as possible
	return hugo.Shutdown()
}

// Build the Hugo command - root and all its children
// (every other command (verb) is attached as a child command)
func buildCommand(hugo *core.Hugo) *hugoCmd {

	// Create a new Hugo object and create the root "hugo" command
	gohugo := buildHugoCommand(hugo)
	cmd := gohugo.cmd

	// Add all the sub-commands (sub-commands of sub-commands will add
	// their own children)
	cmd.AddCommand(buildHugoBenchmarkCmd(hugo).cmd)
	cmd.AddCommand(buildHugoCheckCmd(hugo).cmd)
	cmd.AddCommand(buildHugoConfigCmd(hugo).cmd)
	cmd.AddCommand(buildHugoConvertCmd(hugo).cmd)
	cmd.AddCommand(buildHugoEnvCmd(hugo).cmd)
	cmd.AddCommand(buildHugoGenCmd(hugo).cmd)
	cmd.AddCommand(buildHugoImportCmd(hugo).cmd)
	cmd.AddCommand(buildHugoListCmd(hugo).cmd)
	cmd.AddCommand(buildHugoNewCmd(hugo).cmd)
	cmd.AddCommand(buildHugoReleaseCmd(hugo).cmd)
	cmd.AddCommand(buildHugoServerCmd(hugo).cmd)
	cmd.AddCommand(buildHugoVersionCmd(hugo).cmd)

	// Add global flags apply to all commands
	cmd.PersistentFlags().Bool("debug", false, "debug output")
	cmd.PersistentFlags().String("config", "", "config file (default: ./config.(yaml|json|toml))")
	cmd.PersistentFlags().Bool("log", false, "enable Logging")
	cmd.PersistentFlags().String("logFile", "", "log File path (if set, also enable Logging)")
	cmd.PersistentFlags().Bool("quiet", false, "build in quiet mode")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	cmd.PersistentFlags().Bool("verboseLog", false, "verbose logging")

	// Set bash-completion on a few of our global flags
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	cmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
	cmd.PersistentFlags().SetAnnotation("logFile", cobra.BashCompFilenameExt, []string{})

	// Rewrite flags to follow standards (e.g. change --baseUrl into --baseURL)
	cmd.SetGlobalNormalizationFunc(normalizeHugoFlags)

	// We don't want usage spit out all the time (but we end up doing this
	// ourselves, so I'm not sure exactly what this is for).
	// Note - this is automatically inherited, so setting it on the root
	// command means all sub-commands are silenced too.
	cmd.SilenceUsage = true

	return gohugo
}

// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
func addHugoBuilderFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("baseURL", "b", "", "hostname (and path) to the root, e.g. http://spf13.com/")
	cmd.Flags().BoolP("buildDrafts", "D", false, "include content marked as draft")
	cmd.Flags().BoolP("buildExpired", "E", false, "include expired content")
	cmd.Flags().BoolP("buildFuture", "F", false, "include content with publishdate in the future")
	cmd.Flags().String("cacheDir", "", "filesystem path to cache directory (default: $TMPDIR/hugo_cache/)")
	cmd.Flags().Bool("canonifyURLs", false, "(deprecated) if true, all relative URLs canonicalized against baseURL")
	cmd.Flags().Bool("cleanDestinationDir", false, "before build, remove files from destination not found in static directories")
	cmd.Flags().StringP("contentDir", "c", "", "filesystem path to content directory")
	cmd.Flags().StringP("destination", "d", "", "filesystem path to write files to")
	cmd.Flags().StringSlice("disableKinds", []string{}, "disable different kind of pages (home, RSS etc.)")
	cmd.Flags().Bool("enableGitInfo", false, "add Git revision, date and author info to the pages")
	cmd.Flags().BoolP("forceSyncStatic", "", false, "copy all files when static is changed.")
	cmd.Flags().Bool("gc", false, "if true, run cleanup tasks (like 'remove unused cache files) after the build")
	cmd.Flags().BoolP("i18n-warnings", "", false, "print missing translations")
	cmd.Flags().Bool("ignoreCache", false, "ignores the cache directory")
	cmd.Flags().StringP("layoutDir", "l", "", "filesystem path to layout directory")
	cmd.Flags().BoolP("noChmod", "", false, "don't sync permission mode of files")
	cmd.Flags().BoolP("noTimes", "", false, "don't sync modification time of files")
	cmd.Flags().Bool("pluralizeListTitles", true, "(deprecated) pluralize titles in lists using inflect")
	cmd.Flags().Bool("preserveTaxonomyNames", false, `(deprecated) preserve taxonomy names as written ("Gérard Depardieu" vs "gerard-depardieu")`)
	cmd.Flags().StringP("source", "s", "", "filesystem path to read files relative from")
	cmd.Flags().Bool("templateMetrics", false, "display metrics about template executions")
	cmd.Flags().Bool("templateMetricsHints", false, "calculate some improvement hints when combined with --templateMetrics")
	cmd.Flags().StringP("theme", "t", "", "theme to use (located in /themes/THEMENAME/)")
	cmd.Flags().String("themesDir", "", "filesystem path to themes directory")
	cmd.Flags().Bool("uglyURLs", false, "(deprecated) if true, use /filename.html instead of /filename/")

	// This is a global in a package not under my control, so I'm leaving it as a global write,
	// even though it could be wrapped more nicely. And maybe whole-application performance
	// monitoring means global state isn't a flaw but a feature.
	cmd.Flags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")

	// Set bash-completion.
	// Each flag must first be defined before using the SetAnnotation() call.
	cmd.Flags().SetAnnotation("cacheDir", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("destination", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("theme", cobra.BashCompSubdirsInDir, []string{"themes"})
}

// normalizeHugoFlags facilitates transitions of Hugo command-line flags,
// e.g. --baseUrl to --baseURL, --uglyUrls to --uglyURLs
func normalizeHugoFlags(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "baseUrl":
		name = "baseURL"
		break
	case "uglyUrls":
		name = "uglyURLs"
		break
	}
	return pflag.NormalizedName(name)
}

// ----------------------------------------------------------------------------------------------

// Build Hugo root command.
func buildHugoCommand(hugo *core.Hugo) *hugoCmd {
	h := &hugoCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "hugo",
		Short: "hugo builds your site",
		Long: `hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io/.`,
		RunE: h.hugo,
	}

	// Add flags for the "hugo" command
	h.cmd.Flags().BoolVar(&h.renderToMemory, "renderToMemory", false, "render to memory (useful for benchmark testing)")
	h.cmd.Flags().BoolVarP(&h.buildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")

	// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
	addHugoBuilderFlags(h.cmd)

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoCmd struct {
	*core.Hugo
	cmd *cobra.Command

	renderToMemory bool
	buildWatch bool

	//visitedURLs *types.EvictingStringQueue
	running bool
}

func (h *hugoCmd) hugo(cmd *cobra.Command, args []string) error {
	//h.visitedURLs = types.NewEvictingStringQueue(10)
	h.running = h.buildWatch

	fmt.Println("hugo - hugo (build) code goes here")
	return nil
}

// ----------------------------------------------------------------------------------------------

func (h *hugoCmd) Reset() {
	//h.Hugo.HugoSites = nil
}

func (h *hugoCmd) resetAndBuildSites() (err error) {
	if err = h.initSites(); err != nil {
		return
	}
	if !h.Config.GetBool("quiet") {
		h.Logger.FEEDBACK.Println("Started building sites ...")
	}
	//return Hugo.Build(hugolib.BuildCfg{ResetState: true})
	return nil
}

func (h *hugoCmd) initSites() error {
	return nil
}
