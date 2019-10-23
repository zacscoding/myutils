package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os/exec"
)

var (
	scpCommand = cli.Command{
		Action:   executeScpCommand,
		Name:     "scp",
		Usage:    "command for scp",
		Category: "SSH COMMANDS",
	}
)

func executeScpCommand(ctx *cli.Context) error {
	var args = make([]string, ctx.NArg())
	for _, arg := range ctx.Args() {
		args = append(args, arg)
	}
	out, err := exec.Command("scp", args...).Output()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}
