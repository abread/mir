// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/cobra"
)

var (
	id             string
	membershipFile string
	verbose        bool
	cpuprofile     string
	memprofile     string

	rootCmd = &cobra.Command{
		Use:   "bench",
		Short: "Mir benchmarking tool",
		Long: "Mir benchmarking tool can run a Mir nodes, measuring latency and" +
			"throughput. The tool can also generate and submit requests to the Mir" +
			"cluster at a specified rate.",
	}
)

func Execute() {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create CPU profile: %v", err)
		}

		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "could not start CPU profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not create memory profile: %v", err)
		}

		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "could not write memory profile: %v", err)
		}
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "node/client ID")
	_ = rootCmd.MarkPersistentFlagRequired("id")
	rootCmd.PersistentFlags().StringVarP(&membershipFile, "membership", "m", "", "total number of nodes")
	_ = rootCmd.MarkPersistentFlagRequired("membership")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to file")
	rootCmd.PersistentFlags().StringVar(&memprofile, "memprofile", "", "write memory profile to file")
}
