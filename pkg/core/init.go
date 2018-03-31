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

// Package core contains the core Hugo logic.
package core

import (
	"io/ioutil"
	"os"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/afero"
	"github.com/neurocline/viper" // private fork for now, until merge into mainline
)

// NewHugo creates an instance of the top-level state for Hugo operation.
// Many fields of this are refined after initial creation (e.g. logging
// is immediately initialized but only at a critical level and could be
// updated once config is parsed).
func NewHugo() *Hugo {

	// Create an initial logger (this will be replaced once config is parsed)
	startupLogger := jww.NewNotepad(jww.LevelError, jww.LevelError, os.Stdout, ioutil.Discard, "", 0)

	return &Hugo{
		Logger: startupLogger,
	}
}

type Hugo struct {
	Logger *jww.Notepad
	Config *viper.Viper // TBD do we want an interface here that we can override?
	Fs *afero.Fs
}

// ----------------------------------------------------------------------------------------------

// Init does global-level init like setting up logging.
// If Init has an error, it exits with an error status.
func Init() {

	// We want to have only as many goroutines as machine cores
	// This is no longer necessary as of Go 1.5 (where the default
	// became the number of cores)
	//runtime.GOMAXPROCS(runtime.NumCPU())
}

// Shutdown does clean shutdown, and reports user-level errors that
// were deferred until execution end.
func (h *Hugo) Shutdown() int {

	// If we had any log.ERROR output (even if it was suppressed somehow),
	// we exit with an error code so that systems using hugo as a tool
	// know something went wrong.
	// TBD get rid of all global logging calls.
	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		return -1
	}

	if h.Logger.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		return -1
	}

	return 0
}
