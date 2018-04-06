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

// core/globals.go holds all of the globals declared by the core package.
// This is because we are deprecating globals, so the Hugo code itself
// won't refer to these globals, but they exist until third-party
// packages migrate away from accessing Hugo globals.

package core

// This is a rare legitimate use of program globals - we use -ldflags
// to set the values of these at build time. So consider these constants
// instead of variables. These are set by "mage install", the makefile
// used to build Hugo.

var (
	// CommitHash contains the current Git revision.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string
)

// CurrentHugoVersion represents the current build version.
// This is a global for convenience.
var CurrentHugoVersion = HugoVersion{
	Number:     0.39,
	PatchLevel: 0,
	Suffix:     "-DEV",
}
