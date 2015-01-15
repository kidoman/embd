package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/kidoman/embd"
)

func detect(c *cli.Context) {
	host, rev, err := embd.DetectHost()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("detected host %v (rev %#x)\n", host, rev)
}

var detectCmd = cli.Command{
	Name:   "detect",
	Usage:  "detect and display information about the host",
	Action: detect,
}

func init() {
	registerCommand(detectCmd)
}
