// Package utils contains internal helper functions for commands.
package utils

import (
	"github.com/urfave/cli"
	"os"
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
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "myutils"
	app.Usage = "Usage.."
	app.Author = "zaccoding"
	app.Version = "0.0.1"
	return app
}

// GetDatabasePath returns a db directory.
func GetDatabasePath() (string, error) {
	//dir, err := os.Getwd()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	//fmt.Println("## database open.. ", (dir + "/myutilsdb"))
	return dir + "/myutilsdb", nil
}
