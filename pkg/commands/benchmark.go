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

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build "hugo benchmark" command.
func buildHugoBenchmarkCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmark",
		Short: "Benchmark Hugo by building a site a number of times.",
		Long: `Hugo can build a site many times over and analyze the running process
creating a benchmark.`,
		RunE: h.benchmark,
	}

	// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
	initHugoBuilderFlags(cmd)

	// Add flags shared by benchmarking: "hugo", "hugo benchmark"
	initHugoBenchmarkFlags(cmd)

	// Add flags unique to "hugo benchmark"
	cmd.Flags().String("cpuprofile", "", "path/filename for the CPU profile file")
	cmd.Flags().String("memprofile", "", "path/filename for the memory profile file")
	cmd.Flags().IntP("count", "n", 13, "number of times to build the site")

	return cmd
}

// ----------------------------------------------------------------------------------------------

func (h *hugoCmd) benchmark(cmd *cobra.Command, args []string) error {
	fmt.Println("hugo benchmark - benchmark site goes here")
	return nil
}
