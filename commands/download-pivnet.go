package commands

import (
	"io/ioutil"

	"github.com/calebwashburn/file-downloader/pivnet"
	"gopkg.in/yaml.v2"
)

type DownloadPivnet struct {
	Token       string `long:"token" env:"PIVNET_TOKEN" description:"pivnet api token" required:"true"`
	DownloadDir string `long:"download-dir" env:"PIVNET_DOWNLOAD_DIR" description:"directory to download files to" required:"true"`
	StemcellDir string `long:"stemcell-dir" env:"PIVNET_STEMCELL_DIR" description:"directory to download stemcell" default:"stemcells"`
	CacheDir    string `long:"cache-dir" description:"directory for file cache" default:".pivotal-cache"`
	ConfigFile  string `long:"config-file" description:"file that contains configuration of file and versions to download" required:"true"`
}

type ProductConfig struct {
	Version             string `yaml:"version"`
	Product             string `yaml:"product"`
	FilePattern         string `yaml:"file_pattern"`
	StemcellVersion     string `yaml:"stemcell_version"`
	StemcellFilePattern string `yaml:"stemcell_file_pattern"`
}

func readProductConfig(file string) (*ProductConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	productConfig := &ProductConfig{}
	err = yaml.Unmarshal(data, &productConfig)
	if err != nil {
		return nil, err
	}
	return productConfig, nil
}

//Execute - downloads files from pivnet
func (d *DownloadPivnet) Execute([]string) error {
	productConfig, err := readProductConfig(d.ConfigFile)
	if err != nil {
		return err
	}
	downloader, err := pivnet.NewDownloader(d.Token, d.CacheDir)
	if err != nil {
		return err
	}

	if err := downloader.Download(d.DownloadDir, productConfig.Product, productConfig.Version, productConfig.FilePattern); err != nil {
		return err
	}
	if productConfig.StemcellVersion != "" {
		if err := downloader.Download(d.StemcellDir, "stemcells", productConfig.StemcellVersion, productConfig.StemcellFilePattern); err != nil {
			return err
		}
	}
	return nil
}
