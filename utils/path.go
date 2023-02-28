package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var (
	DefaultFolder     = getDefaultFolder()
	DefaultLoggerFile = "logger.json"
)

func getDefaultFolder() string {
	defaultFolder := ".logger_go"
	home := os.Getenv("HOME")
	if home == "" {
		usr, err := user.Current()
		if err != nil {
			panic(fmt.Sprintf("cannot get current user: %s", err))
		}
		home = usr.HomeDir
	}
	return filepath.Join(home, defaultFolder)
}

func MakeAllDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}
