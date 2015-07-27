package propane

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandServer,
	commandNew,
}

var commandServer = cli.Command{
	Name:  "server",
	Usage: "c3h8 server [options]",
	Description: `
`,
	Action: doServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bind-addr, a",
			Value: "0.0.0.0",
			Usage: "--bind-addr=0.0.0.0",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 8088,
			Usage: "--port=8088",
		},
	},
}

var commandNew = cli.Command{
	Name:  "new",
	Usage: "c3h8 new FILENAME",
	Description: `
`,
	Action: doNewFile,
	Flags:  []cli.Flag{},
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

func doNewFile(c *cli.Context) {
	args := c.Args()
	if len(args) > 0 {
		filename := args[0]
		now := time.Now()
		newfilename := now.Format("20060102150405") + "-" + filename + ".md"
		newfilepath := filepath.Join(CurDir(), newfilename)
		src := []byte("#\n")
		err := ioutil.WriteFile(newfilepath, src, 0644)
		if err != nil {
			panic(err)
		}
		fmt.Println("Created " + newfilename)
	} else {
		fmt.Println(`Please input FILENAME

USAGE:

    $ c3h8 new new-file

`)
	}
}
