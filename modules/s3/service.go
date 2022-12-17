package s3

import (
	"bytes"
	"encoding/json"
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
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
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
	files model.ProjectFiles,
) (model.ProjectFiles, error) {

	json, err := serializeProjectFiles(files)

	if err != nil {
		return nil, err
	}

	_, err = svc.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
		Key: aws.String(
			fmt.Sprintf("%s/%s/%s", userId, projectId, "project.json"),
		),
		Body: bytes.NewReader(json),
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

/*
GetProjectFiles gets a map of files contained in a project on s3. First it gets a list of all the files in the project, then it downloads each file concurrently.
*/
func (svc *S3Service) GetProjectFiles(userId, projectId string) (model.ProjectFiles, error) {

	json, err := svc.GetFile(fmt.Sprintf("%s/%s/%s", userId, projectId, "project.json"))

	if err != nil {
		return nil, err
	}

	files, err := deserializeProjectFiles(
		[]byte(json),
	)

	if err != nil {
		return nil, err
	}

	return files, nil

}

func (svc *S3Service) GetFile(path string) (string, error) {

	downloader := s3manager.NewDownloader(currentSession)

	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
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
			Bucket: aws.String(env.Get(env.S3_BUCKET)),
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
	errs := make(chan error)
	files := []string{"project.json", "main.wasm"}

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		go func(file string) {
			defer wg.Done()
			errs <- svc.DeleteFile(userId, projectId, file)
		}(file)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func serializeProjectFiles(files model.ProjectFiles) (json.RawMessage, error) {
	data, err := json.Marshal(files)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize project files: %w", err)
	}

	return data, nil
}

func deserializeProjectFiles(data []byte) (model.ProjectFiles, error) {
	var files model.ProjectFiles
	err := json.Unmarshal(data, &files)
	if err != nil {
		return model.ProjectFiles{}, err
	}

	return files, nil
}
