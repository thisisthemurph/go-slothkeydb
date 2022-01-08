package repository

import (
	"errors"
	"fmt"
)

func makeUnknownKeyError(key string) error {
	return errors.New(fmt.Sprintf("A key of [%s] could not be found.", key))
}