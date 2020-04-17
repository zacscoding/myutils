package main

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/urfave/cli"
	"github.com/zacscoding/myutils/host"
	"github.com/zacscoding/myutils/remote"
	"os"
	"path/filepath"
	"strings"
)

var (
	scpCommand = cli.Command{
		Action:    executeScpCommand,
		Name:      "scp",
		Usage:     "command for scp",
		Category:  "SCP COMMANDS",
		ArgsUsage: "[[hostname]:source] [[hostname]:destination]",
	}
)

func executeScpCommand(ctx *cli.Context) error {
	if ctx.NArg() != 2 {
		return errors.New("required args [[hostname]:source] [[hostname]:destination]")
	}

	// parse source | destination
	srcHost, srcPath := splitPath(ctx.Args()[0])
	destHost, destPath := splitPath(ctx.Args()[1])
	if srcHost == "" && destHost == "" {
		return errors.New("cannot find a host in paths")
	}

	upload := true
	hostName := srcHost

	if srcHost != "" {
		upload = true
		hostName = srcHost
	}
	// getting host
	h, err := host.GetHost(app.db, hostName)
	if err != nil {
		return err
	}
	// create sftp client
	sc, err := remote.CreateSSHClient(h)
	if err != nil {
		return err
	}
	defer sc.Close()

	client, err := sftp.NewClient(sc)
	if err != nil {
		return err
	}
	defer client.Close()

	if upload {
		return uploadFiles(srcPath, destPath, client)
	}
	return downloadFiles(srcPath, destPath, client)
}

// splitPath returns a pair of "hostName" and ""
func splitPath(path string) (string, string) {
	idx := strings.IndexRune(path, ':')
	if idx == -1 {
		return "", path
	}
	return path[:idx], path[idx+1:]
}

// uploadFiles upload src file or directory to dest
func uploadFiles(src, dest string, client *sftp.Client) error {
	fi, err := client.Open(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	defer fi.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New(fmt.Sprintf("cannot access a file %s. %v\n", path, err))
		}
		//client.Create()
		fmt.Println("Path :", path, ", isDir :", info.IsDir())
		return nil
	})
}

func downloadFiles(src, dest string, client *sftp.Client) error {
	return nil
}
