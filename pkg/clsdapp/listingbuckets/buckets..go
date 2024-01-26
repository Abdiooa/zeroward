package listingbuckets

import (
	"context"
	"fmt"
	"time"

	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ListBuckets(awsRegion string, accessKeyId string, accessKeySecret string) (err error) {

	client, err := common.SetupS3Client(awsRegion, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("listing buckets failed: %s", err)
	}

	header := []string{"Bucket", "Creation Time"}
	var rows [][]string

	for _, bucket := range result.Buckets {
		// Add a row for each bucket
		rows = append(rows, []string{
			aws.ToString(bucket.Name),
			bucket.CreationDate.Format(time.RFC3339),
		})
	}

	common.RenderTable(header, rows)

	return nil
}
