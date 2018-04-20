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
	"net"
	"net/url"
	"strings"

	"github.com/neurocline/cobra"
	"github.com/spf13/viper"
)

// Build "hugo server" command.
func buildHugoServerCmd(hugo *commandeer) *hugoServerCmd {
	h := &hugoServerCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:     "server",
		Aliases: []string{"serve"},
		Short:   "A high performance webserver",
		Long: `Hugo provides its own webserver which builds and serves the site.
While hugo server is high performance, it is a webserver with limited options.
Many run it in production, but the standard behavior is for people to use it
in development and use a more full featured server such as Nginx or Caddy.

'hugo server' will avoid writing the rendered and served content to disk,
preferring to store it in memory.

By default hugo will also watch your files for any changes you make and
automatically rebuild the site. It will then live reload any open browser pages
and push the latest content to them. As most Hugo sites are built in a fraction
of a second, you will be able to save and see your changes nearly instantly.`,
		RunE: h.server,
	}

	// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
	addHugoBuilderFlags(h.cmd)

	// Add flags for "hugo server"
	h.cmd.Flags().BoolVar(&h.serverAppend, "appendPort", true, "append port to baseURL")
	h.cmd.Flags().StringVar(&h.serverInterface, "bind", "127.0.0.1", "interface to which the server will bind")
	h.cmd.Flags().BoolVar(&h.disableFastRender, "disableFastRender", false, "enables full re-renders on changes")
	h.cmd.Flags().BoolVar(&h.disableLiveReload, "disableLiveReload", false, "watch without enabling live browser reload on rebuild")
	h.cmd.Flags().IntVar(&h.liveReloadPort, "liveReloadPort", -1, "port for live reloading (i.e. 443 in HTTPS proxy situations)")
	h.cmd.Flags().String("memstats", "", "log memory usage to this file")
	h.cmd.Flags().String("meminterval", "100ms", "interval to poll memory usage (requires --memstats), valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	h.cmd.Flags().BoolVar(&h.noHTTPCache, "noHTTPCache", false, "prevent HTTP caching")
	h.cmd.Flags().BoolVar(&h.navigateToChanged, "navigateToChanged", false, "navigate to changed content file on live browser reload")
	h.cmd.Flags().IntVarP(&h.serverPort, "port", "p", 1313, "port on which the server will listen")
	h.cmd.Flags().BoolVar(&h.renderToDisk, "renderToDisk", false, "render to Destination path (default is render to memory & serve from there)")
	h.cmd.Flags().BoolVarP(&h.serverWatch, "watch", "w", true, "watch filesystem for changes and recreate as needed")

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoServerCmd struct {
	c *commandeer
	cmd *cobra.Command

	disableFastRender bool
	disableLiveReload bool
	liveReloadPort    int
	navigateToChanged bool
	noHTTPCache       bool
	renderToDisk      bool
	serverAppend      bool
	serverInterface   string
	serverPort        int
	serverWatch       bool

	//visitedURLs *types.EvictingStringQueue
	running bool
}

func (h *hugoServerCmd) server(cmd *cobra.Command, args []string) error {
	// stuff

	cfgInit := func(c *commandeer) error {
		// lots of stuff
		return nil
	}

	err := h.c.InitializeConfig(cfgInit, h.cmd)
	h.c.setRunning(true)

	// more stuff

	return err
}

// fixURL massages the baseURL into a form needed for serving
// all pages correctly.
// Added serverAppend var (it was previously a global variable)
func fixURL(/*cfg config.Provider*/ cfg *viper.Viper, s string, port int, serverAppend bool) (string, error) {
	useLocalhost := false
	if s == "" {
		s = cfg.GetString("baseURL")
		useLocalhost = true
	}

	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}

	// do an initial parse of the input string
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	// if no Host is defined, then assume that no schema or double-slash were
	// present in the url.  Add a double-slash and make a best effort attempt.
	if u.Host == "" && s != "/" {
		s = "//" + s

		u, err = url.Parse(s)
		if err != nil {
			return "", err
		}
	}

	if useLocalhost {
		if u.Scheme == "https" {
			u.Scheme = "http"
		}
		u.Host = "localhost"
	}

	if serverAppend {
		if strings.Contains(u.Host, ":") {
			u.Host, _, err = net.SplitHostPort(u.Host)
			if err != nil {
				return "", fmt.Errorf("Failed to split baseURL hostpost: %s", err)
			}
		}
		u.Host += fmt.Sprintf(":%d", port)
	}

	return u.String(), nil
}
