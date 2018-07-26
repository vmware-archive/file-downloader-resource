package types

type Version struct {
	Ref string `json:"ref"`
}
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type CheckResponse []Version

type InRequest struct {
	Source  Source   `json:"source"`
	Version Version  `json:"version"`
	Params  InParams `json:"params"`
}

type InParams struct {
	Product  string `json:"product"`
	Stemcell bool   `json:"stemcell"`
}

type InResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type OutRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type OutResponse struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Source struct {
	ConfigProvider       ConfigProviderEnum `json:"config_provider"`
	FileProvider         FileProviderEnum   `json:"file_provider"`
	VersionRoot          string             `json:"version_root"`
	URI                  string             `json:"uri"`
	Branch               string             `json:"branch"`
	PrivateKey           string             `json:"private_key"`
	Username             string             `json:"username"`
	Password             string             `json:"password"`
	PivnetToken          string             `json:"pivnet_token"`
	Bucket               string             `json:"bucket"`
	AccessKeyID          string             `json:"access_key_id"`
	SecretAccessKey      string             `json:"secret_access_key"`
	RegionName           string             `json:"region_name"`
	Endpoint             string             `json:"endpoint"`
	DisableSSL           bool               `json:"disable_ssl"`
	SkipSSLVerification  bool               `json:"skip_ssl_verification"`
	ServerSideEncryption string             `json:"server_side_encryption"`
	UseV2Signing         bool               `json:"use_v2_signing"`
}

type ConfigProviderEnum string

const (
	ConfigProviderUnspecified ConfigProviderEnum = ""
	ConfigProviderGit         ConfigProviderEnum = "git"
)

type FileProviderEnum string

const (
	FileProviderUnspecified FileProviderEnum = ""
	FileProviderPivnet      FileProviderEnum = "pivnet"
	FileProviderS3          FileProviderEnum = "s3"
)

type VersionInfo struct {
	Version             string `yaml:"version"`
	PivotalProduct      string `yaml:"product"`
	FilePattern         string `yaml:"file_pattern"`
	StemcellVersion     string `yaml:"stemcell_version"`
	StemcellFilePattern string `yaml:"stemcell_file_pattern"`
}
