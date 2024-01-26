package uploading

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/common"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

const S3_ACL = "public-read"

func UploadFile(awsRegion, accessKeyId, accessKeySecret, bucketName, filePath string, objectKey string, metadata map[string]string) error {

	client, err := common.SetupS3Client(awsRegion, accessKeyId, accessKeySecret)
	if err != nil {
		return fmt.Errorf("failed to set up S3 client: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening the file: %s", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("error closing the file: %v\n", closeErr)
		}
	}()

	fileName := filepath.Base(filePath)
	objectKey = filepath.Join(objectKey, fileName) // Using specified path as object key

	// Check if the bucket exists
	_, err = client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				return fmt.Errorf("bucket not found: %s", bucketName)
			default:
				// Handle other errors
				return fmt.Errorf("error checking if the bucket exists: %v", err)
			}
		}
	}
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:   &bucketName,
		Key:      &objectKey,
		ACL:      S3_ACL,
		Body:     file,
		Metadata: metadata,
	})

	if err != nil {
		return fmt.Errorf("error uploading the file to cloud storage: %v", err)
	}

	fmt.Printf("File '%s' uploaded successfully to S3://%s/%s\n", fileName, bucketName, objectKey)

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error removing the file: %v", err)
	}

	return nil
}
