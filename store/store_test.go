package store_test

import (
	"log"
	"slothkeydb/store"
	"testing"
)

const STORE_PATH string = "..\\teststore.slo"

func TestGet(t *testing.T) {
	myStore, _ := store.Open(STORE_PATH)
	myStore.Add("name", "William")
	myStore.Add("age", "twenty")
	myStore.Add("name", "Billy")
	myStore.Add("sex", "m")
	
	t.Run("adding and getting", func(t *testing.T) {
		key := "name"
		want := "Billy"
		got, _ := myStore.Get(key)

		if got != want {
			log.Fatalf(`Got "%s" expected "%s"`, got, want)
		}	
	})
	
	t.Run("unknown key", func(t *testing.T) {
		got, err := myStore.Get("UNKNOWN KEY")
		
		if got != "" {
			log.Fatalf("Expected empty string, got %s", got)
		}
		
		if err == nil {
			log.Fatal("Expected error response, no error returned")
		}
	})
}

func TestDelete(t *testing.T) {
	myStore, _ := store.Open(STORE_PATH)
	myStore.Add("name", "William")
	myStore.Add("name", "Billy")
	myStore.Add("name", "Michael")
	myStore.Add("name", "Mike")
	
	t.Run("deleting a key", func(t *testing.T) {
		key := "name"

		myStore.Delete(key)
		
		want := ""
		got, err := myStore.Get(key)
		
		if (got != want) {
			log.Fatalf("Expected empty string, got %s", got)
		}

		if err == nil {
			log.Fatal("Expected error response, no error returned")
		}
	})
}
