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

// Caveat with Viper - merging between config files or between
// levels of the Viper sources hierarchy works for individual
// values, but not maps or slices. This is not a concern for the
// built-in Hugo config (which is all individual values), but it
// should be noted somewhere, and perhaps a verbose warning that
// merging threw data away

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/neurocline/viper"
	"github.com/spf13/cobra"

	"github.com/neurocline/drouet/pkg/z"
)

// InitializeConfig creates a default config and then updates it with
// values from a config file and from command-line flags.
func InitializeConfig(h *Hugo, cmds ...*cobra.Command) (*viper.Viper, error) {
	fmt.Fprintf(z.Log, "core.InitializeConfig\n%s\n", z.Stack())

	// First, create a default config with Viper
	v := viper.New()
	//fmt.Printf("------------\nviper.New()\n----\n%s", v.Spew())

	// Get all the default values
	if err := defaultSettings(v); err != nil {
		return nil, err
	}
	//fmt.Printf("------------\ndefaultSettings\n----\n%s", v.Spew())

	// Then, load any config files we can find (there can be more than one)
	if err := configFileSettings(v, cmds...); err != nil {
		return nil, err
	}
	//fmt.Printf("------------\nconfigFileSettings\n----\n%s", v.Spew())

	// Then, add any values from command-line flags. We add all of them
	// (TBD we need some aliasing here)
	for _, cmd := range cmds {
		v.BindPFlags(cmd.Flags())
		v.BindPFlags(cmd.PersistentFlags())
	}
	//fmt.Printf("------------\nadd flags\n----\n%s", v.Spew())

	// Make sure cacheDir points to something useful
//	if v.Get("cacheDir") == "" {
//		v.Set("cacheDir", helpers.GetTempDir("hugo_cache", sourceFs))
//	} else {
//
//	}

/*	// Check for deprecated settings being used
	// useModTimeAsFallback is a deprecated config item
	// scheduled to be removed in Hugo 0.39
	if v.GetBool("useModTimeAsFallback") {
		helpers.Deprecated("Site config", "useModTimeAsFallback", `Replace with this in your config.toml:

[frontmatter]
date = [ "date",":fileModTime", ":default"]
lastmod = ["lastmod" ,":fileModTime", ":default"]
`, false)

	}*/

	//fmt.Printf("------------\nfinal state\n----\n%s", v.Spew())
	fmt.Fprintf(z.Log, "---------\nViper config\n---------\n%s", v.Spew())

	h.Config = v
	return v, nil
}

func configFileSettings(v *viper.Viper, cmds ...*cobra.Command) error {

	// If --source was used on the command-line, then this points to the entire
	// site directory and is where the config file is located. Note that --source
	// has no effect if --config is used.
	// TBD hey, Cobra, this dance is awkward. It should really be just
	// basePath, changed = cmds[0].Flags().Get("source"). Or is there something
	// even simpler? What about the Pythonesque
	// basePath = cmds[0].Flags.Get("source", os.Getwd())? I like that.
	var basePath string
	if cmds[0].Flags().Changed("source") {
		basePath, _ = cmds[0].Flags().GetString("source")
		basePath, _ = filepath.Abs(basePath)
	} else {
		basePath, _ = os.Getwd()
	}

	// If --config was used on the command-line, then this points to
	// one or more config files, comma-separated. Due to the behavior
	// of strings.Split on an empty string, we will always have at
	// least one entry in the slice, so we will always have at least one
	// entry to pass to ReadInConfig
	var configPath string
	if cmds[0].Flags().Changed("config") {
		configPath, _ = cmds[0].Flags().GetString("config")
	}
	configs := strings.Split(configPath, ",")

	v.AutomaticEnv()
	v.SetEnvPrefix("hugo")

	// If v.configFile is a non-empty string, it will be used as the path to the config file.
	// Otherwise, we search in v.configPaths for the default name followed by a likely extension.
	v.SetConfigFile(configs[0])
	v.AddConfigPath(basePath)

	// Read in the first config file. We must have at least one.
	var configFiles []string
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return err
		}
		return fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site.\n       Run `hugo help new` for details. (%s)\n", err)
	}

	// Remember the path to the config file we just loaded
	if cf := v.ConfigFileUsed(); cf != "" {
		configFiles = append(configFiles, cf)
	}

	// If there are more config files, add them as well
	for _, configFile := range configs[1:] {
		v.SetConfigFile(configFile)
		if err = v.MergeInConfig(); err != nil {
			return fmt.Errorf("Unable to parse/merge Config file (%s).\n (%s)\n", configFile, err)
		}
		configFiles = append(configFiles, configFile)
	}

	// Memoize basePath as workingDir
	v.Set("workingDir", basePath)

	return nil
}

