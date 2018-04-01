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

// Wrap the core.Hugo type so we can declare methods that take a *core.Hugo receiver
type hugoCmd struct {
	*core.Hugo
}

// Execute builds a command processor and runs the user command.
func Execute() int {

	// Do basic system init
	core.Init()

	// Create Hugo object to hold Hugo state, and build command processor
	// (we pass the Hugo object in to every subcommand processor so each
	// subcommand has access to the Hugo object).
	hugo := &hugoCmd{core.NewHugo()}
	cmd := buildCommand(hugo)

	// Run the command processor; show usage message if there is an error
	// (TBD clean this a bit, full usage is very long for some commands)
	if c, err := cmd.ExecuteC(); err != nil {
		c.Println("")
		c.Println(c.UsageString())
		return -1
	}

	// Shut everything down as cleanly as possible
	return hugo.Shutdown()
}

// Build the Hugo command - root and all its children
// (every other command (verb) is attached as a child command)
func buildCommand(h *hugoCmd) *cobra.Command {
	root := buildHugoCommand(h)

	root.AddCommand(buildHugoBenchmarkCmd(h))
	root.AddCommand(buildHugoCheckCmd(h))
	root.AddCommand(buildHugoConfigCmd(h))
	root.AddCommand(buildHugoConvertCmd(h))
	root.AddCommand(buildHugoEnvCmd(h))
	root.AddCommand(buildHugoGenCmd(h))
	root.AddCommand(buildHugoImportCmd(h))
	root.AddCommand(buildHugoListCmd(h))
	root.AddCommand(buildHugoNewCmd(h))
	root.AddCommand(buildHugoReleaseCmd(h).cmd)
	root.AddCommand(buildHugoServerCmd(h))
	root.AddCommand(buildHugoVersionCmd(h))
	return root
}

// Build Hugo root command.
func buildHugoCommand(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hugo",
		Short: "hugo builds your site",
		Long: `hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io/.`,
		RunE: h.hugo,
	}

	// Global flags apply to all commands
	cmd.PersistentFlags().Bool("debug", false, "debug output")
	cmd.PersistentFlags().String("config", "", "config file (default: ./config.(yaml|json|toml))")
	cmd.PersistentFlags().Bool("log", false, "enable Logging")
	cmd.PersistentFlags().String("logFile", "", "log File path (if set, also enable Logging)")
	cmd.PersistentFlags().Bool("quiet", false, "build in quiet mode")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	cmd.PersistentFlags().Bool("verboseLog", false, "verbose logging")

	// Set bash-completion
	// Each flag must first be defined before using the SetAnnotation() call.
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	cmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
	cmd.PersistentFlags().SetAnnotation("logFile", cobra.BashCompFilenameExt, []string{})

	// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
	initHugoBuilderFlags(cmd)

	// Add flags shared by benchmarking: "hugo", "hugo benchmark"
	initHugoBenchmarkFlags(cmd)

	// Add flags unique to the "hugo" command
	cmd.Flags().BoolP("watch", "w", false, "watch filesystem for changes and recreate as needed")

	// Update flags to latest convention
	// (TBD maybe we should be setting more aliases?)
	cmd.SetGlobalNormalizationFunc(normalizeHugoFlags)

	// We don't want usage spit out all the time (but we end up doing this
	// ourselves, so I'm not sure exactly what this is for)
	cmd.SilenceUsage = true

	//hugoCmdV = cmd
	return cmd
}

// ----------------------------------------------------------------------------------------------

// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
func initHugoBuilderFlags(cmd *cobra.Command) {
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

// Add flags shared by benchmarking: "hugo", "hugo benchmark"
func initHugoBenchmarkFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("renderToMemory", false, "render to memory (only useful for benchmark testing)")
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

// "hugo" with no verb is "hugo build", build a site
func (h *hugoCmd) hugo(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo - build site goes here")
	fmt.Printf("Hugo: %+v\n", *h.Hugo)
	return nil
}
