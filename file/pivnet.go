package file

import (
	"crypto/sha256"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"log"
	"os"

	"github.com/fatih/color"
	pivnetapi "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
)

type PivnetProvider struct {
	client         pivnetapi.Client
	progressWriter io.Writer
	logger         *logshim.LogShim
}

func NewPivnetProvider(token string) (Provider, error) {
	color.NoColor = false
	logWriter := os.Stderr
	logger := log.New(logWriter, "", log.LstdFlags)
	config := pivnetapi.ClientConfig{
		Host:      pivnetapi.DefaultHost,
		Token:     token,
		UserAgent: "file-downloader",
	}
	ls := logshim.NewLogShim(logger, logger, false)
	client := pivnetapi.NewClient(config, ls)
	return &PivnetProvider{
		client:         client,
		progressWriter: os.Stderr,
		logger:         ls,
	}, nil
}

//DownloadFile - Downloads file based on version info
func (p *PivnetProvider) DownloadFile(targetDirectory, productSlug, version, pattern string) error {

	releases, err := p.client.Releases.List(productSlug)
	if err != nil {
		return err
	}

	for _, release := range releases {
		if release.Version == version {
			productFiles, err := p.client.ProductFiles.ListForRelease(productSlug, release.ID)
			if err != nil {
				return err
			}
			err = p.client.EULA.Accept(productSlug, release.ID)
			if err != nil {
				return err
			}
			return p.downloadFiles(targetDirectory, pattern, productFiles, productSlug, release.ID)

		}
	}
	return fmt.Errorf("Release Version %s of product %s not found", version, productSlug)
}

func (p *PivnetProvider) downloadFiles(
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
		file, err := os.Create(targetFile)
		if err != nil {
			return err
		}
		err = p.client.ProductFiles.DownloadForRelease(file, productSlug, releaseID, pf.ID, p.progressWriter)
		if err != nil {
			return err
		}
	}
	return nil
}

func sumFile(filepath string) (string, error) {
	fileToSum, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer fileToSum.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, fileToSum)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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
