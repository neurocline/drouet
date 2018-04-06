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

package core

import (
	"fmt"
	"regexp"

	jww "github.com/spf13/jwalterweatherman"
)

// ----------------------------------------------------------------------------------------------

// commandError is an error used to signal different error situations in command handling.
type commandError struct {
	s         string
	userError bool
}

func (c commandError) Error() string {
	return c.s
}

func (c commandError) isUserError() bool {
	return c.userError
}

func NewUserError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: true}
}

func NewSystemError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: false}
}

func NewSystemErrorF(format string, a ...interface{}) commandError {
	return commandError{s: fmt.Sprintf(format, a...), userError: false}
}

// Catch some of the obvious user errors from Cobra.
// We don't want to show the usage message for every error.
// The below may be to generic. Time will show.
var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand|unknown command")

func IsUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}

// Deprecated informs about a deprecation, but only once for a given set of arguments' values.
// If the err flag is enabled, it logs as an ERROR (will exit with -1) and the text will
// point at the next Hugo release.
// The idea is two remove an item in two Hugo releases to give users and theme authors
// plenty of time to fix their templates.
func Deprecated(object, item, alternative string, err bool) {
	if err {
		DistinctErrorLog.Printf("%s's %s is deprecated and will be removed in Hugo %s. %s", object, item, CurrentHugoVersion.Next().ReleaseVersion(), alternative)

	} else {
		// Make sure the users see this while avoiding build breakage. This will not lead to an os.Exit(-1)
		DistinctFeedbackLog.Printf("WARNING: %s's %s is deprecated and will be removed in a future release. %s", object, item, alternative)
	}
}

func CheckErr(logger *jww.Notepad, err error, s ...string) {
	if err == nil {
		return
	}
	if len(s) == 0 {
		logger.CRITICAL.Println(err)
		return
	}
	for _, message := range s {
		logger.ERROR.Println(message)
	}
	logger.ERROR.Println(err)
}
