package misc

import (
	"errors"
	"os"
)

func IsFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}
