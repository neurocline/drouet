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
	"errors"

	"github.com/neurocline/drouet/pkg/releaser"

	"github.com/neurocline/cobra"
)

// Build "hugo release" command.
// Note: This is a command only meant for internal use and must be run
// via "go run -tags release main.go release" on the actual code base that is in the release.
func buildHugoReleaseCmd(hugo *commandeer) *hugoReleaseCmd {
	h := &hugoReleaseCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:    "release",
		Short:  "Release a new version of Hugo.",
		Hidden: true,
		RunE:   h.release,
	}

	h.cmd.PersistentFlags().StringVarP(&h.version, "rel", "r", "", "new release version, i.e. 0.25.1")
	h.cmd.PersistentFlags().BoolVarP(&h.skipPublish, "skip-publish", "", false, "skip all publishing pipes of the release")
	h.cmd.PersistentFlags().BoolVarP(&h.try, "try", "", false, "simulate a release, i.e. no changes")

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoReleaseCmd struct {
	c *commandeer
	cmd *cobra.Command

	version string
	skipPublish bool
	try         bool
}

func (h *hugoReleaseCmd) release(cmd *cobra.Command, args []string) error {
	if h.version == "" {
		return errors.New("must set the --rel flag to the relevant version number")
	}

	return releaser.New(h.version, h.skipPublish, h.try).Run()
}
