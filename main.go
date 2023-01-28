package main

import (
	"log"
	"os"

	"github.com/strategicpause/cgstat/command/list"
	"github.com/strategicpause/cgstat/command/view"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Commands: RegisterCommands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RegisterCommands() cli.Commands {
	return cli.Commands{
		list.Register(),
		view.Register(),
	}
}
