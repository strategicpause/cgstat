package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/strategicpause/cgstat/controller"
	"github.com/strategicpause/cgstat/stats"
	"os"
	"path/filepath"
)

func main() {
	cgstatArgs := ParseArguments()
	err := ValidateArguments(cgstatArgs)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = PrintStats(cgstatArgs)
	if err != nil {
		fmt.Println(err)
	}
}

func PrintStats(args *stats.CgstatArgs) error {
	controller, err := controller.NewCgroupStatsController(args)
	if err != nil {
		return err
	}
	return controller.Start()
}

func ParseArguments() *stats.CgstatArgs {
	cgstatArgs := stats.CgstatArgs{}

	flag.StringVar(&cgstatArgs.CgroupName, "name", "", "Name of cgroup")
	flag.StringVar(&cgstatArgs.CgroupPrefix, "prefix", "", "Cgroup prefix")
	flag.BoolVar(&cgstatArgs.VerboseOutput, "verbose", false, "Prints verbose information about a single cgroup")
	flag.StringVar(&cgstatArgs.OutputFile, "out", "", "Writes to a given file if provided.")
	flag.BoolVar(&cgstatArgs.FollowMode, "follow", false, "Refreshes the output every interval")
	flag.Float64Var(&cgstatArgs.RefreshInterval, "refresh", 1.0, "Refresh interval in seconds")

	flag.Parse()

	return &cgstatArgs
}

func ValidateArguments(args *stats.CgstatArgs) error {
	if args.CgroupName == "" && args.CgroupPrefix == "" {
		return errors.New("cgroup name or prefix must be specified")
	}
	if args.VerboseOutput && args.CgroupPrefix != "" {
		return errors.New("you must specify a cgroup name when using verbose output")
	}
	if args.RefreshInterval < 0.0 {
		return errors.New("you must specify a non-negative refresh interval")
	}
	if args.HasOutputFile() {
		base, err := filepath.Abs(args.OutputFile)
		if err != nil {
			return err
		}
		_, err = os.Stat(filepath.Dir(base))
		if err != nil {
			return err
		}
	}
	return nil
}
