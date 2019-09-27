// host is management host from data store.
package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zacscoding/myutils/db"
	"github.com/zacscoding/myutils/types"
	"log"
)

// AddHost save a given host into local db
func AddHost(db *db.Database, host *types.Host) error {
	if !host.HasCredentials() {
		hostStr := ""
		b, err := json.Marshal(host)
		if err == nil {
			hostStr = string(b)
		}
		return errors.New("must have at least password or key path :" + hostStr)
	}

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

// GetHost returns a host given hostname
func GetHost(db *db.Database, hostname string) (*types.Host, error) {
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

// GetHosts returns list of hosts from db
func GetHosts(db *db.Database) ([]*types.Host, error) {
	itr := db.NewIteratorWithPrefix([]byte(types.HostPrefix))
	var hosts []*types.Host

	for itr.Next() {
		var h *types.Host
		if err := json.Unmarshal(itr.Value(), &h); err != nil {
			fmt.Println("Failed to unmarshal host.", err)
			continue
		}

		hosts = append(hosts, h)
	}
	return hosts, nil
}

// UpdateHost update a given host into local stored.
func UpdateHost(db *db.Database, h *types.Host) error {
	has, err := db.Has(getHostKey(h.Name))
	if err != nil {
		return err
	}
	if !has {
		return errors.New("Not exist host with name " + h.Name)
	}

	err = AddHost(db, h)
	if err == nil {
		log.Println("Success to update")
	}
	return err
}

// DeleteHost delete a host with given name.
func DeleteHost(db *db.Database, hostname string) error {
	return db.Delete(getHostKey(hostname))
}

// getHostKey returns a key given host with prefix("host.")
func getHostKey(hostname string) []byte {
	return []byte(types.HostPrefix + hostname)
}
