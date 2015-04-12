package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Nodester"
	app.Version = "0.0.5"
	app.Author = "Rasmus Kildevæld"
	app.Email = "rasmuskildevaeld@gmail.com"

	app.Commands = Commands()

	app.Run(os.Args)

}

// Commands is
func Commands() []cli.Command {
	n := &NodeCli{}
	return []cli.Command{
		cli.Command{
			Name:      "use",
			Before:    n.init,
			Action:    n.Use,
			ShortName: "u",
			Usage:     "Use nodejs version",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "force, f",
				},
			},
		},
		cli.Command{
			Name:      "install",
			ShortName: "i",
			Before:    n.init,
			Action:    n.Install,
		},
		cli.Command{
			Name:      "remove",
			ShortName: "rm",
			Before:    n.init,
			Action:    n.Remove,
		},
		cli.Command{
			Name:   "clean",
			Before: n.init,
			Action: n.Clear,
		},
		cli.Command{
			Name:      "list",
			ShortName: "ls",
			Before:    n.init,
			Action:    n.List,
		},
		cli.Command{
			Name:      "list-remote",
			ShortName: "lsr",
			Before:    n.init,
			Action:    n.ListRemote,
		},
		cli.Command{
			Name:   "current",
			Before: n.init,
			Action: n.Current,
		},
	}
}
