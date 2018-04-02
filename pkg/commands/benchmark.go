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
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/neurocline/drouet/pkg/core"
	"github.com/neurocline/drouet/pkg/z"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
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

	// Add flags for the "hugo benchmark" command
	cmd.Flags().IntP("count", "n", 13, "number of times to build the site")
	cmd.Flags().String("cpuprofile", "", "path/filename for the CPU profile file")
	cmd.Flags().String("memprofile", "", "path/filename for the memory profile file")
	cmd.Flags().Bool("renderToMemory", false, "render to memory (useful for benchmark testing)")

	// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
	initHugoBuilderFlags(cmd)

	return cmd
}

// ----------------------------------------------------------------------------------------------

func (h *hugoCmd) benchmark(cmd *cobra.Command, args []string) error {

	// Load config
	var err error
	h.Config, err = core.InitializeConfig(h.Hugo, cmd)
	if err != nil {
		return err
	}
	fmt.Fprintf(z.Log, "commands.benchmark()%s\n%s\n", z.Stack(), h.Config.Spew())

	// If --memprofile=<FILE>, then create log file for memory profiling
	var memProf *os.File
	memProfileFile := h.Config.GetString("memprofile")
	if memProfileFile != "" {
		memProf, err = os.Create(memProfileFile)
		if err != nil {
			return err
		}
		defer memProf.Close()
	}

	// If --cpuprofile=<FILE>, then create log file for cpu profiling
	var cpuProf *os.File
	cpuProfileFile := h.Config.GetString("cpuprofile")
	if cpuProfileFile != "" {
		cpuProf, err = os.Create(cpuProfileFile)
		if err != nil {
			return err
		}
		defer cpuProf.Close()
	}

	// Get the baseline for memory allocations
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memAllocated := memStats.TotalAlloc
	mallocs := memStats.Mallocs

	// Start CPU profiling if requested
	if cpuProf != nil {
		pprof.StartCPUProfile(cpuProf)
	}

	// Build site N times, measuring total elapsed time
	benchmarkTimes := h.Config.GetInt("count")
	t := time.Now()
	for i := 0; i < benchmarkTimes; i++ {
		//if err = c.resetAndBuildSites(); err != nil {
		//	return err
		//}
	}
	totalTime := time.Since(t)

	// Stop profiling and write results
	if memProf != nil {
		pprof.WriteHeapProfile(memProf)
	}
	if cpuProf != nil {
		pprof.StopCPUProfile()
	}

	// Show runtime and allocations, averaged over N iterations
	runtime.ReadMemStats(&memStats)
	totalMemAllocated := memStats.TotalAlloc - memAllocated
	totalMallocs := memStats.Mallocs - mallocs

	jww.FEEDBACK.Println()
	jww.FEEDBACK.Printf("Average time per operation: %vms\n", int(1000*totalTime.Seconds()/float64(benchmarkTimes)))
	jww.FEEDBACK.Printf("Average memory allocated per operation: %vkB\n", totalMemAllocated/uint64(benchmarkTimes)/1024)
	jww.FEEDBACK.Printf("Average allocations per operation: %v\n", totalMallocs/uint64(benchmarkTimes))

	return nil
}
