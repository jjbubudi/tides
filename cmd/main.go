package main

import (
	"log"
	"os"

	"github.com/jjbubudi/tides/cmd/fetch"
	"github.com/jjbubudi/tides/cmd/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tides"
	app.HelpName = "tides"
	app.Usage = "Fetches tidal data from the Hong Kong Observatory"
	app.Commands = []cli.Command{
		server.NewServerCommand(),
		fetch.NewFetchCommand(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
