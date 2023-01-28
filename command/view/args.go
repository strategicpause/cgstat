package view

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"time"
)

const (
	ArgName            = "name"
	ArgPrefix          = "prefix"
	ArgVerbose         = "verbose"
	ArgOut             = "out"
	ArgFollow          = "follow"
	ArgRefreshInterval = "refresh-interval"
)

type Args struct {
	CgroupName      string
	CgroupPrefix    string
	VerboseOutput   bool
	OutputFile      string
	FollowMode      bool
	RefreshInterval float64
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Name of cgroup",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "Cgroup prefix",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Prints verbose information about a single cgroup.",
		},
		cli.StringFlag{
			Name:  "out",
			Usage: "Writes to a given file if provided.",
		},
		cli.BoolFlag{
			Name:  "follow",
			Usage: "Refreshes the output every interval.",
		},
		cli.Float64Flag{
			Name:  "refresh-interval",
			Usage: "Refresh interval in seconds",
			Value: 1.0,
		},
	}
}

func parseArgs(cCtx *cli.Context) (*Args, error) {
	viewArgs := &Args{
		CgroupName:      cCtx.String(ArgName),
		CgroupPrefix:    cCtx.String(ArgPrefix),
		VerboseOutput:   cCtx.Bool(ArgVerbose),
		OutputFile:      cCtx.String(ArgOut),
		FollowMode:      cCtx.Bool(ArgFollow),
		RefreshInterval: cCtx.Float64(ArgRefreshInterval),
	}

	if err := validateArguments(viewArgs); err != nil {
		return nil, fmt.Errorf("error parsing list args: %s", err)
	}

	return viewArgs, nil
}

func validateArguments(args *Args) error {
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

func (a *Args) HasPrefix() bool {
	return a.CgroupPrefix != ""
}

func (a *Args) HasOutputFile() bool {
	return a.OutputFile != ""
}

func (a *Args) GetRefreshInterval() time.Duration {
	return time.Duration(a.RefreshInterval * float64(time.Second))
}
