package s3

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sammyhass/web-ide/server/env"
	"github.com/sammyhass/web-ide/server/model"
)

var currentSession *session.Session

func InitSession() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
		Credentials: credentials.NewStaticCredentials(
			env.Get(env.S3_ACCESS_KEY_ID),
			env.Get(env.S3_SECRET_ACCESS_KEY),
			"",
		),
	})

	if err != nil {
		log.Fatal(err)
	}

	currentSession = sess
}

type Service struct {
	uploader *s3manager.Uploader
	s3       *s3.S3
}

func NewS3Service() *Service {
	return &Service{
		uploader: s3manager.NewUploader(currentSession),
		s3:       s3.New(currentSession),
	}
}

func (svc *Service) UploadFile(
	dir, fileName, content string,
) (string, error) {
	loc, err := svc.Upload(dir, fileName, strings.NewReader(content))

	if err != nil {
		return "", err
	}

	return loc, nil
}

func (svc *Service) Upload(
	dir, fileName string,
	r io.Reader,
) (string, error) {
	out, err := svc.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
		Key: aws.String(
			fmt.Sprintf("%s/%s", dir, fileName),
		),
		Body: r,
	})

	if err != nil {
		return "", err
	}

	return out.Location, nil
}

/*
GetFiles gets a map of files contained in a directory in the s3 bucket
*/
func (svc *Service) GetFiles(dir string) (map[string]string, error) {
	res, err := svc.s3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
		Prefix: aws.String(dir),
	})

	if err != nil {
		return nil, err
	}

	files := make(map[string]string)

	var wg sync.WaitGroup
	var mu sync.Mutex

	var errs chan error

	wg.Add(len(res.Contents))

	for _, obj := range res.Contents {
		go func(obj *s3.Object) {
			defer wg.Done()

			content, err := svc.GetFile(*obj.Key)

			if err != nil {
				errs <- err
				return
			}

			// get file name from key at the last index
			spl := strings.Split(*obj.Key, "/")
			fileName := spl[len(spl)-1]
			if fileName == "" {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			files[fileName] = content
		}(obj)
	}

	wg.Wait()

	for err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (svc *Service) Get(path string) (io.Reader, error) {
	downloader := s3manager.NewDownloader(currentSession)

	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
		Key:    aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(buf.Bytes())

	return r, nil
}

func (svc *Service) GetFile(path string) (string, error) {

	r, err := svc.Get(path)

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(r)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

/*
DeleteDir deletes all the files in a directory in s3.
*/
func (svc *Service) DeleteDir(dir string) error {
	res, err := svc.s3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
		Prefix: aws.String(dir),
	})

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var errs chan error

	wg.Add(len(res.Contents))

	for _, obj := range res.Contents {
		go func(obj *s3.Object) {
			defer wg.Done()
			if _, err := svc.s3.DeleteObject(
				&s3.DeleteObjectInput{
					Bucket: aws.String(env.Get(env.S3_BUCKET)),
					Key:    obj.Key,
				},
			); err != nil {
				errs <- err
			}
		}(obj)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) UploadFiles(dir string, files model.ProjectFiles) error {
	var wg sync.WaitGroup

	errs := make(chan error)

	wg.Add(len(files))

	for name, content := range files {
		go func(name, content string) {
			defer wg.Done()
			_, err := svc.UploadFile(dir, name, content)
			if err != nil {
				errs <- err
			}
		}(name, content)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) GenPresignedURL(path string, exp time.Duration) (string, error) {
	contentType := mime.TypeByExtension(filepath.Ext(path))
	req, _ := svc.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket:              aws.String(env.Get(env.S3_BUCKET)),
		Key:                 aws.String(path),
		ResponseContentType: aws.String(contentType),
	})

	url, err := req.Presign(exp)

	if err != nil {
		return "", err
	}

	return url, nil
}
