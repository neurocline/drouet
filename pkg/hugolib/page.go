// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"github.com/neurocline/drouet/pkg/source"
)

// Stubs for now
type Page struct {
	Source
}

func (p *Page) IsDraft() bool {
	return false
}

func (p *Page) IsFuture() bool {
	return false
}

func (p *Page) IsExpired() bool {
	return false
}

type Source struct {
	Frontmatter []byte
	Content     []byte
	source.File
}

type Pages []*Page
