package main

import (
	"encoding/json"
	"errors"
	"github.com/urfave/cli"
	"github.com/zacscoding/myutils/host"
	"github.com/zacscoding/myutils/types"
	"github.com/zacscoding/myutils/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var (
	hostFlags = []cli.Flag{
		utils.HostNameFlag,
		utils.HostUserFlag,
		utils.HostAddressFlag,
		utils.HostPortFlag,
		utils.HostPasswordFlag,
		utils.HostPemPathFlag,
		utils.HostDescriptionFLag,
	}

	hostCommand = cli.Command{
		Action:   ShowSubCommand,
		Name:     "host",
		Usage:    "manage hosts such as add | get | gets | update | delete",
		Category: "HOST COMMANDS",
		Subcommands: []cli.Command{
			{
				Name:   "import",
				Usage:  "Import hosts json file to local store",
				Action: importHosts,
				Flags: []cli.Flag{
					utils.PathFlag,
				},
			},
			{
				Name:   "export",
				Usage:  "Export to hosts json file from local store",
				Action: exportHosts,
				Flags: []cli.Flag{
					utils.PathFlag,
				},
			},
			{
				Name:   "add",
				Usage:  "Adds a host",
				Action: addHost,
				Flags:  hostFlags,
			},
			{
				Name:   "get",
				Usage:  "Get a host",
				Action: showHost,
				Flags:  hostFlags,
			},
			{
				Name:   "gets",
				Usage:  "Get hosts",
				Action: showHosts,
				Flags:  hostFlags,
			},
			{
				Name:   "update",
				Usage:  "Update a host",
				Action: updateHost,
				Flags:  hostFlags,
			},
			{
				Name:   "delete",
				Usage:  "Delete a host",
				Action: deleteHost,
				Flags:  hostFlags,
			},
		},
	}
)

// exportHosts export hosts data in local store to json file.
func exportHosts(ctx *cli.Context) error {
	path := ctx.String(utils.PathFlag.Name)
	if path == "" {
		return errors.New(`path must not be ""`)
	}

	// 1) exist file
	// 	1-1) directory
	//		=> use default filename
	// 	1-2) file
	//		=> return error
	// 2) not exist file
	// 2-1) create a new file
	fi, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else {
		if !fi.IsDir() {
			return errors.New("already exist file :" + path)
		}
		log.Println("use default filename : hosts.json")
		path = filepath.Join(path, "hosts.json")
	}

	hosts, err := host.GetHosts(app.db)
	if err != nil {
		return err
	}
	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].Name < hosts[j].Name
	})

	b, err := json.MarshalIndent(hosts, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}
	log.Println("success to export hosts. destination :", path)
	return nil
}

// importHosts import hosts data from json.
func importHosts(ctx *cli.Context) error {
	path := ctx.String(utils.PathFlag.Name)
	if path == "" {
		return errors.New(`path must not be ""`)
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	readBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var hosts []*types.Host
	err = json.Unmarshal(readBytes, &hosts)
	if err != nil {
		return err
	}

	var failures []string
	for _, h := range hosts {
		err = host.AddHost(app.db, h)

		if err != nil {
			failures = append(failures, h.Name)
			continue
		}
	}

	log.Printf("import hosts result >> try : %d / failures : %d. >>>> %v\n", len(hosts), len(failures), failures)
	return nil
}

// addHost save a host to local db
func addHost(ctx *cli.Context) error {
	h, err := parseHost(ctx)
	if err != nil {
		return err
	}
	return host.AddHost(app.db, h)
}

// showHost display a host given query.
func showHost(ctx *cli.Context) error {
	h, err := parseHost(ctx)
	if err != nil {
		return err
	}
	if h.Name == "" {
		return errors.New("hostname must be not empty")
	}

	h, err = host.GetHost(app.db, h.Name)
	if err != nil {
		return nil
	}
	displayHost(h)
	return nil
}

// showHosts display all hosts from local store.
func showHosts(ctx *cli.Context) error {
	hosts, err := host.GetHosts(app.db)
	if err != nil {
		return nil
	}
	displayHost(hosts...)
	return nil
}

// updateHost update a host parsed from cli into local stored.
func updateHost(ctx *cli.Context) error {
	h, err := parseHost(ctx)
	if err != nil {
		return err
	}
	return host.UpdateHost(app.db, h)
}

// deleteHost delete a host parsed from cli.
func deleteHost(ctx *cli.Context) error {
	h, err := parseHost(ctx)
	if err != nil {
		return err
	}
	return host.DeleteHost(app.db, h.Name)
}

// Parse host from given cli.Context.
func parseHost(ctx *cli.Context) (*types.Host, error) {
	host := &types.Host{
		Name:        ctx.String("name"),
		User:        ctx.String("user"),
		Address:     ctx.String("address"),
		Port:        ctx.Int("port"),
		Password:    ctx.String("password"),
		KeyPath:     ctx.String("keypath"),
		Description: ctx.String("description"),
	}
	return host, nil
}

// displayHost show all hosts to console.
func displayHost(hosts ...*types.Host) {
	if hosts == nil || len(hosts) == 0 {
		log.Printf("> empty hosts in local store")
		return
	}

	for i, h := range hosts {
		s, err := json.Marshal(h)
		if err != nil {
			log.Printf("%v -> %s\n", i+1, h.Name)
		} else {
			log.Printf("%v -> %s\n", i+1, string(s))
		}
	}
}
