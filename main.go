package main

import (
	"github.com/strategicpause/cgstat/controller"
	"github.com/strategicpause/cgstat/stats"

	"errors"
	"flag"
	"fmt"
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

func ValidateArguments(cgstatArgs *stats.CgstatArgs) error {
	if cgstatArgs.CgroupName == "" && cgstatArgs.CgroupPrefix == "" {
		return errors.New("cgroup name or prefix must be specified")
	}
	if cgstatArgs.VerboseOutput && cgstatArgs.CgroupPrefix != "" {
		return errors.New("you must specify a cgroup name when using verbose output")
	}
	// TODO: Validate path
	// TODO: Validate interval
	return nil
}
