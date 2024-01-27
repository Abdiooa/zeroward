package remove

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Remove(client *s3.Client, bucketName, objectKey string) error {
	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	if err != nil {
		return fmt.Errorf("error removing the file from cloud storage: %v", err)
	}

	return nil
}
