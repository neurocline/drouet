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

package source

import (
	//"os"
	//"path/filepath"
	"regexp"

	//"github.com/neurocline/drouet/core"

	"github.com/spf13/afero"
	//"github.com/spf13/cast"
)

// SourceSpec abstracts language-specific file creation.
// TODO(bep) rename to Spec
type SourceSpec struct {
	//*helpers.PathSpec

	Fs afero.Fs

	// This is set if the ignoreFiles config is set.
	ignoreFilesRe []*regexp.Regexp

	Languages              map[string]interface{}
	DefaultContentLanguage string
	DisabledLanguages      map[string]bool
}
