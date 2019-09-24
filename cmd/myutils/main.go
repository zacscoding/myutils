// myutils main package
package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/zacscoding/myutils/utils"
	"os"
)

var (
	app = utils.NewApp()
)

func init() {
	app.Action = func(ctx *cli.Context) error {
		return cli.ShowAppHelp(ctx)
	}

	app.Commands = []cli.Command{
		hostCommand,
	}
}

func main() {
	// os.Args = append(os.Args, "host")
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
