# File Downloader Resource

A resource that pulls configuration from git to allow externalizing what files to download for a given provider vs having to reconfigure resource parameters on traditional [pivnet resource](https://github.com/pivotal-cf/pivnet-resource)

## Maintainer

* [Caleb Washburn](https://github.com/calebwashburn)

## Support

file downloader resource is a community supported concourse resource.  Opening issues for questions, feature requests and/or bugs is the best path to getting "support".  We strive to be active in keeping this tool working and meeting your needs in a timely fashion.

## Source Configuration

* `config_provider`: *Optional. Default `git`.* The provider used to pull configuration

* `file_provider`: *Optional. Default `pivnet`.* The provider used to download files from

### Configuration Provider

There is 1 supported configuration provider

#### `git` Configuration provider

The `git` provider will retrieve new configuration when a commit to repository occurs

* `uri`: *Required.* The repository URL.

* `branch`: *Required.* The branch the file lives on.

* `version_root`: *Required* The path of where to find configuration files

* `private_key`: *Optional.* The SSH private key to use when pulling from/pushing to to the repository.

* `username`: *Optional.* Username for HTTP(S) auth when pulling/pushing.
   This is needed when only HTTP/HTTPS protocol for git is available (which does not support private key auth)
   and auth is required.

* `password`: *Optional.* Password for HTTP(S) auth when pulling/pushing.

* `path`: *Optional.* Path to look for changes in

### File Provider

There is 2 supported file provider(2)

### `pivnet` provider

The `pivnet` provider works by downloading files based on configuration.

* `pivnet_token`: *Required.* Token used to authenticate to pivnet (network.pivotal.io)

### `s3` provider

The `s3` provider works by downloading file from s3

* `bucket`: *Required.* The name of the bucket.

* `access_key_id`: *Optional.* The AWS access key to use when accessing the
  bucket.

* `secret_access_key`: *Optional.* The AWS secret key to use when accessing
  the bucket.

* `region_name`: *Optional.* The region the bucket is in. Defaults to
  `us-east-1`.

* `endpoint`: *Optional.* Custom endpoint for using S3 compatible provider.

* `disable_ssl`: *Optional.* Disable SSL for the endpoint, useful for S3
  compatible providers without SSL.

* `skip_ssl_verification`: *Optional.* Skip SSL verification for S3 endpoint. Useful for S3 compatible providers using self-signed SSL certificates.

* `use_v2_signing`: *Optional.* Use signature v2 signing, useful for S3 compatible providers that do not support v4.

### Sample Configuration

``` yaml
version: 2.1.5
product: elastic-runtime
file_pattern: srt-*.pivotal
stemcell_version: "3541.25"
stemcell_file_pattern: "bosh-stemcell-*-vsphere-esxi-ubuntu-trusty-go_agent.tgz"
```

* `version`: *Required.* Version of file to download

* `product`: *Required.* Product name (for pivnet this is the product slug)

* `file_pattern`: *Required.* File Pattern (for pivnet this is the product glob)

* `stemcell_version`: *Optional.* Version of stemcell to download for a given product

* `stemcell_file_pattern`: *Optional.* Stemcell File Pattern (for pivnet this is the product glob)


### Example

With the following resource configuration:

``` yaml
resource_types:
- name: file-downloader
  type: docker-image
  source:
    repository: pivotalservices/file-downloader-resource
```

Using Pivnet file provider
``` yaml
resources:
- name: pivnet-files
  type: file-downloader
  source:
    config_provider: git
    version_root: {{folder_path_in_git_repo}}
    uri: git@github.com:pivotalservices/your_repo.git
    private_key: {{git_private_key}}
    branch: master
    file_provider: pivnet
    pivnet_token: {{pivnet_token}}
```

Using s3 file provider
``` yaml
resources:
- name: pivnet-files
  type: file-downloader
  source:
    config_provider: git
    version_root: {{folder_path_in_git_repo}}
    uri: git@github.com:pivotalservices/your_repo.git
    private_key: {{git_private_key}}
    branch: master
    file_provider: s3
    bucket: {{s3_bucket}}
    access_key_id: {{s3_access_key}}
    secret_access_key: {{s3_secret_access_key}}
```

To retrieve files for a opsman product with a `get`:

``` yaml
plan:
- get: image
  resource: pivnet-files
  params:
    product: opsman
```

Opsman product configuration in git

``` yaml
version: 2.1.3
product: ops-manager
file_pattern: pcf-vsphere-*.ova
```


Or, to pull .pivotal and stemcells for cf product with `get`:

``` yaml
plan:
- aggregate:
  - get: cf
    resource: pivnet-files
    params:
      product: cf
  - get: cf-stemcell
    resource: pivnet-files
    params:
      product: cf
      stemcell: true
```

cf product configuration in git

``` yaml
version: 2.1.5
product: elastic-runtime
file_pattern: cf-*.pivotal
stemcell_version: "3541.25"
stemcell_file_pattern: "bosh-stemcell-*-vsphere-esxi-ubuntu-trusty-go_agent.tgz"
```

## Behavior

### `check`: Report the current version based on configuration provider.

Detects new versions

### `in`: Provide the file based on get parameters

Based on the `version` from the check will parse the configuration file for given product

#### Parameters

* `product`: *Required.* name of .yml file in `version_root`

* `stemcell`: *optional. default false* true/false indicates where to download stemcell based on `stemcell_version` and `stemcell_file_pattern` in the <product>.yml file

### `out`: No-op command

### Running the tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t pivotalservices/file-downloader-resource .
```

### Contributing

Please make all pull requests to the `master` branch and ensure tests pass
locally.
