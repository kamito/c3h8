package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/kamito/c3h8/propane"
)

func main() {
	app := cli.NewApp()
	app.Name = "c3h8"
	app.Version = propane.Version
	app.Usage = ""
	app.Author = "Shinichirow KAMITO"
	app.Email = "updoor@gmail.com"
	app.Commands = propane.Commands

	app.Run(os.Args)
}
