// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"io/ioutil"
	"log"
	"os"
	//"path/filepath"
	"sync"
	"time"

	"github.com/neurocline/drouet/pkg/core"

	//"github.com/neurocline/drouet/pkg/common/types"
	//"github.com/neurocline/drouet/pkg/deps"
	//"github.com/neurocline/drouet/pkg/helpers"
	//"github.com/neurocline/drouet/pkg/hugofs"
	//src "github.com/neurocline/drouet/pkg/source"
	//"github.com/neurocline/drouet/pkg/hugolib"
	//"github.com/neurocline/drouet/pkg/utils"

	"github.com/neurocline/cobra"

	"github.com/bep/debounce"
	//"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
)

type commandeer struct {
	*core.HugoBaseConfig

	subCmdVs []*cobra.Command

//	pathSpec    *helpers.PathSpec
	visitedURLs *core.EvictingStringQueue

//	staticDirsConfig []*src.Dirs

	// We watch these for changes.
//	configFiles []string

	// Called each time InitializeConfig/c.loadConfig is executed
	configCallback func(c *commandeer) error

	// We can do this only once.
	fsCreate sync.Once

	// Used in cases where we get flooded with events in server mode.
	debounce func(f func())

	serverPorts []int
//	languages   helpers.Languages

	configured bool
}

// Create our single instance of our command processor controller
func newCommandeer() *commandeer {

	return &commandeer{
		HugoBaseConfig: &core.HugoBaseConfig{},
		visitedURLs:    core.NewEvictingStringQueue(10),
	}
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func (c *commandeer) InitializeConfig(inject func(c *commandeer) error, subCmdVs ...*cobra.Command) error {

	c.configCallback = inject
	c.subCmdVs = subCmdVs

	return c.loadConfig()
}

func (c *commandeer) setRunning(running bool) {
	if running {
		// The time value used is tested with mass content replacements in a fairly big Hugo site.
		// It is better to wait for some seconds in those cases rather than get flooded
		// with rebuilds.
		c.debounce, _ = debounce.New(4 * time.Second)
	}
	c.HugoBaseConfig.Running = running
}

func (c *commandeer) loadConfig() error {

	cfg := c.HugoBaseConfig
	c.configured = false

	var dir string
	if source != "" {
		dir, _ = filepath.Abs(source)
	} else {
		dir, _ = os.Getwd()
	}

	var sourceFs afero.Fs = hugofs.Os
	if c.HugoBaseConfig.Fs != nil {
		sourceFs = c.HugoBaseConfig.Fs.Source
	}

	config, configFiles, err := hugolib.LoadConfig(hugolib.ConfigSourceDescriptor{Fs: sourceFs, Path: source, WorkingDir: dir, Filename: cfgFile})
	if err != nil {
		return err
	}

	c.Cfg = config
	c.configFiles = configFiles

	for _, cmdV := range c.subCmdVs {
		c.initializeFlags(cmdV)
	}

	if l, ok := c.Cfg.Get("languagesSorted").(helpers.Languages); ok {
		c.languages = l
	}

	if baseURL != "" {
		config.Set("baseURL", baseURL)
	}

	if c.configCallback != nil {
		err = c.configCallback(c)
	}

	if err != nil {
		return err
	}

	if len(disableKinds) > 0 {
		c.Set("disableKinds", disableKinds)
	}

	logger, err := createLogger(cfg.Cfg)
	if err != nil {
		return err
	}

	cfg.Logger = logger

	config.Set("logI18nWarnings", logI18nWarnings)

	if theme != "" {
		config.Set("theme", theme)
	}

	if themesDir != "" {
		config.Set("themesDir", themesDir)
	}

	if destination != "" {
		config.Set("publishDir", destination)
	}

	config.Set("workingDir", dir)

	if contentDir != "" {
		config.Set("contentDir", contentDir)
	}

	if layoutDir != "" {
		config.Set("layoutDir", layoutDir)
	}

	if cacheDir != "" {
		config.Set("cacheDir", cacheDir)
	}

	createMemFs := config.GetBool("renderToMemory")

	if createMemFs {
		// Rendering to memoryFS, publish to Root regardless of publishDir.
		config.Set("publishDir", "/")
	}

	c.fsCreate.Do(func() {
		fs := hugofs.NewFrom(sourceFs, config)

		// Hugo writes the output to memory instead of the disk.
		if createMemFs {
			fs.Destination = new(afero.MemMapFs)
		}

		err = c.initFs(fs)
	})

	if err != nil {
		return err
	}

	cacheDir = config.GetString("cacheDir")
	if cacheDir != "" {
		if hugofs.FilePathSeparator != cacheDir[len(cacheDir)-1:] {
			cacheDir = cacheDir + hugofs.FilePathSeparator
		}
		isDir, err := helpers.DirExists(cacheDir, sourceFs)
		utils.CheckErr(cfg.Logger, err)
		if !isDir {
			mkdir(cacheDir)
		}
		config.Set("cacheDir", cacheDir)
	} else {
		config.Set("cacheDir", helpers.GetTempDir("hugo_cache", sourceFs))
	}

	cfg.Logger.INFO.Println("Using config file:", config.ConfigFileUsed())

	themeDir := c.PathSpec().GetThemeDir()
	if themeDir != "" {
		if _, err := sourceFs.Stat(themeDir); os.IsNotExist(err) {
			return newSystemError("Unable to find theme Directory:", themeDir)
		}
	}

	themeVersionMismatch, minVersion := c.isThemeVsHugoVersionMismatch(sourceFs)

	if themeVersionMismatch {
		cfg.Logger.ERROR.Printf("Current theme does not support Hugo version %s. Minimum version required is %s\n",
			helpers.CurrentHugoVersion.ReleaseVersion(), minVersion)
	}

	return nil
}

func (c *commandeer) Set(key string, value interface{}) {
	if c.configured {
		panic("commandeer cannot be changed")
	}
	c.Cfg.Set(key, value)
}

func (c *commandeer) createLogger() error {
	var (
		logHandle       = ioutil.Discard
		logFile         = c.Cfg.GetString("logFile")
		verboseLog      = c.Cfg.GetBool("verboseLog")
		logging         = c.Cfg.GetBool("log")
		isLogging       = false
	)

	// Create a logfile if asked for directly or implicitly
	var err error
	if logFile != "" {
		isLogging = true
		logHandle, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return core.NewSystemError("Failed to open log file:", logFile, err)
		}
	} else if verboseLog || logging {
		isLogging = true
		logHandle, err = ioutil.TempFile("", "hugo")
		if err != nil {
			return core.NewSystemError(err)
		}
	}

	// Set the appropriate logging level for console output
	// must do --verbose --debug to get debug console out when logging
	// must do --verbose --trace to get trace console out when logging
	var verbose = c.Cfg.GetBool("verbose")
	var trace = c.Cfg.GetBool("trace")
	var debug = c.Cfg.GetBool("debug")
	var quiet = c.Cfg.GetBool("quiet")

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
	core.InitDistinctLoggers()

	c.Logger = jww.NewNotepad(stdoutThreshold, logThreshold, os.Stdout, logHandle, "", log.Ldate|log.Ltime)
	return nil
}
