package cfg

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
)

func Dir() string {
	home := os.Getenv("JABBA_HOME")
	if home != "" {
		return filepath.Clean(home)
	}
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, ".jabba")
}

func Index() string {
	registry := os.Getenv("JABBA_INDEX")
	if registry == "" {
		// registry = "https://github.com/shyiko/jabba/raw/master/index.json"
		registry = "https://gitee.com/molihuan/jabba/raw/master/index.json"
	}
	return registry
}
