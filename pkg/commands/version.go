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
	"runtime"
	"strings"

	"github.com/neurocline/drouet/pkg/core"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// Build "hugo version" command.
func buildHugoVersionCmd(hugo *core.Hugo) *hugoVersionCmd {
	h := &hugoVersionCmd{Hugo: hugo}

	h.cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's.`,
		RunE:  h.version,
	}

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoVersionCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func (h *hugoVersionCmd) version(cmd *cobra.Command, args []string) error {
	showVersion(h.Hugo.Logger)
	return nil
}

func showVersion(Logger *jww.Notepad) {
	// Create Hugo version string with optional commit hash
	vers := fmt.Sprintf("%s", core.CurrentHugoVersion)
	if core.CommitHash != "" {
		vers = fmt.Sprintf("%s-%s", vers, strings.ToUpper(core.CommitHash))
	}
	os_arch := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	Logger.FEEDBACK.Printf("Hugo Static Site Generator v%s %s BuildDate: %s\n", vers, os_arch, core.BuildDate)
}
