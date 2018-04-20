// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"github.com/neurocline/drouet/pkg/core"
)

// HugoSites represents the sites to build. Each site represents a language.
type HugoSites struct {
	Sites []*Site
}

// NewHugoSites creates HugoSites from the given config.
func NewHugoSites(cfg core.HugoBaseConfig) (*HugoSites, error) {
	return nil, nil
}

// Pages returns all pages for all sites.
func (h *HugoSites) Pages() Pages {
	return h.Sites[0].AllPages
}

// BuildCfg holds build options used to, as an example, skip the render step.
type BuildCfg struct {
	// Reset site state before build. Use to force full rebuilds.
	ResetState bool
	// Re-creates the sites from configuration before a build.
	// This is needed if new languages are added.
	CreateSitesFromConfig bool
	// Skip rendering. Useful for testing.
	SkipRender bool
	// Use this to indicate what changed (for rebuilds).
//	whatChanged *whatChanged
	// Recently visited URLs. This is used for partial re-rendering.
	RecentlyVisited map[string]bool
}
