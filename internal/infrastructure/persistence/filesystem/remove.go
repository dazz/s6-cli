package filesystem

import (
	"errors"
	"fmt"
	"os"
)

func (fs *Filesystem) RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}

func (fs *Filesystem) RemoveFile(path string) error {
	if err := os.Remove(path); err != nil {
		return errors.New(fmt.Sprintf("Error removing script file: %s", err))
	}
	return nil
}
