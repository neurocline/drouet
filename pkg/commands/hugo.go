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

// Package commands handles Hugo command-line processing.
//
// Hugo commands and flags are implemented using Cobra.
package commands

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

// Build command processor and execute
func Execute() {
    cmd := buildCommand()

    if c, err := cmd.ExecuteC(); err != nil {
        c.Println("")
        c.Println(c.UsageString())
        os.Exit(-1)
    }
}

// Build the Hugo command - root and all its children
// (every other command (verb) is attached as a child command)
func buildCommand() *cobra.Command {
    root := buildHugoCommand()

//    root.AddCommand(buildHugoBenchmarkCmd)
//    root.AddCommand(buildHugoCheckCmd)
//    root.AddCommand(buildHugoConfigCmd)
//    root.AddCommand(buildHugoConvertCmd)
//    root.AddCommand(buildHugoEnvCmd)
//    root.AddCommand(buildHugoGenCmd)
//    root.AddCommand(buildHugoImportCmd)
//    root.AddCommand(buildHugoListCmd)
//    root.AddCommand(buildHugoNewCmd)
//    root.AddCommand(buildHugoServerCmd)
//    root.AddCommand(buildHugoVersionCmd)

    return root
}

// Build Hugo root command.
func buildHugoCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "hugo",
        Short: "hugo builds your site",
        Long: `hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io/.`,
        RunE: hugo,
    }

    return cmd
}

// ----------------------------------------------------------------------------------------------

// "hugo" with no verb is "hugo build", build a site
func hugo(cmd *cobra.Command, args []string) error {
    fmt.Println("hugo")
    return nil
}
