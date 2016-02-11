package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

type NodeCli struct {
	Node *NodeManager
}

func (n *NodeCli) Run(c *cli.Context) {
	n.init(c)

}

func (n *NodeCli) Use(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		fmt.Println("  Please specify a version\n")
		os.Exit(1)
	}
	version := args.First()
	force := c.Bool("force")
	migrate := c.Bool("migrate")
	current := n.Node.Current()

	if !n.Node.Has(version) {
		if !force {
			fmt.Printf("  %s not installed\n", version)
			os.Exit(1)
		} else {
			n.Install(c)
		}
	}

	if migrate && current != "" {
		err := NewProcess(fmt.Sprintf("  Migrating modules from %s to %s...", current, version), func() error {
			return n.Node.Migrate(current, version)
		})
		if err != nil {
			fmt.Printf("  Migration failed: %s\n", err.Error())
			os.Exit(1)
		}

	}

	n.Node.Use(version)
	fmt.Printf("  Current version: %s\n", version)

}

func (n *NodeCli) Install(c *cli.Context) {

	args := c.Args()
	if len(args) == 0 {
		fmt.Println("  Please specify a version\n")
		os.Exit(1)
	}
	for _, version := range args {
		err := NewProgress(fmt.Sprintf("  Downloading %s ...", version), func(fn func(str string)) error {
			_, err := n.Node.Download(version, func(p DownloadProgress) {
				str := fmt.Sprintf("%d/%d kb", p.Progress/1024, p.Total/1014)
				fn(str)
			})

			return err
		})

		if err != nil {
			fmt.Printf("  Could not download '%s' due to: %s\n", version, err.Error())
			os.Exit(1)
		}

		NewProcess(fmt.Sprintf("  Installing %s ...", version), func() error {
			return n.Node.Install(version)
		})

	}
}

func (n *NodeCli) Remove(c *cli.Context) {

	args := c.Args()
	if len(args) == 0 {
		fmt.Println("  Plesase sepcify a version\n")
		os.Exit(1)
	}
	version := args.First()

	NewProcess("  Removing "+version+" ...", func() error {
		return n.Node.Remove(version)
	})

}

func (n *NodeCli) Migrate(c *cli.Context) {

	if len(c.Args()) != 2 {
		fmt.Println("  Please spcify a from and to version\n")
		os.Exit(1)
	}

	from := c.Args().First()
	to := c.Args()[1]

	err := NewProcess(fmt.Sprintf("  Migrating modules from %s to %s...", from, to), func() error {
		return n.Node.Migrate(from, to)
	})
	if err != nil {
		fmt.Printf("  Migration failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func (n *NodeCli) Clear(c *cli.Context) {

	NewProcess("  Clearing cache...", func() error {
		return n.Node.CleanCache()
	})

}

func (n *NodeCli) List(c *cli.Context) {

	versions := n.Node.List()

	fmt.Printf("Versions: %s\n", versions)
}

func (n *NodeCli) ListRemote(c *cli.Context) {
	p := &Process{
		Msg: "  Fetching remote list ...",
	}
	p.Start()

	remote, err := n.Node.ListRemote()
	if err != nil {
		p.Done("error")
		os.Exit(1)
	} else {
		p.Done("ok")
	}
	fmt.Printf("  Remote Versions: %s\n\n", remote)
}

func (n *NodeCli) Current(c *cli.Context) {
	fmt.Printf("  Current: %s\n", n.Node.Current())
}

func (n *NodeCli) init(c *cli.Context) (err error) {

	if n.Node == nil {

		path := os.Getenv("NODESTER_ROOT")

		if path == "" {
			home, e := homedir.Dir()
			if e != nil {
				err = e
				return
			}

			defaultPath := filepath.Join(home, ".nodester")

			if exists(defaultPath) {
				n.Node = NewNodeManager(defaultPath)
				return nil
			}

			var tmp string
			var set bool
			for {
				os.Stdout.WriteString("NODSTER_ROOT not set. Should I use default directory: ~/.nodester? [Y/n]")
				fmt.Scanf("%s", &tmp)
				tmp = strings.ToLower(tmp)
				if tmp == "" {
					set = true
					break
				}

				if !strings.Contains("yn", tmp) {
					os.Stdout.WriteString("\033[2K\r")
					continue
				}
				if tmp == "y" {
					set = true
				} else {
					set = false
				}
				break
			}

			if !set {
				err = errors.New("Node root path not defined")
			} else {
				p := filepath.Join(home, ".nodester")
				n.Node = NewNodeManager(p)
			}

		} else {
			n.Node = NewNodeManager(path)
		}
	}

	return err

}
