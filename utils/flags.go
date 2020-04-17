// Package utils contains internal helper functions for commands.
package utils

import (
	"github.com/urfave/cli"
	"os/user"
	"path/filepath"
)

var (
	PathFlag = cli.StringFlag{
		Name:  "path",
		Usage: "path of config file.",
	}
	HostNameFlag = cli.StringFlag{
		Name:  "name, n",
		Usage: "name of the host.",
	}
	HostUserFlag = cli.StringFlag{
		Name:  "user, u",
		Usage: "host username for ssh.",
	}
	HostAddressFlag = cli.StringFlag{
		Name:  "address, a",
		Usage: "host address for ssh.",
	}
	HostPortFlag = cli.IntFlag{
		Name:  "port, p",
		Usage: "host port for ssh.",
		Value: 22,
	}
	HostPasswordFlag = cli.StringFlag{
		Name:  "password, pwd",
		Usage: "host password for ssh.",
	}
	HostPemPathFlag = cli.StringFlag{
		Name:  "keypath",
		Usage: "host key file path for ssh.",
	}
	HostDescriptionFLag = cli.StringFlag{
		Name:  "description, d",
		Usage: "description of host.",
	}
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "myutils"
	app.Usage = "Usage.."
	app.Author = "zaccoding"
	app.Version = "0.0.2"
	return app
}

// GetDatabasePath returns a db directory i.e workspace/myutilsdb
func GetDatabasePath() (string, error) {
	workspace, err := GetWorkspace()
	if err != nil {
		return "", err
	}
	return filepath.Join(workspace, "myutilsdb"), nil
}

// GetWorkspace returns myutils workspace i.e ~/myutils
func GetWorkspace() (string, error) {
	cu, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(cu.HomeDir, "myutils"), nil
}
