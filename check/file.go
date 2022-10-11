package check

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"
)

// getDbFilesPath returns OS specific absolute path to filename. On Unix it will
// be $HOME/.checkip/<filename>
func getDbFilesPath(filename string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(usr.HomeDir, ".checkip")

	// Make sure dir exists.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0750)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(dir, filename), nil
}

var mu sync.Mutex

// updateFile updates file from url if the file is older than a week. If file
// does not exist it downloads and creates it. compressFmt is the compression
// format of the file to download; gz or tgz. Empty string means no compression.
// Only one instance on updateFile can run at any given time.
func updateFile(file, url string, compressFmt string) error {
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Stat(file)

	if os.IsNotExist(err) {
		r, err := downloadFile(url)
		if err != nil {
			return err
		}
		if err := extractFile(file, r, compressFmt); err != nil {
			return err
		}

		return nil // don't check ModTime if file does not exist
	}

	if isOlderThanOneWeek(f.ModTime()) {
		r, err := downloadFile(url)
		if err != nil {
			return err
		}
		if err := extractFile(file, r, compressFmt); err != nil {
			return err
		}
	}

	return nil
}

func isOlderThanOneWeek(t time.Time) bool {
	return time.Since(t) > 7*24*time.Hour
}

func downloadFile(url string) (r io.ReadCloser, err error) {
	log.Printf("downloading %s", url)
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

// extractFile decompress r into filename. Supported compression formats are gz
// and tgz. Empty string means no compression.
func extractFile(filename string, r io.ReadCloser, compressFmt string) error {
	switch compressFmt {
	case "gz":
		if err := extractGzFile(filename, r); err != nil {
			return err
		}
	case "tgz":
		if err := extractTgzFile(filename, r); err != nil {
			return err
		}
	case "":
		if err := storeFile(filename, r); err != nil {
			return err
		}
	default:
		return fmt.Errorf("don't know ho to extract a %s file", compressFmt)
	}
	return nil
}

func storeFile(outFilename string, r io.ReadCloser) error {
	defer r.Close() // let's close resp.Body

	outFile, err := os.Create(outFilename)
	if err != nil {
		return nil
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, r); err != nil {
		return err
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

	for {
		tarHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if tarHeader.Typeflag == tar.TypeReg {
			if filepath.Base(tarHeader.Name) == filepath.Base(outFile) {
				if err := writeTarEntry(outFile, tarReader); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func writeTarEntry(outFile string, tarReader *tar.Reader) error {
	f, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, tarReader)
	return err
}
