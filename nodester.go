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
		fmt.Println("Wrong usage! You must specify a version")
		os.Exit(1)
	}
	version := args.First()
	force := c.Bool("force")

	if !n.Node.Has(version) {
		if !force {
			fmt.Printf("node: %s not installed", version)
			os.Exit(1)
		} else {
			n.Install(c)
		}
	}
	fmt.Printf("Use version: %s\n", version)
	n.Node.Use(version)

}

func (n *NodeCli) Install(c *cli.Context) {

	args := c.Args()
	if len(args) == 0 {
		fmt.Println("Wrong usage! You must specify a version")
		os.Exit(1)
	}
	for _, version := range args {
		_, err := n.Node.Download(version, func(p DownloadProgress) {
			str := fmt.Sprintf("Downloading... %d/%d kb\r", p.Progress/1024, p.Total/1014)
			os.Stdout.Write([]byte(str))
			if p.Progress == p.Total {
				os.Stdout.WriteString("\033[2K\rDownloading... Done\n")
			}
		})

		if err != nil {
			os.Stdout.WriteString("\033[2K\rDownloading... Error:" + err.Error() + "\n")
			os.Exit(1)
		}
		os.Stdout.WriteString("Installing...")
		err = n.Node.Install(version)
		if err != nil {
			os.Stdout.WriteString(" Error!\n")
		} else {
			os.Stdout.WriteString(" Done!\n")
		}
	}
	//version := args.First()

}

func (n *NodeCli) Remove(c *cli.Context) {

	args := c.Args()
	if len(args) == 0 {
		fmt.Println("Wrong usage! You must specify a version")
		os.Exit(1)
	}
	version := args.First()

	os.Stdout.WriteString(fmt.Sprintf("Removing %s...", version))
	err := n.Node.Remove(version)
	if err != nil {
		os.Stdout.WriteString(" Error!\n")
	} else {
		os.Stdout.WriteString(" Done!\n")
	}

}

func (n *NodeCli) Clear(c *cli.Context) {

	os.Stdout.WriteString("Clearing cache...")
	err := n.Node.CleanCache()
	if err != nil {
		os.Stdout.WriteString(" Error!")
	} else {
		os.Stdout.WriteString(" Done!\n")
	}
}

func (n *NodeCli) List(c *cli.Context) {

	versions := n.Node.List()

	fmt.Printf("Versions: %s\n", versions)
}

func (n *NodeCli) ListRemote(c *cli.Context) {

	os.Stdout.WriteString("Fetching remote list...")
	remote, err := n.Node.ListRemote()
	if err != nil {
		os.Stdout.WriteString(" Error!\n")
		os.Exit(1)
	} else {
		os.Stdout.WriteString(" Done!\n")
	}
	fmt.Printf("Remote Versions: %s\n", remote)
}

func (n *NodeCli) Current(c *cli.Context) {
	fmt.Printf("Current: %s\n", n.Node.Current())
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
