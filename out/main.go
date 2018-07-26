package main

import (
	"encoding/json"
	"os"

	"github.com/pivotalservices/file-downloader-resource/types"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <source>")
		os.Exit(1)
	}

	var request types.OutRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	outVersion := request.Version

	json.NewEncoder(os.Stdout).Encode(types.OutResponse{
		Version: outVersion,
	})
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
