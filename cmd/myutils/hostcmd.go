package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zacscoding/myutils/datastore"
	"github.com/zacscoding/myutils/types"
	"github.com/zacscoding/myutils/utils"
	"io/ioutil"
	"log"
	"os"
)

var (
	hostFlags = []cli.Flag{
		utils.HostNameFlag,
		utils.HostUserFlag,
		utils.HostAddressFlag,
		utils.HostPortFlag,
		utils.HostPasswordFlag,
		utils.HostPemPathFlag,
	}

	hostCommand = cli.Command{
		Action:   showSubCommand,
		Name:     "host",
		Usage:    "manage hosts such as add | get | gets | update | delete",
		Category: "HOST COMMANDS",
		Subcommands: []cli.Command{
			{
				Name:   "import",
				Usage:  "Import hosts",
				Action: importHosts,
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

// display subcommands
func showSubCommand(ctx *cli.Context) error {
	return cli.ShowSubcommandHelp(ctx)
}

func importHosts(ctx *cli.Context) error {
	path := ctx.String(utils.PathFlag.Name)
	if path == "" {
		return errors.New("invalid config path")
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

	db, err := getDatastore()
	if err != nil {
		return err
	}
	defer db.Close()

	var failures []string
	for _, host := range hosts {
		if !host.HasCredentials() {
			failures = append(failures, host.Name)
			continue
		}

		key := getHostKey(host.Name)
		encoded, err := json.Marshal(host)
		if err != nil {
			failures = append(failures, host.Name)
			continue
		}

		err = db.Put(key, encoded)
		if err != nil {
			failures = append(failures, host.Name)
			continue
		}
	}

	log.Printf("import hosts result >> try : %d / failures : %d. >>>> %v\n", len(hosts), len(failures), failures)
	return nil
}

// AddHost save a given host to local datastore
func AddHost(host *types.Host) error {
	if !host.HasCredentials() {
		return errors.New("must have at least password or key path")
	}

	db, err := getDatastore()
	if err != nil {
		return err
	}
	defer db.Close()

	key := getHostKey(host.Name)
	encoded, err := json.Marshal(host)
	if err != nil {
		return err
	}

	err = db.Put(key, encoded)
	if err != nil {
		return err
	}
	log.Println("Success to save a host : ", string(encoded))
	return nil
}

// addHost save a host to local datastore
func addHost(ctx *cli.Context) error {
	host, err := parseHost(ctx)
	if err != nil {
		return err
	}
	return AddHost(host)
}

// GetHost returns a host given hostname
func GetHost(hostname string) (*types.Host, error) {
	db, err := getDatastore()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	val, err := db.Get(getHostKey(hostname))
	if err != nil {
		return nil, err
	}

	var h *types.Host
	err = json.Unmarshal(val, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// showHost display a host given query.
func showHost(ctx *cli.Context) error {
	host, err := parseHost(ctx)
	if err != nil {
		return err
	}
	if host.Name == "" {
		return errors.New("hostname must be not empty")
	}

	host, err = GetHost(host.Name)
	if err != nil {
		return nil
	}

	displayHost(host)
	return nil
}

// GetHosts returns list of hosts.
func GetHosts() ([]*types.Host, error) {
	db, err := getDatastore()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	itr := db.NewIteratorWithPrefix([]byte(types.HostPrefix))
	var hosts []*types.Host

	for itr.Next() {
		var h *types.Host
		err = json.Unmarshal(itr.Value(), &h)
		if err != nil {
			fmt.Println("Failed to unmarshal host.", err)
			continue
		}

		hosts = append(hosts, h)
	}
	return hosts, nil
}

// showHosts display all hosts from local store.
func showHosts(ctx *cli.Context) error {
	hosts, err := GetHosts()
	if err != nil {
		return nil
	}
	displayHost(hosts...)
	return nil
}

// UpdateHost update a given host into local stored.
func UpdateHost(host *types.Host) error {
	if host.Name == "" {
		return errors.New("invalid host name")
	}
	if !host.HasCredentials() {
		return errors.New("must have at least password or key path")
	}

	db, err := getDatastore()
	if err != nil {
		return nil
	}
	defer db.Close()

	has, err := db.Has(getHostKey(host.Name))
	if err != nil {
		return err
	}
	if !has {
		return errors.New("Not exist host with name " + host.Name)
	}

	err = AddHost(host)
	if err != nil {
		return err
	}
	log.Println("Success to update")
	return nil
}

// UpdateHost update a host parsed from cli into local stored.
func updateHost(ctx *cli.Context) error {
	host, err := parseHost(ctx)
	if err != nil {
		return err
	}
	return UpdateHost(host)
}

// DeleteHost delete a host with given name.
func DeleteHost(hostname string) error {
	db, err := getDatastore()
	if err != nil {
		return nil
	}
	defer db.Close()

	err = db.Delete(getHostKey(hostname))
	if err != nil {
		return err
	}
	return nil
}

// deleteHost delete a host parsed from cli.
func deleteHost(ctx *cli.Context) error {
	host, err := parseHost(ctx)
	if err != nil {
		return err
	}
	return DeleteHost(host.Name)
}

// Parse host from given cli.Context.
func parseHost(ctx *cli.Context) (*types.Host, error) {
	host := &types.Host{
		Name:     ctx.String("name"),
		User:     ctx.String("user"),
		Address:  ctx.String("address"),
		Port:     ctx.Int("port"),
		Password: ctx.String("password"),
		KeyPath:  ctx.String("keypath"),
	}

	return host, nil
}

// getDatastore returns a new datastore
func getDatastore() (*datastore.Datastore, error) {
	path, err := utils.GetDatastorePath()
	if err != nil {
		return nil, err
	}
	return datastore.NewDatastore(path, nil)
}

// getHostKey returns a key given host with prefix("host.")
func getHostKey(hostname string) []byte {
	return []byte(types.HostPrefix + hostname)
}

func displayHost(hosts ...*types.Host) {
	for i, h := range hosts {
		s, err := json.Marshal(h)
		if err != nil {
			log.Printf("%v -> %s\n", i+1, h.Name)
		} else {
			log.Printf("%v -> %s\n", i+1, string(s))
		}
	}
}
