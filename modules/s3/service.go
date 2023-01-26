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
	userId, projectId, fileName, content string,
) (string, error) {

	res, err := svc.uploader.Upload(&s3manager.UploadInput{
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
func (svc *Service) UploadProjectFiles(
	userId, projectId string,
	files model.ProjectFiles,
) (model.ProjectFiles, error) {
	err := svc.UploadFiles(userId, projectId, files)

	if err != nil {
		return nil, err
	}

	return files, nil
}

/*
GetProjectFiles gets a map of files contained in a project on s3. First it gets a list of all the files in the project, then it downloads each file concurrently.
*/
func (svc *Service) GetProjectFiles(userId, projectId string) (model.ProjectFiles, error) {
	res, err := svc.s3.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(env.Get(env.S3_BUCKET)),
		Prefix: aws.String(fmt.Sprintf("%s/%s", userId, projectId)),
	})

	if err != nil {
		return nil, err
	}

	files := make(model.ProjectFiles)

	for _, obj := range res.Contents {
		key := strings.Split(*obj.Key, "/")[2]
		fmt.Printf("%v / %s", key, *obj.Key)
		file, err := svc.GetFile(*obj.Key)
		if err != nil {
			return nil, err
		}
		files[key] = file
	}

	return files, nil
}

func (svc *Service) GetFile(path string) (string, error) {

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

func (svc *Service) DeleteFile(userId, projectId, fileName string) error {
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
DeleteProjectFiles deletes all the files in a project in s3.
*/
func (svc *Service) DeleteProjectFiles(userId, projectId string) error {
	errs := make(chan error)
	files := model.DefaultFiles

	var wg sync.WaitGroup
	wg.Add(len(files))

	for name, file := range files {
		go func(file string) {
			defer wg.Done()
			errs <- svc.DeleteFile(userId, projectId, name)
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

func (svc *Service) UploadFiles(userId string, projectId string, files model.ProjectFiles) error {
	var wg sync.WaitGroup

	errs := make(chan error)

	wg.Add(len(files))

	for name, content := range files {
		go func(name, content string) {
			defer wg.Done()
			_, err := svc.UploadFile(userId, projectId, name, content)
			if err != nil {
				errs <- err
			}
		}(name, content)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
