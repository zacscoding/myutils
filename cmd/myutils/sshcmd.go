package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zacscoding/myutils/host"
	"github.com/zacscoding/myutils/remote"
	"github.com/zacscoding/myutils/types"
	"log"
	"strings"
	"sync"
)

var (
	sshCommand = cli.Command{
		Action:   ShowSubCommand,
		Name:     "ssh",
		Usage:    "command for ssh [shell]",
		Category: "SSH COMMANDS",
		Subcommands: []cli.Command{
			{
				Name:      "shell",
				Usage:     "open remote shell",
				Action:    openRemoteShell,
				ArgsUsage: "[host name]",
			},
			{
				Name:      "command",
				Usage:     "execute given comment to a host",
				Action:    executeCommands,
				ArgsUsage: "[comma separated list of host name]",
			},
		},
	}
)

// openRemoteShell start to open remote shell given cli context
func openRemoteShell(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("invalid arguments")
	}

	h, err := host.GetHost(app.db, ctx.Args()[0])
	if err != nil {
		return err
	}

	conn, err := remote.CreateSSHClient(h)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer conn.Close()

	return remote.OpenRemoteShell(conn)
}

// executeCommands execute given command to hosts
func executeCommands(ctx *cli.Context) error {
	if ctx.NArg() != 2 {
		return errors.New("invalid arguments")
	}

	hostNames := strings.Split(ctx.Args()[0], ",")
	command := ctx.Args()[1]

	var hosts []*types.Host
	for _, hostName := range hostNames {
		h, err := host.GetHost(app.db, hostName)
		if err != nil {
			fmt.Println("failed to find a host. name :", hostName)
			continue
		}
		hosts = append(hosts, h)
	}

	commandGen := func(h *types.Host) string {
		return command
	}

	mux := &sync.Mutex{}
	var successes, failures []string

	resultHandler := func(result remote.HostCmdResult) {
		var res string
		if result.Err != nil {
			mux.Lock()
			failures = append(failures, result.Host.Name)
			mux.Unlock()
			res = "fail"
		} else {
			mux.Lock()
			successes = append(successes, result.Host.Name)
			mux.Unlock()
			res = "success"
		}
		var out bytes.Buffer
		out.WriteString("// ------------------------------------------------\n")
		out.WriteString(fmt.Sprintf("host : %s, result : %s, command : %s\n", result.Host.Name, res, result.Command))
		if result.Err != nil {
			out.WriteString(fmt.Sprintf("> error :%v\n", result.Err))
		} else {
			out.WriteString(fmt.Sprintf("> standard output:\n%s\n", result.Result.StdOut))
			out.WriteString(fmt.Sprintf("> standard error:\n%s\n", result.Result.StdErr))
		}
		out.WriteString("--------------------------------------------------- //")
		fmt.Println(out.String())
	}
	remote.ExecutesCommand(hosts, commandGen, resultHandler)
	fmt.Printf(">> Success : %v, Fail : %v\n", successes, failures)
	return nil
}
