package datastore

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Datastore struct {
	db   *leveldb.DB
	path string
}

// NewDatastore returns a new datastore
func NewDatastore(path string, o *opt.Options) (*Datastore, error) {
	if path == "" {
		return nil, errors.New("invalid path")
	}

	var opts opt.Options
	if o != nil {
		opts = opt.Options(*o)
	}

	db, err := leveldb.OpenFile(path, &opts)
	if errors.IsCorrupted(err) {
		db, err = leveldb.RecoverFile(path, &opts)
		if err != nil {
			return nil, err
		}
	}

	return &Datastore{
		db:   db,
		path: path,
	}, nil
}

// Has returns true if its present in the store, otherwise false.
func (db *Datastore) Has(key []byte) (bool, error) {
	return db.db.Has(key, nil)
}

// Put inserts the given value into the store.
func (db *Datastore) Put(key []byte, value []byte) error {
	return db.db.Put(key, value, nil)
}

// Get returns a value given key.
func (db *Datastore) Get(key []byte) ([]byte, error) {
	val, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// NewIterator returns the entire key space.
func (db *Datastore) NewIterator() iterator.Iterator {
	return db.db.NewIterator(new(util.Range), nil)
}

// NewIteratorWithStart returns the key space given subset start.
func (db *Datastore) NewIteratorWithStart(start []byte) iterator.Iterator {
	return db.db.NewIterator(&util.Range{Start: start}, nil)
}

// NewIteratorWithPrefix returns the key space given prefix.
func (db *Datastore) NewIteratorWithPrefix(p []byte) iterator.Iterator {
	return db.db.NewIterator(util.BytesPrefix(p), nil)
}

// Delete removes the key from the store.
func (db *Datastore) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

// Path returns the path to datastore directory
func (db *Datastore) Path() string {
	return db.path
}

func (db *Datastore) Close() {
	db.db.Close()
}
