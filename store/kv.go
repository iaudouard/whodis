package store

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	diskFileName = "store.wdb"
)

type KVStore struct {
	store map[string]string
}

func (kv KVStore) Get(key string) string {
	return kv.store[key]
}

func (kv *KVStore) Set(key string, value string) {
	kv.store[key] = value
}

func (kv *KVStore) Delete(key string) {
	delete(kv.store, key)
}

func (kv KVStore) toBytes() []byte {
	var data []byte
	for key, value := range kv.store {
		data = append(data, []byte(fmt.Sprintf("%s %s\n", key, value))...)
	}
	return data
}

func (kv KVStore) WriteToDisk() error {
	return saveToDisk(kv.toBytes(), diskFileName)
}

func (kv *KVStore) LoadFromDisk() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fp := filepath.Join(cwd, diskFileName)
	file, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	var key, value string
	for {
		_, err := fmt.Fscanf(file, "%s %s\n", &key, &value)
		if err != nil {
			break
		}
		kv.store[key] = value
	}

	return nil
}

func NewKVStore() KVStore {
	store := KVStore{
		store: make(map[string]string),
	}
	store.LoadFromDisk()
	return store
}
