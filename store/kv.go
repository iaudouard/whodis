package store

import (
	"encoding/json"
	"fmt"
	"os"
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

func (kv KVStore) WriteToDisk() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := json.MarshalIndent(kv.store, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/backup.json", path), file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewKVStore() KVStore {
	store := KVStore{
		store: make(map[string]string),
	}
	store.Set("hello", "world")
	return store
}
