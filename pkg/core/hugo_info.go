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

package core

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/spf13/cast"
)

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
// This should be the only one.
var CurrentHugoVersion = HugoVersion{
	Number:     0.38,
	PatchLevel: 0,
	Suffix:     "-DEV",
}

// CurrentHugoGenerator returns the Hugo meta generator string
func CurrentHugoGenerator() template.HTML {
	s := fmt.Sprintf(`<meta name="generator" content="Hugo %s" />`, CurrentHugoVersion.String())
	return template.HTML(s)
}

// ----------------------------------------------------------------------------------------------

// This is in my opinion an overly-complicated system to manipulate version numbers, given
// that it's almost entirely used to generate a few strings, which could just be hand-edited
// or transformed from a base string.

// HugoVersion represents the Hugo build version.
type HugoVersion struct {
	// Major and minor version.
	Number float32

	// Increment this for bug releases
	PatchLevel int

	// HugoVersionSuffix is the suffix used in the Hugo version string.
	// It will be blank for release versions.
	Suffix string
}

type HugoVersionString string

func (v HugoVersion) String() string {
	return hugoVersion(v.Number, v.PatchLevel, v.Suffix)
}

func (v HugoVersion) Version() HugoVersionString {
	return HugoVersionString(v.String())
}

func (h HugoVersionString) String() string {
	return string(h)
}

func hugoVersion(version float32, patchVersion int, suffix string) string {
	if patchVersion > 0 {
		return fmt.Sprintf("%.2f.%d%s", version, patchVersion, suffix)
	}
	return fmt.Sprintf("%.2f%s", version, suffix)
}

// Implements compare.Eqer
func (h HugoVersionString) Eq(other interface{}) bool {
	s, err := cast.ToStringE(other)
	if err != nil {
		return false
	}
	return s == h.String()
}

// --------------------------------

// MustParseHugoVersion parses a version string
// and panics if any error occurs.
func MustParseHugoVersion(s string) HugoVersion {
	vv, err := ParseHugoVersion(s)
	if err != nil {
		panic(err)
	}
	return vv
}

// ParseHugoVersion parses a version string.
func ParseHugoVersion(s string) (HugoVersion, error) {
	var versionSuffixes = []string{"-test", "-DEV"}

	var vv HugoVersion
	for _, suffix := range versionSuffixes {
		if strings.HasSuffix(s, suffix) {
			vv.Suffix = suffix
			s = strings.TrimSuffix(s, suffix)
		}
	}

	v, p := parseVersion(s)

	vv.Number = v
	vv.PatchLevel = p

	return vv, nil
}

func parseVersion(s string) (float32, int) {
	var (
		v float32
		p int
	)

	if strings.Count(s, ".") == 2 {
		li := strings.LastIndex(s, ".")
		p = cast.ToInt(s[li+1:])
		s = s[:li]
	}

	v = float32(cast.ToFloat64(s))

	return v, p
}

// --------------------------------

// ReleaseVersion represents the release version.
func (v HugoVersion) ReleaseVersion() HugoVersion {
	v.Suffix = ""
	return v
}

// Next returns the next Hugo release version.
func (v HugoVersion) Next() HugoVersion {
	return HugoVersion{Number: v.Number + 0.01}
}

// NextPatchLevel returns the next patch/bugfix Hugo version.
// This will be a patch increment on the previous Hugo version.
// TBD only used by test code at the moment.
func (v HugoVersion) NextPatchLevel(level int) HugoVersion {
	return HugoVersion{Number: v.Number - 0.01, PatchLevel: level}
}

// --------------------------------

func compareVersionsWithSuffix(inVersion float32, inPatchVersion int, suffix string, in interface{}) int {
	var c int
	switch d := in.(type) {
	case float64:
		c = compareFloatVersions(inVersion, float32(d))
	case float32:
		c = compareFloatVersions(inVersion, d)
	case int:
		c = compareFloatVersions(inVersion, float32(d))
	case int32:
		c = compareFloatVersions(inVersion, float32(d))
	case int64:
		c = compareFloatVersions(inVersion, float32(d))
	default:
		s, err := cast.ToStringE(in)
		if err != nil {
			return -1
		}

		v, err := ParseHugoVersion(s)
		if err != nil {
			return -1
		}

		if v.Number == inVersion && v.PatchLevel == inPatchVersion {
			return strings.Compare(suffix, v.Suffix)
		}

		if v.Number < inVersion || (v.Number == inVersion && v.PatchLevel < inPatchVersion) {
			return -1
		}

		return 1
	}

	if c == 0 && suffix != "" {
		return 1
	}

	return c
}

func compareFloatVersions(version float32, v float32) int {
	if v == version {
		return 0
	}
	if v < version {
		return -1
	}
	return 1
}

// --------------------------------

// Only used by test code (move to test code?)
func compareVersions(inVersion float32, inPatchVersion int, in interface{}) int {
	return compareVersionsWithSuffix(inVersion, inPatchVersion, "", in)
}
