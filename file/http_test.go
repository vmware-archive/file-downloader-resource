package file_test

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/pivotalservices/file-downloader-resource/file"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestHTTPProviderProvider(t *testing.T) {
	spec.Run(t, "HttpProvider", testHTTPProvider, spec.Report(report.Terminal{}))
}

func testHTTPProvider(t *testing.T, when spec.G, it spec.S) {
	var provider *file.HTTPProvider
	var targetFile string
	var httpClient *http.Client
	var server *ghttp.Server
	//targetURL := "http://example.com/products/elastic-runtime/2.3.0/cf-2.3.0.pivotal"
	it.Before(func() {
		RegisterTestingT(t)
		server = ghttp.NewServer()
		targetFile = filepath.Join(os.TempDir(), "temp.file")
		httpClient = http.DefaultClient
		logWriter := os.Stderr
		logger := log.New(logWriter, "", log.LstdFlags)
		provider = &file.HTTPProvider{
			BaseURL:        server.URL(),
			HTTPClient:     httpClient,
			Logger:         logshim.NewLogShim(logger, logger, false),
			ProgressWriter: logWriter,
		}
	})
	it.After(func() {
		server.Close()
		os.RemoveAll(targetFile)
	})
	when("Getting content url", func() {
		it("returns a url", func() {
			targetURL := fmt.Sprintf("%s/elastic-runtime/2.3.0/cf-2.3.0.pivotal", server.URL())
			Expect(provider.ContentURL("elastic-runtime", "2.3.0", "cf-*.pivotal")).Should(BeEquivalentTo(targetURL))
		})
	})

	when("Getting file name", func() {
		it("returns a filename", func() {
			Expect(provider.FileName("2.3.0", "cf-*.pivotal")).Should(Equal("cf-2.3.0.pivotal"))
		})
	})

	when("Download", func() {
		it("successfully downloads", func() {
			bytes := make([]byte, 40000)
			rand.Read(bytes)
			length := strconv.Itoa(len(bytes))
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("HEAD", "/elastic-runtime/2.3.0/cf-2.3.0.pivotal"),
					func(w http.ResponseWriter, r *http.Request) {
						w.Header().Add("Content-Length", length)
					},
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/elastic-runtime/2.3.0/cf-2.3.0.pivotal"),
					func(w http.ResponseWriter, r *http.Request) {
						w.Write(bytes)
						w.Header().Add("Content-Length", length)
					},
				),
			)
			targetURL := fmt.Sprintf("%s/elastic-runtime/2.3.0/cf-2.3.0.pivotal", server.URL())
			err := provider.Download(targetFile, targetURL)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

}
