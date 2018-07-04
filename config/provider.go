package config

import (
	"fmt"

	"github.com/calebwashburn/file-downloader/types"
)

//go:generate counterfeiter -o fakes/fake_provider.go provider.go Provider

// Provider - defines the interface for how to fetch configuration
type Provider interface {
	LatestVersion() (*types.Version, error)
	GetVersionInfo(revision, productName string) (*types.VersionInfo, error)
}

// FromSource - factory to return appropriate driver based on configuration
func FromSource(source types.Source) (Provider, error) {

	switch source.ConfigProvider {
	case types.ConfigProviderUnspecified, types.ConfigProviderGit:

		return &GitProvider{
			VersionRoot: source.VersionRoot,
			URI:         source.URI,
			Branch:      source.Branch,
			PrivateKey:  source.PrivateKey,
			Username:    source.Username,
			Password:    source.Password,
		}, nil

	default:
		return nil, fmt.Errorf("unknown provider: %s", source.ConfigProvider)
	}
}
