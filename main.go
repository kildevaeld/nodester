package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

type Config struct {
	Root    string
	Default string
}

func main() {

	/*fi, _ := os.Stdout.Stat() // get the FileInfo struct describing the standard input.

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is from pipe")
		// do things for data from pipe

		bytes, _ := ioutil.ReadAll(os.Stdin)
		str := string(bytes)
		fmt.Println(str)

	} else {
		fmt.Println("data is from terminal")
		// do things from data from terminal

		ConsoleReader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter your name : ")

		input, err := ConsoleReader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Your name is : ", input)

	}*/

	app := cli.NewApp()
	app.Name = "Nodester"
	app.Version = "0.0.6"
	app.Author = "Rasmus Kildevæld"
	app.Email = "rasmuskildevaeld@gmail.com"

	app.Commands = Commands()

	err := app.Run(os.Args)

	if err != nil {
		fmt.Printf(err.Error())
	}

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
			Usage:     "Use nodejs (0.12.0) or IO.js (io@2.0.0) version",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "force, f",
				},
				cli.BoolFlag{
					Name: "migrate, m",
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
		cli.Command{
			Name:      "migrate",
			ShortName: "m",
			Before:    n.init,
			Action:    n.Migrate,
		},
	}

}
