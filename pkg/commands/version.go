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
func buildHugoVersionCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's.`,
		RunE:  h.version,
	}

	return cmd
}

// ----------------------------------------------------------------------------------------------

func (h *hugoCmd) version(cmd *cobra.Command, args []string) error {

	os_arch_date := fmt.Sprintf("%s/%s BuildDate: %s", runtime.GOOS, runtime.GOARCH, core.BuildDate)
	vers := fmt.Sprintf("%s", core.CurrentHugoVersion())
	if core.CommitHash != "" {
		vers = fmt.Sprintf("%s-%s", vers, strings.ToUpper(core.CommitHash))
	}

	jww.FEEDBACK.Printf("Hugo Static Site Generator v%s %s\n", vers, os_arch_date)

	//if hugolib.CommitHash == "" {
	//	jww.FEEDBACK.Printf("Hugo Static Site Generator v%s %s/%s BuildDate: %s\n", helpers.CurrentHugoVersion, runtime.GOOS, runtime.GOARCH, hugolib.BuildDate)
	//} else {
	//	jww.FEEDBACK.Printf("Hugo Static Site Generator v%s-%s %s/%s BuildDate: %s\n", helpers.CurrentHugoVersion, strings.ToUpper(hugolib.CommitHash), runtime.GOOS, runtime.GOARCH, hugolib.BuildDate)
	//}

	return nil
}
