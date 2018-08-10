package main

import (
	"encoding/json"
	"os"

	"github.com/pivotalservices/file-downloader-resource/config"
	"github.com/pivotalservices/file-downloader-resource/file"
	"github.com/pivotalservices/file-downloader-resource/types"
)

var VERSION = "0.0.0-dev"

func main() {
	if len(os.Args) < 2 {
		println("version: " + VERSION)
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fatal("creating destination", err)
	}

	var request types.InRequest
	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	configProvider, err := config.FromSource(request.Source)
	if err != nil {
		fatal("constructing config provider", err)
	}
	versionInfo, err := configProvider.GetVersionInfo(request.Version.Ref, request.Params.Product)
	if err != nil {
		fatal("getting version info", err)
	}

	fileProvider, err := file.FromSource(request.Source)
	if err != nil {
		fatal("constructing file provider", err)
	}
	if request.Params.Stemcell {
		err = fileProvider.DownloadFile(destination, versionInfo.StemcellProductPath(), versionInfo.StemcellVersion, versionInfo.StemcellFilePattern)
		if err != nil {
			fatal("downloading stemcell file", err)
		}
	} else {
		err = fileProvider.DownloadFile(destination, versionInfo.PivotalProduct, versionInfo.Version, versionInfo.FilePattern)
		if err != nil {
			fatal("downloading file", err)
		}
	}
	if request.Params.Stemcell {
		json.NewEncoder(os.Stdout).Encode(types.InResponse{
			Version: request.Version,
			Metadata: types.Metadata{
				{Name: "resource_version", Value: VERSION},
				{Name: "ref", Value: request.Version.Ref},
				{Name: "product", Value: "stemcells"},
				{Name: "product_version", Value: versionInfo.StemcellVersion},
				{Name: "file_pattern", Value: versionInfo.StemcellFilePattern},
			},
		})
	} else {
		json.NewEncoder(os.Stdout).Encode(types.InResponse{
			Version: request.Version,
			Metadata: types.Metadata{
				{Name: "resource_version", Value: VERSION},
				{Name: "ref", Value: request.Version.Ref},
				{Name: "product", Value: versionInfo.PivotalProduct},
				{Name: "product_version", Value: versionInfo.Version},
				{Name: "file_pattern", Value: versionInfo.FilePattern},
			},
		})
	}
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
