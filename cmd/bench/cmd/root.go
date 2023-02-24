// Copyright Contributors to the Mir project
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
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

	cpuprofileFile *os.File
	memprofileFile *os.File

	rootCmd = &cobra.Command{
		Use:   "bench",
		Short: "Mir benchmarking tool",
		Long: "Mir benchmarking tool can run a Mir nodes, measuring latency and" +
			"throughput. The tool can also generate and submit requests to the Mir" +
			"cluster at a specified rate.",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if cpuprofile != "" {
				fmt.Printf("saving CPU profile to %s\n", cpuprofile)

				cpuprofileFile, err = os.Create(cpuprofile)
				if err != nil {
					return fmt.Errorf("could not create CPU profile: %w", err)
				}

				err = pprof.StartCPUProfile(cpuprofileFile)
				if err != nil {
					return fmt.Errorf("could not start CPU profile: %w", err)
				}
			}
			if memprofile != "" {
				fmt.Printf("saving memory profile to %s\n", memprofile)

				memprofileFile, err = os.Create(memprofile)
				if err != nil {
					return fmt.Errorf("could not create memory profile: %w", err)
				}
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if cpuprofileFile != nil {
				pprof.StopCPUProfile()
				if err := cpuprofileFile.Close(); err != nil {
					return fmt.Errorf("could not close CPU profile file: %w", err)
				}
			}
			if memprofileFile != nil {
				runtime.GC() // get up-to-date statistics
				if err := pprof.WriteHeapProfile(memprofileFile); err != nil {
					return fmt.Errorf("could not write memory profile: %w", err)
				}
				if err := memprofileFile.Close(); err != nil {
					return fmt.Errorf("could not close memory profile file: %w", err)
				}
			}

			return nil
		},
	}
)

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
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
