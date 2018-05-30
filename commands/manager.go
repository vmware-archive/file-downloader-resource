package commands

type manager struct {
	DownloadPivnet DownloadPivnet `command:"download-pivnet" description:"downloads specific file version from pivnet"`
}

var Manager manager
