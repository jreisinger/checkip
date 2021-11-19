package check

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GetConfigValue tries to get value for key first from an environment variable
// then from a configuration file at $HOME/.checkip.yaml.
var GetConfigValue = func(key string) (string, error) {
	var v string

	// Try to get the key from environment.
	if v = os.Getenv(key); v != "" {
		return v, nil
	}

	// Try to get the key from the config file.
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	confFile := filepath.Join(usr.HomeDir, ".checkip.yaml")
	buf, err := ioutil.ReadFile(confFile)
	if err != nil {
		return "", err
	}

	var data map[string]string
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		return "", err
	}

	if v, ok := data[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("%s not found in %s", key, confFile)
}
