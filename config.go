package checkip

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/kylelemons/go-gypsy/yaml"
)

// GetConfigValue tries to get value for key first from an environment variable
// then from a configuration file at $HOME/.checkip.yaml
func GetConfigValue(key string) (string, error) {
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
	cfg, err := yaml.ReadFile(confFile)
	if err != nil {
		return "", err
	}
	v, err = cfg.Get(key)
	if err != nil {
		return "", fmt.Errorf("%s not found in %s", key, confFile)
	}

	return v, nil
}
