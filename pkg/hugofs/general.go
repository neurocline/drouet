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

package hugofs

import (
	"os"
	"path/filepath"

	jww "github.com/spf13/jwalterweatherman"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// This really doesn't belong here. Note that this doesn't work through the afero.Fs that
// Hugo uses, it really is intended for local disk operation (e.g. creating a temp dir)
func Mkdir(x ...string) {
	p := filepath.Join(x...)

	err := os.MkdirAll(p, 0777) // before umask
	if err != nil {
		jww.FATAL.Fatalln(err)
	}
}
