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

	"github.com/spf13/cobra"
)

// Build "hugo release" command.
// Note: This is a command only meant for internal use and must be run
// via "go run -tags release main.go release" on the actual code base that is in the release.
func buildHugoReleaseCmd(h *hugoCmd) *releaseCommandeer {
	r := &releaseCommandeer{}

	r.cmd = &cobra.Command{
		Use:    "release",
		Short:  "Release a new version of Hugo.",
		Hidden: true,
		RunE:   func(cmd *cobra.Command, args []string) error {
			return r.release(h)
		},
	}

	r.cmd.PersistentFlags().StringVarP(&r.version, "rel", "r", "", "new release version, i.e. 0.25.1")
	r.cmd.PersistentFlags().BoolVarP(&r.skipPublish, "skip-publish", "", false, "skip all publishing pipes of the release")
	r.cmd.PersistentFlags().BoolVarP(&r.try, "try", "", false, "simulate a release, i.e. no changes")

	return r
}

// ----------------------------------------------------------------------------------------------

type releaseCommandeer struct {
	cmd *cobra.Command

	version string

	skipPublish bool
	try         bool
}

func (r *releaseCommandeer) release(h *hugoCmd) error {
	if r.version == "" {
		return errors.New("must set the --rel flag to the relevant version number")
	}

	return releaser.New(r.version, r.skipPublish, r.try).Run()
}
