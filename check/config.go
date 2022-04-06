package check

import (
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// GetConfigValue tries to get value for key first from an environment variable
// then from a configuration file at $HOME/.checkip.yaml. If value is not found
// an empty string and nil is returned (i.e. it's not considered an error).
var GetConfigValue = func(key string) (string, error) {
	// Try to get the key from environment.
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	// Try to get the key from the config file.
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	confFile := filepath.Join(usr.HomeDir, ".checkip.yaml")
	buf, err := os.ReadFile(confFile)
	if err != nil {
		return "", nil // non existent file is not considered an error
	}

	var data map[string]string
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		return "", err
	}

	if value, ok := data[key]; ok {
		return value, nil
	}

	return "", nil // value not found
}
