// Copyright 2017-present The Hugo Authors. All rights reserved.
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

// Package releaser implements a set of utilities and a wrapper around Goreleaser
// to help automate the Hugo release process.
package releaser

import (
	"fmt"
	"strings"
)

// ReleaseHandler provides functionality to release a new version of Hugo.
type ReleaseHandler struct {
	cliVersion string

	skipPublish bool

	// Just simulate, no actual changes.
	try bool

	git func(args ...string) (string, error)
}

// New initialises a ReleaseHandler.
func New(version string, skipPublish, try bool) *ReleaseHandler {
	// When triggered from CI release branch
	version = strings.TrimPrefix(version, "release-")
	version = strings.TrimPrefix(version, "v")
	rh := &ReleaseHandler{cliVersion: version, skipPublish: skipPublish, try: try}

	if try {
		rh.git = func(args ...string) (string, error) {
			fmt.Println("git", strings.Join(args, " "))
			return "", nil
		}
	} else {
		rh.git = git
	}

	return rh
}

// Run creates a new release.
func (r *ReleaseHandler) Run() error {
	// TBD fill in
	return nil
}
