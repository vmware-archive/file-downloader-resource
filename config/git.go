package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/pivotalservices/file-downloader-resource/types"
)

var gitRepoDir string
var privateKeyPath string
var netRcPath string

var ErrEncryptedKey = errors.New("private keys with passphrases are not supported")

func init() {
	gitRepoDir = filepath.Join(os.TempDir(), "file-downloader-git-repo")
	privateKeyPath = filepath.Join(os.TempDir(), "private-key")
	netRcPath = filepath.Join(os.Getenv("HOME"), ".netrc")
}

type GitProvider struct {
	VersionRoot string
	URI         string
	Branch      string
	PrivateKey  string
	Username    string
	Password    string
	Depth       string
	Path        string
}

//GetVersionInfo - Check returns version of git resource
func (provider *GitProvider) GetVersionInfo(revision, productName string) (*types.VersionInfo, error) {
	err := provider.setUpAuth()
	if err != nil {
		return nil, err
	}

	err = provider.setUpRepo(revision)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(path.Join(gitRepoDir, provider.VersionRoot, fmt.Sprintf("%s.yml", productName)))
	if err != nil {
		return nil, err
	}
	versionInfo := types.VersionInfo{}

	err = yaml.Unmarshal(bytes, &versionInfo)

	return &versionInfo, err
}

//LatestVersion - Check returns version of git resource
func (provider *GitProvider) LatestVersion() (*types.Version, error) {
	err := provider.setUpAuth()
	if err != nil {
		return nil, err
	}

	err = provider.setUpRepo("HEAD")
	if err != nil {
		return nil, err
	}

	var gitVersion *exec.Cmd
	if len(provider.Path) > 0 {
		gitVersion = exec.Command("git", "log", "--format='%H'", "--first-parent", "-1", "--", provider.Path)
	} else {
		gitVersion = exec.Command("git", "log", "--format='%H'", "--first-parent", "-1")
	}
	gitVersion.Dir = gitRepoDir
	gitVersion.Stderr = os.Stderr
	out, err := gitVersion.Output()
	if err != nil {
		return nil, err
	}
	return &types.Version{Ref: strings.Replace(strings.Replace(string(out), "\n", "", -1), "'", "", -1)}, nil
}

func (provider *GitProvider) setUpRepo(revision string) error {
	_, err := os.Stat(gitRepoDir)
	if err != nil {
		gitClone := exec.Command("git", "clone", provider.URI, "--branch", provider.Branch)
		if len(provider.Depth) > 0 {
			gitClone.Args = append(gitClone.Args, "--depth", provider.Depth)
		}
		gitClone.Args = append(gitClone.Args, "--single-branch", gitRepoDir)
		gitClone.Stdout = os.Stderr
		gitClone.Stderr = os.Stderr
		if err := gitClone.Run(); err != nil {
			return err
		}
	} else {
		gitFetch := exec.Command("git", "fetch", "origin", provider.Branch)
		gitFetch.Dir = gitRepoDir
		gitFetch.Stdout = os.Stderr
		gitFetch.Stderr = os.Stderr
		if err := gitFetch.Run(); err != nil {
			return err
		}
	}

	gitReset := exec.Command("git", "reset", "--hard", "origin/"+provider.Branch)
	gitReset.Dir = gitRepoDir
	gitReset.Stdout = os.Stderr
	gitReset.Stderr = os.Stderr
	if err := gitReset.Run(); err != nil {
		return err
	}

	gitCheckout := exec.Command("git", "checkout", "-q", revision)
	gitCheckout.Dir = gitRepoDir
	gitCheckout.Stdout = os.Stderr
	gitCheckout.Stderr = os.Stderr
	if err := gitCheckout.Run(); err != nil {
		return err
	}

	return nil
}

func (provider *GitProvider) setUpAuth() error {
	_, err := os.Stat(netRcPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err := os.Remove(netRcPath)
		if err != nil {
			return err
		}
	}

	if len(provider.PrivateKey) > 0 {
		err := provider.setUpKey()
		if err != nil {
			return err
		}
	}

	if len(provider.Username) > 0 && len(provider.Password) > 0 {
		err := provider.setUpUsernamePassword()
		if err != nil {
			return err
		}
	}

	return nil
}

func (provider *GitProvider) setUpKey() error {
	if strings.Contains(provider.PrivateKey, "ENCRYPTED") {
		return ErrEncryptedKey
	}

	_, err := os.Stat(privateKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile(privateKeyPath, []byte(provider.PrivateKey), 0600)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return os.Setenv("GIT_SSH_COMMAND", "ssh -o StrictHostKeyChecking=no -i "+privateKeyPath)
}

func (provider *GitProvider) setUpUsernamePassword() error {
	_, err := os.Stat(netRcPath)
	if err != nil {
		if os.IsNotExist(err) {
			content := fmt.Sprintf("default login %s password %s", provider.Username, provider.Password)
			err = ioutil.WriteFile(netRcPath, []byte(content), 0600)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
