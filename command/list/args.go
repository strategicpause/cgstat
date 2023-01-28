package list

import "github.com/urfave/cli"

const (
	ArgsPrefix = "prefix"
)

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  ArgsPrefix,
			Usage: "Cgroup prefix",
			Value: "/",
		},
	}
}
