package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/gudmundur/ds"
)

func main() {
	app := cli.NewApp()
	app.Name = "ds"
	app.Usage = "syncs dropbox folders"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token",
			EnvVar: "DROPBOX_TOKEN",
			Usage:  "Dropbox OAuth2 token",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "sync",
			Usage: "Syncs from Dropbox",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "directory",
					Value: "tmp",
					Usage: "Destination directory",
				},
			},
			Action: func(c *cli.Context) {
				token := c.GlobalString("token")
				ds := ds.NewDropboxSync(token)
				ds.Directory = c.String("directory")
				ds.Sync()
			},
		},
	}

	app.Run(os.Args)
}
