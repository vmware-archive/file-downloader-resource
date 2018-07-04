package main

import (
	"encoding/json"
	"os"

	"github.com/calebwashburn/file-downloader/config"
	"github.com/calebwashburn/file-downloader/types"
)

func main() {
	var request types.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	provider, err := config.FromSource(request.Source)
	if err != nil {
		fatal("constructing driver", err)
	}

	version, err := provider.LatestVersion()
	if err != nil {
		fatal("fetching version", err)
	}
	json.NewEncoder(os.Stdout).Encode(types.CheckResponse{*version})
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