func defaultSettings(v *viper.Viper) error {

	// Make 'taxonomies' be available via the alias 'indexes'
	// (TBD when can this be removed? Old term, right?)
	v.RegisterAlias("indexes", "taxonomies")

	// Register simple defaults
	v.SetDefault("archetypeDir", "archetypes")
	v.SetDefault("buildDrafts", false)
	v.SetDefault("buildExpired", false)
	v.SetDefault("buildFuture", false)
	v.SetDefault("canonifyURLs", false)
	v.SetDefault("cleanDestinationDir", false)
	v.SetDefault("contentDir", "content")
	v.SetDefault("defaultContentLanguage", "en")
	v.SetDefault("defaultContentLanguageInSubdir", false)
	v.SetDefault("enableMissingTranslationPlaceholders", false)
	v.SetDefault("enableGitInfo", false)
	v.SetDefault("ignoreFiles", make([]string, 0))
	v.SetDefault("dataDir", "data")
	v.SetDefault("debug", false)
	v.SetDefault("disableAliases", false)
	v.SetDefault("disableFastRender", false)
	v.SetDefault("disableLiveReload", false)
	v.SetDefault("disablePathToLower", false)
	v.SetDefault("enableEmoji", false)
	v.SetDefault("footnoteAnchorPrefix", "")
	v.SetDefault("footnoteReturnLinkContents", "")
	v.SetDefault("forceSyncStatic", false)
	v.SetDefault("hasCJKLanguage", false)
	v.SetDefault("i18nDir", "i18n")
	v.SetDefault("ignoreCache", false)
	v.SetDefault("layoutDir", "layouts")
	v.SetDefault("metaDataFormat", "toml")
	v.SetDefault("newContentEditor", "")
	v.SetDefault("paginate", 10)
	v.SetDefault("paginatePath", "page")
	// v.SetDefault("permalinks", make(PermalinkOverrides, 0))
	v.SetDefault("pluralizeListTitles", true)
	v.SetDefault("preserveTaxonomyNames", false)
	v.SetDefault("publishDir", "public")
	v.SetDefault("relativeURLs", false)
	v.SetDefault("removePathAccents", false)
	v.SetDefault("resourceDir", "resources")
	v.SetDefault("rssLimit", -1)
	v.SetDefault("rSSUri", "index.xml")
	v.SetDefault("sectionPagesMenu", "")
	// v.SetDefault("sitemap", Sitemap{Priority: -1, Filename: "sitemap.xml"})
	v.SetDefault("staticDir", "static")
	v.SetDefault("summaryLength", 70)
	v.SetDefault("taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	v.SetDefault("themesDir", "themes")
	v.SetDefault("titleCaseStyle", "AP")
	v.SetDefault("uglyURLs", false)
	v.SetDefault("useModTimeAsFallback", false)
	v.SetDefault("verbose", false)
	v.SetDefault("watch", false)

	// TBD put this back in with BlackFriday support code
	v.SetDefault("blackFriday.smartypants", true)
	v.SetDefault("blackFriday.angledQuotes", false)
	v.SetDefault("blackFriday.smartypantsQuotesNBSP", false)
	v.SetDefault("blackFriday.fractions", true)
	v.SetDefault("blackFriday.hrefTargetBlank", false)
	v.SetDefault("blackFriday.smartDashes", true)
	v.SetDefault("blackFriday.latexDashes", true)
	v.SetDefault("blackFriday.plainIDAnchors", true)
	v.SetDefault("blackFriday.taskLists", true)

	// TBD put this in with Pygments support code
	// TBD if we don't use Pygments, should we set its defaults?
	v.SetDefault("pygmentsCodeFences", false)
	v.SetDefault("pygmentsCodeFencesGuessSyntax", false)
	v.SetDefault("pygmentsOptions", "")
	v.SetDefault("pygmentsStyle", "monokai")
	v.SetDefault("pygmentsUseClasses", false)
	v.SetDefault("pygmentsUseClassic", false)

	return nil
}
