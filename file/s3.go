package file

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type S3Provider struct {
	Client         s3iface.S3API
	ProgressOutput io.Writer
	BucketName     string
}

func NewS3Provider(accessKeyID, secretAccessKey, regionName, endpoint, bucketName string, skipSSLVerification, disableSSL, useV2Signing bool) (Provider, error) {
	var creds *credentials.Credentials

	if accessKeyID == "" && secretAccessKey == "" {
		creds = credentials.AnonymousCredentials
	} else {
		creds = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	}

	if len(regionName) == 0 {
		regionName = "us-east-1"
	}

	var httpClient *http.Client
	if skipSSLVerification {
		httpClient = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	} else {
		httpClient = http.DefaultClient
	}

	awsConfig := &aws.Config{
		Region:           aws.String(regionName),
		Credentials:      creds,
		S3ForcePathStyle: aws.Bool(true),
		MaxRetries:       aws.Int(maxRetries),
		DisableSSL:       aws.Bool(disableSSL),
		HTTPClient:       httpClient,
	}

	if len(endpoint) != 0 {
		awsConfig.Endpoint = aws.String(endpoint)
	}

	client := s3.New(session.New(awsConfig), awsConfig)

	if useV2Signing {
		setv2Handlers(client)
	}

	return &S3Provider{
		Client:         client,
		BucketName:     bucketName,
		ProgressOutput: os.Stderr,
	}, nil
}

//DownloadFile - Downloads file based on version info
func (p *S3Provider) DownloadFile(targetDirectory, productSlug, version, pattern string) error {

	var (
		localPath     string
		remotePath    string
		contentLength int64
	)
	if err := os.MkdirAll(targetDirectory, os.ModePerm); err != nil {
		return err
	}

	bucketFiles, err := p.Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(p.BucketName),
		Prefix: aws.String(productSlug),
	})

	if err != nil {
		return err
	}

	for _, bucketFile := range bucketFiles.Contents {
		matched, err := filepath.Match(path.Join(productSlug, pattern), *bucketFile.Key)
		if err != nil {
			return err
		}
		if matched {
			fileName := strings.Replace(*bucketFile.Key, productSlug+"/", "", 1)
			localPath = path.Join(targetDirectory, fileName)
			remotePath = *bucketFile.Key
			contentLength = *bucketFile.Size
		}
	}
	if localPath != "" {
		progress := p.newProgressBar(contentLength)

		downloader := s3manager.NewDownloaderWithClient(p.Client)
		localFile, err := os.Create(localPath)
		if err != nil {
			return err
		}
		defer localFile.Close()

		getObject := &s3.GetObjectInput{
			Bucket: aws.String(p.BucketName),
			Key:    aws.String(remotePath),
		}

		progress.Start()
		defer progress.Finish()

		_, err = downloader.Download(progressWriterAt{localFile, progress}, getObject)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("No files found in bucket %s, folder %s matching pattern %s", p.BucketName, productSlug, pattern)
	}

	return nil
}

type progressWriterAt struct {
	io.WriterAt
	*pb.ProgressBar
}

func (pwa progressWriterAt) WriteAt(p []byte, off int64) (int, error) {
	n, err := pwa.WriterAt.WriteAt(p, off)
	if err != nil {
		return n, err
	}

	pwa.ProgressBar.Add(len(p))

	return n, err
}

func (p *S3Provider) newProgressBar(total int64) *pb.ProgressBar {
	progress := pb.New64(total)

	progress.Output = p.ProgressOutput
	progress.ShowSpeed = true
	progress.Units = pb.U_BYTES
	progress.NotPrint = true

	return progress.SetWidth(80)
}
