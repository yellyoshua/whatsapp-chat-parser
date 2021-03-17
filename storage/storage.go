package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/iam"
	stg "cloud.google.com/go/storage"
	"github.com/yellyoshua/whatsapp-chat-parser/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

var credentialsFile string = "keys.json"

// Uploader __
type Uploader interface {
	UploadFiles(map[string]io.Reader) error
}

var (
	// RoleObjectViewer __
	RoleObjectViewer = "roles/storage.objectViewer"
)

type uploader struct {
	bucket         *stg.BucketHandle
	role           string
	serviceAccount string
	bucketName     string
}

// New __
func New() Uploader {
	var bucket string = os.Getenv("GCS_BUCKET")
	var projectID string = os.Getenv("GCS_PROJECT_ID")
	var bucketAttrs *stg.BucketAttrs = &stg.BucketAttrs{
		Location:     "US",
		LocationType: "multi-region",
		StorageClass: "STANDARD",
	}
	var srvAccount = os.Getenv("GCS_IAM_SVC")

	ctx, close := context.WithTimeout(context.TODO(), 5*time.Second)
	defer close()

	storageClient, err := stg.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		logger.Fatal("Error setup client storage -> %s", err)
	}

	var up *uploader = &uploader{
		bucket:         storageClient.Bucket(bucket),
		bucketName:     bucket,
		role:           RoleObjectViewer,
		serviceAccount: srvAccount,
	}

	if err := createGCSBucket(up, projectID, bucketAttrs, storageClient); err != nil {
		logger.Fatal("Error creating bucket -> %s", err)
	}

	return up
}

func (u *uploader) UploadFiles(files map[string]io.Reader) error {
	var chans chan error = make(chan error)
	var wg sync.WaitGroup
	var err error

	uploadFileRoutine := func(ctx context.Context, fullPath string, f io.Reader, c chan error, wg *sync.WaitGroup) {
		obj := u.bucket.Object(fullPath)
		w := obj.NewWriter(context.TODO())

		if _, err := io.Copy(w, f); err != nil {
			c <- err
			return
		}

		defer wg.Done()

		if err := w.Close(); err != nil {
			c <- err
			return
		}

		// Make public the object to internet
		if err := setObjectPublic(ctx, obj); err != nil {
			c <- err
			return
		}

		c <- nil
		return
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

func createGCSBucket(up *uploader, projectID string, bucketAttrs *stg.BucketAttrs, client *stg.Client) error {
	// Setup context and client
	ctx := context.TODO()

	var bucket *stg.BucketHandle = up.bucket
	var bucketName string = up.bucketName

	buckets := client.Buckets(ctx, projectID)
	for {
		attrs, err := buckets.Next()

		// Assume bucket not found if at Iterator end and create
		if err == iterator.Done {
			// Create bucket
			if err := bucket.Create(ctx, projectID, bucketAttrs); err != nil {
				return fmt.Errorf("Failed to create bucket: %v", err)
			}

			// // Make public the bucket to internet
			// if err := setBucketPublic(ctx, up); err != nil {
			// 	return fmt.Errorf("Failed set public bucket: %v", err)
			// }

			log.Printf("Bucket %v created and public.\n", bucketName)
			return nil
		}

		if err != nil {
			return fmt.Errorf("Issues setting up Bucket(%q).Objects(): %v. Double check project id", attrs.Name, err)
		}
		if attrs.Name == bucketName {
			log.Printf("Bucket %v exists.\n", bucketName)
			return nil
		}
	}
}

func setObjectPublic(ctx context.Context, obj *stg.ObjectHandle) error {
	err := obj.ACL().Set(ctx, stg.AllUsers, stg.RoleReader)
	return err
}

func setBucketPublic(ctx context.Context, up *uploader) error {
	var roleObjectViewer string = up.role

	policy, err := up.bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}

	// Binding new policy that make public the bucket to internet
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    roleObjectViewer,
		Members: []string{iam.AllUsers},
	})

	if err := up.bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return err
	}
	return nil
}
