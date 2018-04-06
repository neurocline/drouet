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

package core

import (

	"io/ioutil"
	"os"

	"github.com/neurocline/viper"

	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
)

// NewHugo creates an instance of the top-level state for Hugo operation.
// Many fields of this are refined after initial creation (e.g. logging
// is immediately initialized but only at a critical level and could be
// updated once config is parsed).
func NewHugo() *Hugo {

	// Create an initial logger (this will be replaced once config is parsed)
	startupLogger := jww.NewNotepad(jww.LevelError, jww.LevelError, os.Stdout, ioutil.Discard, "", 0)

	return &Hugo{
		Log: startupLogger,
	}
}

// Hugo is the main object; it contains default behavior
type Hugo struct {
	sites HugoSites

	Log *jww.Notepad
	Config *viper.Viper // this is Config *config.Provider in Hugo
	Fs afero.Fs
}

// HugoSites represents the sites to build. Each site represents a language.
// There is usually just one of these, since it is itself a list of sites.
type HugoSites struct {
	Log *jww.Notepad
}
