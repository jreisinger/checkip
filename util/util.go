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

func isOlderThanOneWeek(t time.Time) bool {
	return time.Now().Sub(t) > 7*24*time.Hour
}

func downloadFile(url string) (r io.ReadCloser, err error) {
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

// ExtractFile decompress r into filename. Supported compression formats are gz
// and tgz.
func ExtractFile(filename string, r io.ReadCloser, compressFmt string) error {
	switch compressFmt {
	case "gz":
		if err := extractGzFile(filename, r); err != nil {
			return err
		}
	case "tgz":
		if err := extractTgzFile(filename, r); err != nil {
			return err
		}
	default:
		return fmt.Errorf("don't know ho to extract a %s file", compressFmt)
	}
	return nil
}

func extractGzFile(outFilename string, r io.ReadCloser) error {
	defer r.Close() // let's close resp.Body

	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outFilename)
	if err != nil {
		return nil
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, gzipReader); err != nil {
		return err
	}

	return nil
}

func extractTgzFile(outFile string, r io.ReadCloser) error {
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

// Update updates file from url if the file is older than a week. If file does
// not exist it downloads and creates it. compressFmt is the compression format
// of the file to download; gz or tgz.
func Update(file, url string, compressFmt string) error {
	f, err := os.Stat(file)

	if os.IsNotExist(err) {
		r, err := downloadFile(url)
		if err != nil {
			return err
		}
		if err := ExtractFile(file, r, compressFmt); err != nil {
			return err
		}

		return nil // don't check ModTime if file does not exist
	}

	if isOlderThanOneWeek(f.ModTime()) {
		r, err := downloadFile(url)
		if err != nil {
			return err
		}
		if err := ExtractFile(file, r, compressFmt); err != nil {
			return err
		}
	}

	return nil
}
