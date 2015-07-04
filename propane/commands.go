package propane

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandServer,
}

var commandServer = cli.Command{
	Name:  "server",
	Usage: "",
	Description: `
`,
	Action: doServer,
	Flags: []cli.Flag {
		cli.StringFlag{
			Name: "bind-addr, a",
			Value: "0.0.0.0",
			Usage: "--bind-addr=0.0.0.0",
		},
		cli.IntFlag{
			Name: "port, p",
			Value: 8088,
			Usage: "--port=8088",
		},
	},
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doServer(c *cli.Context) {
	RunServer(c)
}
