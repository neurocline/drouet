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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/neurocline/drouet/pkg/hugofs"

	"github.com/neurocline/cobra"
	"github.com/neurocline/viper"

	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
)

// InitializeConfig creates a default config and then updates it with
// values from a config file and from command-line flags.
func (h *Hugo) InitializeConfig(cmd *cobra.Command) error {

	// First, create a default config with Viper
	v := viper.New()

	// Get all the default values
	if err := defaultSettings(v); err != nil {
		return err
	}

	// Then, load any config files we can find (there can be more than one)
	if err := configFileSettings(v, cmd); err != nil {
		return err
	}

	// Then, add any values from command-line flags. We add all of them.
	v.BindPFlags(cmd.Flags())
	v.BindPFlags(cmd.PersistentFlags())

	// Set some overrides to match Hugo behavior (these will probably
	// go away). These are pointless, as far as I can tell, because flags are
	// already the second-highest priority next to overrides.
	// However, I probably need to not set these as strings, hmm? I need to
	// preserve types.
	// Hey, maybe this is overriding locale-specific values? Except lang flattening
	// also uses set, so that can't be true.
	for _, flag := range []string{ "baseURL", "logI18nWarnings", "theme", "themesDir" } {
		if cmd.Flags().Changed(flag) {
			v.Set("themesDir", cmd.Flags().Lookup(flag).Value.String())
		}
	}

	// Make sure cacheDir points to something useful
	var sourceFs afero.Fs = hugofs.Os
	cacheDir := v.GetString("cacheDir")
	if cacheDir != "" {
		if hugofs.FilePathSeparator != cacheDir[len(cacheDir)-1:] {
			cacheDir = cacheDir + hugofs.FilePathSeparator
		}
		isDir, err := hugofs.DirExists(cacheDir, sourceFs)
		CheckErr(h.Log, err)
		if !isDir {
			hugofs.Mkdir(cacheDir)
		}
		v.Set("cacheDir", cacheDir)
	} else {
		v.Set("cacheDir", hugofs.GetTempDir("hugo_cache", sourceFs))
	}

	// Check for deprecated settings being used
	checkDeprecated(v, cmd)

	// All done, remember our good config
	h.Config = v

	// Create logger now that we have config loaded (log config could have been
	// set in a config file or environment variable)
	if err := h.createLogger(); err != nil {
		return err
	}

	return nil
}

func checkDeprecated(v *viper.Viper, cmd *cobra.Command) {

	// useModTimeAsFallback is a deprecated config item
	// scheduled to be removed in Hugo 0.39
	if v.GetBool("useModTimeAsFallback") {
		msg := `Replace with this in your config.toml:

[frontmatter]
date = [ "date",":fileModTime", ":default"]
lastmod = ["lastmod" ,":fileModTime", ":default"]
`
		Deprecated("Site config", "useModTimeAsFallback", msg, false)
	}

	for _, key := range []string{ "uglyURLs", "pluralizeListTitles", "preserveTaxonomyNames", "canonifyURLs"} {
		if cmd.Flags().Changed(key) {
			msg := fmt.Sprintf(`Set "%s = true" in your config.toml.
If you need to set this configuration value from the command line, set it via an OS environment variable: "HUGO_%s=true hugo"`, key, strings.ToUpper(key))
			// Remove in Hugo 0.38
			Deprecated("hugo", "--"+key+" flag", msg, true)
		}
	}
}

