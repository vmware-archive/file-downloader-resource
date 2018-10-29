package file

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/pivotal-cf/go-pivnet/download"
	"github.com/pivotal-cf/go-pivnet/logger"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/shirou/gopsutil/disk"
)

//go:generate counterfeiter -o ./fakes/bar.go --fake-name Bar . bar
type bar interface {
	SetTotal(contentLength int64)
	SetOutput(output io.Writer)
	Add(totalWritten int) int
	Kickoff()
	Finish()
	NewProxyReader(reader io.Reader) io.Reader
}

type HTTPProvider struct {
	BaseURL        string
	HTTPClient     *http.Client
	Bar            bar
	ProgressWriter io.Writer
	Logger         logger.Logger
}

const URL_PATTERN = "%s/%s/%s/%s"

func NewHTTPProvider(skipSSLValidation bool, baseURL string) (Provider, error) {
	downloadClient := &http.Client{
		Timeout: 0,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipSSLValidation,
			},
			Proxy: http.ProxyFromEnvironment,
		},
	}
	logWriter := os.Stderr
	logger := log.New(logWriter, "", log.LstdFlags)
	return &HTTPProvider{
		HTTPClient:     downloadClient,
		Logger:         logshim.NewLogShim(logger, logger, false),
		ProgressWriter: logWriter,
		BaseURL:        baseURL,
	}, nil

}

func (h *HTTPProvider) DownloadFile(targetDirectory, productSlug, version, pattern string, unpack bool) error {

	if err := os.MkdirAll(targetDirectory, os.ModePerm); err != nil {
		return err
	}
	fileName := h.FileName(version, pattern)
	contentURL := h.ContentURL(productSlug, version, pattern)
	targetFile := path.Join(targetDirectory, fileName)
	return h.Download(targetFile, contentURL)
}

func (h *HTTPProvider) FileName(version, pattern string) string {
	return strings.Replace(pattern, "-*", fmt.Sprintf("-%s", version), 1)
}
func (h *HTTPProvider) ContentURL(slug, version, pattern string) string {
	return fmt.Sprintf(URL_PATTERN, h.BaseURL, slug, version, h.FileName(version, pattern))
}

func (h *HTTPProvider) Download(
	targetFile string,
	contentURL string,
) error {

	h.Bar = download.NewBar()
	resp, err := h.HTTPClient.Head(contentURL)
	if err != nil {
		return fmt.Errorf("failed to make HEAD request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status for url %s: %d", contentURL, resp.StatusCode)
	}
	contentURL = resp.Request.URL.String()
	diskStats, err := disk.Usage(path.Dir(targetFile))
	if err != nil {
		return fmt.Errorf("failed to get disk free space: %s", err)
	}
	if diskStats.Free < uint64(resp.ContentLength) {
		return fmt.Errorf("file is too big to fit on this drive")
	}

	h.Bar.SetOutput(h.ProgressWriter)
	h.Bar.SetTotal(resp.ContentLength)
	h.Bar.Kickoff()

	defer h.Bar.Finish()
	err = h.retryableRequest(contentURL, targetFile)
	if err != nil {
		return fmt.Errorf("failed during retryable request: %s", err)
	}
	return nil
}

func (h *HTTPProvider) retryableRequest(contentURL string, targetFilePath string) error {
	targetFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	fileInfo, err := targetFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to read information from output file: %s", err)
	}
	fileWriter, err := os.OpenFile(targetFile.Name(), os.O_RDWR, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %s", err)
	}
	defer fileWriter.Close()
Retry:
	resp, err := h.HTTPClient.Get(contentURL)
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			if netErr.Temporary() {
				goto Retry
			}
		}

		return fmt.Errorf("download request failed: %s", err)
	}

	defer resp.Body.Close()

	var proxyReader io.Reader
	proxyReader = h.Bar.NewProxyReader(resp.Body)

	bytesWritten, err := io.Copy(fileWriter, proxyReader)
	if err != nil {
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			h.Logger.Info(fmt.Sprintf("retrying %v", err))
			h.Bar.Add(int(-1 * bytesWritten))
			goto Retry
		}
		oe, _ := err.(*net.OpError)
		if strings.Contains(oe.Err.Error(), syscall.ECONNRESET.Error()) {
			h.Bar.Add(int(-1 * bytesWritten))
			goto Retry
		}
		return fmt.Errorf("failed to write file during io.Copy: %s", err)
	}

	return nil
}
