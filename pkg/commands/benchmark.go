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
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/neurocline/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// Build "hugo benchmark" command.
func buildHugoBenchmarkCmd(hugo *commandeer) *hugoBenchmarkCmd {
	h := &hugoBenchmarkCmd{c: hugo}

	h.cmd = &cobra.Command{
		Use:   "benchmark",
		Args:  cobra.NoArgs,
		Short: "Benchmark Hugo by building a site a number of times.",
		Long: `Hugo can build a site many times over and analyze the running process
creating a benchmark.`,
		RunE: h.benchmark,
	}

	// Add flags for the "hugo benchmark" command
	h.cmd.Flags().IntVarP(&h.benchmarkTimes, "count", "n", 13, "number of times to build the site")
	h.cmd.Flags().StringVar(&h.cpuProfileFile, "cpuprofile", "", "path/filename for the CPU profile file")
	h.cmd.Flags().StringVar(&h.memProfileFile, "memprofile", "", "path/filename for the memory profile file")
	h.cmd.Flags().BoolVar(&h.renderToMemory, "renderToMemory", false, "render to memory (useful for benchmark testing)")

	// Add flags shared by builders: "hugo", "hugo server", "hugo benchmark"
	addHugoBuilderFlags(h.cmd)

	return h
}

// ----------------------------------------------------------------------------------------------

type hugoBenchmarkCmd struct {
	c *commandeer
	cmd *cobra.Command

	benchmarkTimes int
	cpuProfileFile string
	memProfileFile string
	renderToMemory bool
}

func (h *hugoBenchmarkCmd) benchmark(cmd *cobra.Command, args []string) error {
	var err error

	// Load config
	cfgInit := func(c *commandeer) error {
		c.Set("renderToMemory", h.renderToMemory)
		return nil
	}

	if err = h.c.InitializeConfig(cfgInit, h.cmd); err != nil {
		return err
	}

	// If --memprofile=<FILE>, then create log file for memory profiling
	var memProf *os.File
	if h.memProfileFile != "" {
		memProf, err = os.Create(h.memProfileFile)
		if err != nil {
			return err
		}
		defer memProf.Close()
	}

	// If --cpuprofile=<FILE>, then create log file for cpu profiling
	var cpuProf *os.File
	if h.cpuProfileFile != "" {
		cpuProf, err = os.Create(h.cpuProfileFile)
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
	t := time.Now()
	//builder := &hugoCmd{c: h.c, cmd: h.cmd, renderToMemory: h.renderToMemory}
	for i := 0; i < h.benchmarkTimes; i++ {
		//if err = builder.resetAndBuildSites(); err != nil {
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
	jww.FEEDBACK.Printf("Average time per operation: %vms\n", int(1000*totalTime.Seconds()/float64(h.benchmarkTimes)))
	jww.FEEDBACK.Printf("Average memory allocated per operation: %vkB\n", totalMemAllocated/uint64(h.benchmarkTimes)/1024)
	jww.FEEDBACK.Printf("Average allocations per operation: %v\n", totalMallocs/uint64(h.benchmarkTimes))

	return nil
}
