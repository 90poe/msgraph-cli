package main

import (
	"os"

	"github.com/mhoc/msgraph-cli/auth"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "msgraph"
	app.Version = "0.1.0"
	app.Usage = "cli tool for accessing the microsoft graph api"
	app.Commands = []cli.Command{
		{
			Name:        "auth",
			Usage:       "subcommands related to authenticating with the graph api",
			Subcommands: auth.Commands(),
		},
	}
	app.Run(os.Args)
}
