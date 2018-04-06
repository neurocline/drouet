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

// Caveat with Viper - merging between config files or between
// levels of the Viper sources hierarchy works for individual
// values, but not maps or slices. This is not a concern for the
// built-in Hugo config (which is all individual values), but it
// should be noted somewhere, and perhaps a verbose warning that
// merging threw data away

package core

import (
	"fmt"
	"strings"
	"sync"

	jww "github.com/spf13/jwalterweatherman"
)

// ----------------------------------------------------------------------------------------------

// InitLoggers sets up the global distinct loggers.
func InitDistinctLoggers() {
	DistinctErrorLog = NewDistinctErrorLogger()
	DistinctWarnLog = NewDistinctWarnLogger()
	DistinctFeedbackLog = NewDistinctFeedbackLogger()
}

var (
	// DistinctErrorLog can be used to avoid spamming the logs with errors.
	DistinctErrorLog *DistinctLogger = NewDistinctErrorLogger()

	// DistinctWarnLog can be used to avoid spamming the logs with warnings.
	DistinctWarnLog *DistinctLogger = NewDistinctWarnLogger()

	// DistinctFeedbackLog can be used to avoid spamming the logs with info messages.
	DistinctFeedbackLog *DistinctLogger = NewDistinctFeedbackLogger()
)

// NewDistinctErrorLogger creates a new DistinctLogger that logs ERRORs
func NewDistinctErrorLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.ERROR}
}

// NewDistinctWarnLogger creates a new DistinctLogger that logs WARNs
func NewDistinctWarnLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.WARN}
}

// NewDistinctFeedbackLogger creates a new DistinctLogger that can be used
// to give feedback to the user while not spamming with duplicates.
func NewDistinctFeedbackLogger() *DistinctLogger {
	return &DistinctLogger{m: make(map[string]bool), logger: jww.FEEDBACK}
}

// DistinctLogger ignores duplicate log statements.
type DistinctLogger struct {
	sync.RWMutex
	logger logPrinter
	m      map[string]bool
}

type logPrinter interface {
	// Println is the only common method that works in all of JWWs loggers.
	Println(a ...interface{})
}

// Println will log the string returned from fmt.Sprintln given the arguments,
// but not if it has been logged before.
func (l *DistinctLogger) Println(v ...interface{}) {
	// fmt.Sprint doesn't add space between string arguments
	logStatement := strings.TrimSpace(fmt.Sprintln(v...))
	l.print(logStatement)
}

// Printf will log the string returned from fmt.Sprintf given the arguments,
// but not if it has been logged before.
// Note: A newline is appended.
func (l *DistinctLogger) Printf(format string, v ...interface{}) {
	logStatement := fmt.Sprintf(format, v...)
	l.print(logStatement)
}

func (l *DistinctLogger) print(logStatement string) {
	l.RLock()
	if l.m[logStatement] {
		l.RUnlock()
		return
	}
	l.RUnlock()

	l.Lock()
	if !l.m[logStatement] {
		l.logger.Println(logStatement)
		l.m[logStatement] = true
	}
	l.Unlock()
}
