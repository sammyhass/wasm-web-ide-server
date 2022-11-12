package s3

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sammyhass/web-ide/server/modules/env"
	"github.com/sammyhass/web-ide/server/modules/model"
)

var currentSession *session.Session

func InitSession() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
		Credentials: credentials.NewStaticCredentials(
			env.Env.S3_ACCESS_KEY_ID,
			env.Env.S3_SECRET_ACCESS_KEY,
			"",
		),
	})

	if err != nil {
		log.Fatal(err)
	}

	currentSession = sess
}

type S3Service struct {
	uploader *s3manager.Uploader
	s3       *s3.S3
}

func NewS3Service() *S3Service {
	return &S3Service{
		uploader: s3manager.NewUploader(currentSession),
		s3:       s3.New(currentSession),
	}
}

func (s *S3Service) UploadFile(
	userId, projectId, fileName, content string,
) (string, error) {

	res, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(env.Env.S3_BUCKET),
		Key: aws.String(
			fmt.Sprintf("%s/%s/%s", userId, projectId, fileName),
		),
		Body: strings.NewReader(content),
	})

	if err != nil {
		return "", err
	}

	return res.Location, nil
}

/*
UploadProjectFiles uploads a map of files to s3. It uses a wait group to wait for all the files to be uploaded.
*/
func (svc *S3Service) UploadProjectFiles(
	userId, projectId string,
	files map[string]string,
) (map[string]string, error) {
	wg := sync.WaitGroup{}

	wg.Add(len(files))

	res := make(map[string]string)

	errs := make(chan error, len(files))

	for fileName, content := range files {
		go func(fileName, content string) {
			defer wg.Done()

			location, err := svc.UploadFile(userId, projectId, fileName, content)
			if err != nil {
				errs <- err
			}

			res[fileName] = location
		}(fileName, content)
	}

	wg.Wait()

	close(errs)

	for err := range errs {
		return nil, err
	}

	return res, nil
}

/*
GetProjectFiles gets a map of files contained in a project on s3. First it gets a list of all the files in the project, then it downloads each file concurrently.
*/
func (svc *S3Service) GetProjectFiles(userId, projectId string) (model.ProjectFiles, error) {
	out, err := svc.s3.ListObjects(
		&s3.ListObjectsInput{
			Bucket: aws.String(env.Env.S3_BUCKET),
			Prefix: aws.String(fmt.Sprintf("%s/%s", userId, projectId)),
		},
	)

	if err != nil {
		return nil, err
	}

	files := make(model.ProjectFiles)

	wg := sync.WaitGroup{}
	errs := make(chan error, len(out.Contents))

	wg.Add(len(out.Contents))
	for _, obj := range out.Contents {
		go func(obj *s3.Object) {
			defer wg.Done()

			fileName := strings.Split(*obj.Key, "/")[2]

			content, err := svc.GetFile(*obj.Key)

			if err != nil {
				errs <- err
			}

			files[fileName] = content

		}(obj)
	}

	wg.Wait()

	close(errs)

	for err := range errs {
		return nil, err
	}

	return files, nil

}

func (svc *S3Service) GetFile(path string) (string, error) {

	downloader := s3manager.NewDownloader(currentSession)

	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(env.Env.S3_BUCKET),
		Key:    aws.String(path),
	})

	if err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}

func (svc *S3Service) DeleteFile(userId, projectId, fileName string) error {
	_, err := svc.s3.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(env.Env.S3_BUCKET),
			Key: aws.String(
				fmt.Sprintf("%s/%s/%s", userId, projectId, fileName),
			),
		},
	)

	return err
}

/*
*

	DeleteProjectFiles deletes all the files in a project in s3.
*/
func (svc *S3Service) DeleteProjectFiles(userId, projectId string) error {

	wg := sync.WaitGroup{}
	errs := make(chan error, len(model.DefaultFiles))
	wg.Add(len(model.DefaultFiles))

	for fname := range model.DefaultFiles {
		go func(fileName string) {
			defer wg.Done()
			err := svc.DeleteFile(userId, projectId, fileName)
			if err != nil {
				errs <- err
			}
		}(fname)
	}

	wg.Wait()

	close(errs)

	for err := range errs {
		return err
	}

	return nil
}
