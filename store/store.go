package store

import (
	"slothkeydb/repository"
)

type Store interface {
	Add(key, value string)
	Get(key string) (string, error)
	Delete(key string)
}

type store struct {
	log repository.LogReaderWriter
}

// Add a new key-value pair to the store
func (s store) Add(key, value string) {
	s.log.Write(key, value, repository.LET_Live)
}

// Retrieve the value of the given key
func (s store) Get(key string) (string, error) {
	return s.log.Read(key)
}

// Delete a given key from the store
func (s store) Delete(key string) {
	s.log.Write(key, "", repository.LET_Deleted)
}

// Opens a store - creates if not exists
func Open(path string) (Store, error) {
	repoLog, err := repository.MakeLogFile(path)
	if err != nil {
		return nil, err
	}

	return store{repoLog}, nil
}