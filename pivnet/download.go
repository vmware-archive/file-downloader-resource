package pivnet

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"os"
	"path/filepath"
	"strings"

	pivnetapi "github.com/pivotal-cf/go-pivnet"

	"github.com/pivotal-cf/go-pivnet/logshim"
)

type Downloader struct {
	cacheDir       string
	client         pivnetapi.Client
	progressWriter io.Writer
}

func NewDownloader(token, cacheDir string) (*Downloader, error) {
	config := pivnetapi.ClientConfig{
		Host:      pivnetapi.DefaultHost,
		Token:     token,
		UserAgent: "file-downloader",
	}
	stdoutLogger := log.New(os.Stdout, "", log.LstdFlags)
	stderrLogger := log.New(os.Stderr, "", log.LstdFlags)

	verbose := false
	logger := logshim.NewLogShim(stdoutLogger, stderrLogger, verbose)

	client := pivnetapi.NewClient(config, logger)
	return &Downloader{
		client:         client,
		progressWriter: os.Stdout,
		cacheDir:       cacheDir,
	}, nil
}

func (d *Downloader) Download(targetDirectory, slug, version string, pattern string) error {

	err := os.MkdirAll(d.cacheDir, os.ModePerm)
	if err != nil {
		return err
	}
	releases, err := d.client.Releases.List(slug)
	if err != nil {
		return err
	}

	for _, release := range releases {
		if release.Version == version {
			productFiles, err := d.client.ProductFiles.ListForRelease(slug, release.ID)
			if err != nil {
				return err
			}
			return d.downloadFiles(targetDirectory, pattern, productFiles, slug, release.ID)

		}
	}

	return fmt.Errorf("Release Version %s of product %s not found", version, slug)

}

func (d *Downloader) downloadFiles(
	targetDirectory string,
	pattern string,
	productFiles []pivnetapi.ProductFile,
	productSlug string,
	releaseID int,
) error {

	filtered := productFiles

	// If globs were not provided, download everything without filtering.
	if pattern != "" {
		var err error
		filtered, err = productFileKeysByGlobs(productFiles, pattern)
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(targetDirectory, os.ModePerm); err != nil {
		return err
	}

	for _, pf := range filtered {
		parts := strings.Split(pf.AWSObjectKey, "/")
		fileName := parts[len(parts)-1]
		targetFile := filepath.Join(targetDirectory, fileName)
		cacheFile := filepath.Join(d.cacheDir, fileName)
		if _, err := os.Stat(cacheFile); err == nil {
			fileBytes, err := ioutil.ReadFile(cacheFile)
			if err != nil {
				return err
			}
			sum := sha256.Sum256(fileBytes)
			if pf.SHA256 != fmt.Sprintf("%x", sum) {
				fmt.Println(fmt.Sprintf("Removing corrupt cache %s", fileName))
				err = os.Remove(cacheFile)
				if err != nil {
					return err
				}
			}
		}

		if _, err := os.Stat(cacheFile); err != nil {
			fmt.Println(fmt.Sprintf("Not found in cache %s, downloading", fileName))
			file, err := os.Create(cacheFile)
			if err != nil {
				return err
			}

			err = d.client.ProductFiles.DownloadForRelease(file, productSlug, releaseID, pf.ID, d.progressWriter)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(fmt.Sprintf("Found in cache %s", fileName))
		}

		if err := copyFileContents(cacheFile, targetFile); err != nil {
			return err
		}
	}
	return nil
}

func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()

}

func productFileKeysByGlobs(
	productFiles []pivnetapi.ProductFile,
	pattern string,
) ([]pivnetapi.ProductFile, error) {

	filtered := []pivnetapi.ProductFile{}

	for _, p := range productFiles {
		parts := strings.Split(p.AWSObjectKey, "/")
		fileName := parts[len(parts)-1]

		matched, err := filepath.Match(pattern, fileName)
		if err != nil {
			return nil, err
		}

		if matched {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 && pattern != "" {
		return nil, fmt.Errorf("no match for pattern: '%s'", pattern)
	}

	return filtered, nil
}
