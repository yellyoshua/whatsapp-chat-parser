package storage

import (
	"context"
	"io"
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

	mySession, errSessionAWS := session.NewSession(&aws.Config{
		Region:      aws.String(constants.S3BucketRegion),
		Credentials: credentials.NewEnvCredentials(),
	})
	logger.CheckError("error creating session s3", errSessionAWS)

	s3Client := s3.New(mySession)
	s3Uploader := s3manager.NewUploaderWithClient(s3Client)

	return &uploader{
		s3Uploader: s3Uploader,
		s3Bucket:   aws.String(constants.S3BucketName),
	}
}

func (u *uploader) UploadFiles(files map[string]io.Reader) error {
	var chans chan error = make(chan error)
	var wg sync.WaitGroup
	var err error

	uploadFileRoutine := func(ctx context.Context, fullPath string, f io.Reader, c chan error, wg *sync.WaitGroup) {
		upParams := &s3manager.UploadInput{
			Bucket: u.s3Bucket,
			Key:    &fullPath,
			ACL:    aws.String("public-read"), // TODO: Set public files
			Body:   f,
		}

		defer wg.Done()

		_, err := u.s3Uploader.UploadWithContext(ctx, upParams)
		c <- err
	}

	for fullPath, f := range files {
		wg.Add(1)
		ctx := context.Background()
		go uploadFileRoutine(ctx, fullPath, f, chans, &wg)
	}

	go func() {
		for c := range chans {
			if c != nil {
				err = c
			}
		}
	}()

	wg.Wait()
	close(chans)

	return err
}
