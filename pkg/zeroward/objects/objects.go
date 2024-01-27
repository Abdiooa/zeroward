package listingobjects

import (
	"context"
	"fmt"
	"log"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListObjects(awsRegion string, accessKeyId string, accessKeySecret string, bucketName string) (err error) {

	client, err := common.SetupS3Client(awsRegion, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	result, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Fatal(err)
	}

	header := []string{"Object", "Size (Bytes)", "Last Modified"}
	var rows [][]string

	for _, object := range result.Contents {
		// Add a row for each object
		rows = append(rows, []string{
			aws.ToString(object.Key),
			fmt.Sprintf("%d", object.Size),
			object.LastModified.Format("2006-01-02 15:04:05 Monday"),
		})
	}

	common.RenderTable(header, rows)

	return nil
}
