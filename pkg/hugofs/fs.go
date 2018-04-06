// Copyright 2016 The Hugo Authors. All rights reserved.
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

// Package hugofs provides the file systems used by Hugo.
package hugofs

import (
	"github.com/spf13/afero"
)

// Os points to an Os Afero file system.
var Os = &afero.OsFs{}

// Fs abstracts the file system to separate source and destination file systems
// and allows both to be mocked for testing.
type Fs struct {
	// Source is Hugo's source file system.
	Source afero.Fs

	// Destination is Hugo's destination file system.
	Destination afero.Fs

	// Os is an OS file system.
	// NOTE: Field is currently unused.
	Os afero.Fs

	// WorkingDir is a read-only file system
	// restricted to the project working dir.
	WorkingDir *afero.BasePathFs
}
