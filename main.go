package main

import (
	"slothkeydb/store"
)

// Opens the log or creats if one doesn't exists
func Open(path string) (store.Store, error) {
	return store.Open(path)
}

