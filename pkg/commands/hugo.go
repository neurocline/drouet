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

// Package commands defines and implements command-line commands and flags
// used by Hugo. Commands and flags are implemented using Cobra.
package commands

import (
	"github.com/neurocline/drouet/pkg/hugolib"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// The Response value from Execute.
type Response struct {
	// The build Result will only be set in the hugo build command.
	Result *hugolib.HugoSites

	// Err is set when the command failed to execute.
	Err error

	// The command that was executed.
	Cmd *cobra.Command
}

func (r Response) IsUserError() bool {
	return r.Err != nil && isUserError(r.Err)
}

// Execute adds all child commands to the root command HugoCmd and sets flags appropriately.
// The args are usually filled with os.Args[1:].
func Execute(args []string) Response {
	hugoCmd := newCommandsBuilder().addAll().build()
	cmd := hugoCmd.getCommand()
	cmd.SetArgs(args)

	c, err := cmd.ExecuteC()

	var resp Response

	if c == cmd && hugoCmd.c != nil {
		// Root command executed
		resp.Result = hugoCmd.c.hugo
	}

	if err == nil {
		errCount := jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError)
		if errCount > 0 {
			err = fmt.Errorf("logged %d errors", errCount)
		} else if resp.Result != nil {
			errCount = resp.Result.Log.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError)
			if errCount > 0 {
				err = fmt.Errorf("logged %d errors", errCount)
			}
		}

	}

	resp.Err = err
	resp.Cmd = c

	return resp
}