func configFileSettings(v *viper.Viper, cmd *cobra.Command) error {

	// If --source was used on the command-line, then this points to the entire
	// site directory and is where the config file is located. Note that --source
	// has no effect if --config is used.
	// TBD hey, Cobra, this dance is awkward. It should really be just
	// basePath, changed = cmd.Flags().Get("source"). Or is there something
	// even simpler? What about the Pythonesque
	// basePath = cmd.Flags.Get("source", os.Getwd())? I like that.
	var basePath string
	if cmd.Flags().Changed("source") {
		basePath, _ = cmd.Flags().GetString("source")
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
	if cmd.Flags().Changed("config") {
		configPath, _ = cmd.Flags().GetString("config")
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
	//v.SetDefault("blackFriday.smartypants", true)
	//v.SetDefault("blackFriday.angledQuotes", false)
	//v.SetDefault("blackFriday.smartypantsQuotesNBSP", false)
	//v.SetDefault("blackFriday.fractions", true)
	//v.SetDefault("blackFriday.hrefTargetBlank", false)
	//v.SetDefault("blackFriday.smartDashes", true)
	//v.SetDefault("blackFriday.latexDashes", true)
	//v.SetDefault("blackFriday.plainIDAnchors", true)
	//v.SetDefault("blackFriday.taskLists", true)

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

// ----------------------------------------------------------------------------------------------

func (h *Hugo) createLogger() error {
	var (
		logHandle       = ioutil.Discard
		logFile         = h.Config.GetString("logFile")
		verboseLog      = h.Config.GetBool("verboseLog")
		logging         = h.Config.GetBool("log")
		isLogging       = false
	)

	// Create a logfile if asked for directly or implicitly
	var err error
	if logFile != "" {
		isLogging = true
		logHandle, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return NewSystemError("Failed to open log file:", logFile, err)
		}
	} else if verboseLog || logging {
		isLogging = true
		logHandle, err = ioutil.TempFile("", "hugo")
		if err != nil {
			return NewSystemError(err)
		}
	}

	// Set the appropriate logging level for console output
	// must do --verbose --debug to get debug console out when logging
	// must do --verbose --trace to get trace console out when logging
	var verbose = h.Config.GetBool("verbose")
	var trace = h.Config.GetBool("trace")
	var debug = h.Config.GetBool("debug")
	var quiet = h.Config.GetBool("quiet")

	var stdoutThreshold = jww.LevelError
	switch {
	case quiet:
	case (!isLogging || verbose) && trace:
		stdoutThreshold = jww.LevelTrace
	case (!isLogging || verbose) && debug:
		stdoutThreshold = jww.LevelDebug
	case verbose:
		stdoutThreshold = jww.LevelInfo
	}

	// Set the appropriate logging level for file output (quiet is ignored here)
	// Note 1: verboseLog should be round-tripped to config so it can be set in config
	// Note 2: log vs logging should be resolved so it can round-trip to config
	var logThreshold = jww.LevelWarn
	switch {
	case verboseLog && trace:
		logThreshold = jww.LevelTrace
	case verboseLog && debug:
		logThreshold = jww.LevelDebug
	case verboseLog:
		logThreshold = jww.LevelInfo
	}

	// The global logger is used in some few cases.
	jww.SetLogOutput(logHandle)
	jww.SetLogThreshold(logThreshold)
	jww.SetStdoutThreshold(stdoutThreshold)
	InitDistinctLoggers()

	h.Log = jww.NewNotepad(stdoutThreshold, logThreshold, os.Stdout, logHandle, "", log.Ldate|log.Ltime)
	return nil
}

// InitLoggers sets up the global distinct loggers.
func InitDistinctLoggers() {
	DistinctErrorLog = NewDistinctErrorLogger()
	DistinctWarnLog = NewDistinctWarnLogger()
	DistinctFeedbackLog = NewDistinctFeedbackLogger()
}

var (
	// DistinctErrorLog can be used to avoid spamming the logs with errors.
	DistinctErrorLog *DistinctLogger = NewDistinctErrorLogger()

	// DistinctWarnLog can be used to avoid spamming the logs with warnings.
	DistinctWarnLog *DistinctLogger = NewDistinctWarnLogger()

	// DistinctFeedbackLog can be used to avoid spamming the logs with info messages.
	DistinctFeedbackLog *DistinctLogger = NewDistinctFeedbackLogger()
)

// NewDistinctErrorLogger creates a new DistinctLogger that logs ERRORs
func NewDistinctErrorLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.ERROR}
}

// NewDistinctWarnLogger creates a new DistinctLogger that logs WARNs
func NewDistinctWarnLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.WARN}
}

// NewDistinctFeedbackLogger creates a new DistinctLogger that can be used
// to give feedback to the user while not spamming with duplicates.
func NewDistinctFeedbackLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.FEEDBACK}
}

// DistinctLogger ignores duplicate log statements.
type DistinctLogger struct {
	sync.RWMutex
	logger logPrinter
	m      map[string]bool
}

type logPrinter interface {
	// Println is the only common method that works in all of JWWs loggers.
	Println(a ...interface{})
}

// Println will log the string returned from fmt.Sprintln given the arguments,
// but not if it has been logged before.
func (l *DistinctLogger) Println(v ...interface{}) {
	// fmt.Sprint doesn't add space between string arguments
	logStatement := strings.TrimSpace(fmt.Sprintln(v...))
	l.print(logStatement)
}

// Printf will log the string returned from fmt.Sprintf given the arguments,
// but not if it has been logged before.
// Note: A newline is appended.
func (l *DistinctLogger) Printf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.print(logStatement)
}

func (l *DistinctLogger) print(logStatement string) {
	l.RLock()
	if l.m[logStatement] {
		l.RUnlock()
		return
	}
	l.RUnlock()

	l.Lock()
	if !l.m[logStatement] {
		l.logger.Println(logStatement)
		l.m[logStatement] = true
	}
	l.Unlock()
}
