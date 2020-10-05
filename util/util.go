package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/kylelemons/go-gypsy/yaml"
)

// GetConfigValue tries to get value for key first from an environment variable
// then from a configuration file at $HOME/.checkip.yaml
func GetConfigValue(key string) (string, error) {
	var v string
	if v = os.Getenv(key); v != "" {
		return v, nil
	}

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
		return "", err
	}
	return v, nil
}

func IsOlderThanOneWeek(t time.Time) bool {
	return time.Now().Sub(t) > 7*24*time.Hour
}

func DownloadFile(url string) (r io.ReadCloser, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Check the server response.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't download %v: %v", url, resp.Status)
	}

	return resp.Body, nil
}

func ExtractFile(outFile string, r io.ReadCloser) error {
	defer r.Close() // let's close resp.Body

	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for true {
		tarHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if tarHeader.Typeflag == tar.TypeReg {
			if filepath.Base(tarHeader.Name) == filepath.Base(outFile) {
				outFile, err := os.Create(outFile)
				if err != nil {
					return err
				}
				defer outFile.Close()
				if _, err := io.Copy(outFile, tarReader); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
