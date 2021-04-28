package storage

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/yellyoshua/whatsapp-chat-parser/constants"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
)

// Uploader __
type Uploader interface {
	UploadFiles(map[string]io.Reader) error
}

type uploader struct {
	s3Uploader *s3manager.Uploader
	s3Bucket   *string
}

// New __
func New() Uploader {
	s3BucketRegion := os.Getenv("AWS_REGION")
	currentSession, errSessionAWS := session.NewSession(&aws.Config{
		Region:      aws.String(s3BucketRegion),
		Credentials: credentials.NewEnvCredentials(),
	})

	if errSessionAWS != nil {
		logger.Fatal("error trying connect to storage -> " + errSessionAWS.Error())
	}

	s3Client := s3.New(currentSession)
	s3Uploader := s3manager.NewUploaderWithClient(s3Client)

	return &uploader{
		s3Uploader: s3Uploader,
		s3Bucket:   aws.String(constants.S3BucketName),
	}
}

func (u *uploader) UploadFiles(files map[string]io.Reader) error {
	var chError chan error = make(chan error)
	var wg sync.WaitGroup
	var err error

	uploadFileRoutine := func(fullPath string, f io.Reader, chError chan error, wg *sync.WaitGroup) {
		ctx := context.TODO()
		upParams := &s3manager.UploadInput{
			Bucket: u.s3Bucket,
			Key:    &fullPath,
			ACL:    aws.String("public-read"), // TODO: Set public files
			Body:   f,
		}
		defer wg.Done()

		_, err := u.s3Uploader.UploadWithContext(ctx, upParams)

		chError <- err
	}

	for fullPath, f := range files {
		wg.Add(1)
		go uploadFileRoutine(fullPath, f, chError, &wg)
	}

	go func() {
		wg.Wait()
		close(chError)
	}()

	for e := range chError {
		if e != nil {
			err = e
		}
	}

	return err
}
