package config_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pivotalservices/file-downloader-resource/config"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestGitProvider(t *testing.T) {
	spec.Run(t, "GitProvider", testGitProvider, spec.Report(report.Terminal{}))
}

func testGitProvider(t *testing.T, when spec.G, it spec.S) {
	var (
		gitRepoDir, privateKeyPath, netRcPath, tempRepo string
		provider                                        *config.GitProvider
	)
	it.Before(func() {
		RegisterTestingT(t)
		gitRepoDir = filepath.Join(os.TempDir(), "file-downloader-git-repo")
		tempRepo = filepath.Join(os.TempDir(), "test-repo")
		privateKeyPath = filepath.Join(os.TempDir(), "private-key")
		netRcPath = filepath.Join(os.Getenv("HOME"), ".netrc")
		gitClone := exec.Command("git", "clone", "fixtures/test.bundle", "--branch", "master")
		gitClone.Args = append(gitClone.Args, tempRepo)
		_, err := execCommand(".", gitClone)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = execCommand(tempRepo, exec.Command("git", "config", "--add", "receive.denyCurrentBranch", "ignore"))
		Expect(err).ShouldNot(HaveOccurred())
		provider = &config.GitProvider{
			URI:    tempRepo,
			Branch: "master",
		}
	})
	it.After(func() {
		os.RemoveAll(gitRepoDir)
		os.RemoveAll(tempRepo)
		os.RemoveAll(privateKeyPath)
		os.RemoveAll(netRcPath)
	})
	when("Repo hasn't been cloned", func() {
		it.Before(func() {
			os.RemoveAll(gitRepoDir)
		})
		it("returns a version", func() {
			version, err := provider.LatestVersion()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(version).ShouldNot(BeNil())
			Expect(version.Ref).Should(Equal("a2b5630e85d4a72280fd825da5fddad7398aa8e3"))
		})
	})

	when("new version after a commit", func() {
		it.Before(func() {
			os.RemoveAll(gitRepoDir)
		})
		it("returns a version", func() {
			originalVersion, err := provider.LatestVersion()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(originalVersion).ShouldNot(BeNil())
			Expect(originalVersion.Ref).Should(Equal("a2b5630e85d4a72280fd825da5fddad7398aa8e3"))

			versionCreated, err := createCommit(gitRepoDir)
			Expect(err).ShouldNot(HaveOccurred())

			newVersion, err := provider.LatestVersion()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(newVersion).ShouldNot(BeNil())
			Expect(newVersion.Ref).Should(Equal(strings.Replace(versionCreated, "\n", "", -1)))
		})
	})

	when("Repo has been cloned", func() {
		it.Before(func() {
			os.RemoveAll(gitRepoDir)
		})
		it("returns a version", func() {
			_, err := provider.LatestVersion()
			Expect(err).ShouldNot(HaveOccurred())
			version, err := provider.LatestVersion()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(version).ShouldNot(BeNil())
			Expect(version.Ref).Should(Equal("a2b5630e85d4a72280fd825da5fddad7398aa8e3"))
		})
	})
}

func createCommit(gitRepoDir string) (string, error) {
	_, err := execCommand(gitRepoDir, exec.Command("touch", "foo.txt"))
	if err != nil {
		return "", err
	}

	_, err = execCommand(gitRepoDir, exec.Command("git", "add", "--all"))
	if err != nil {
		return "", err
	}
	_, err = execCommand(gitRepoDir, exec.Command("git", "commit", "-m", "Testing"))
	if err != nil {
		return "", err
	}
	output, err := execCommand(gitRepoDir, exec.Command("git", "push"))
	fmt.Println(output)
	if err != nil {
		return "", err
	}

	return execCommand(gitRepoDir, exec.Command("git", "rev-parse", "HEAD"))
}

func execCommand(directory string, command *exec.Cmd) (string, error) {
	command.Dir = directory
	command.Stderr = os.Stderr
	output, err := command.Output()
	return string(output), err
}
