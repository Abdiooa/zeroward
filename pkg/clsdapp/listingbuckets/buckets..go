package listingbuckets

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kataras/tablewriter"
)

const (
	awsYCEndpoint = "https://storage.yandexcloud.net"
	PatitionYCID  = "yc"
)

func ListBuckets(awsRegion string, accessKeyId string, accessKeySecret string) (err error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "auto" {
			return aws.Endpoint{
				PartitionID:   PatitionYCID,
				URL:           awsYCEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	s3Cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return fmt.Errorf("cannot load the AWS configs: %s", err)
	}

	client := s3.NewFromConfig(s3Cfg)

	result, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("listing buckets failed: %s", err)
	}

	// Create a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Bucket", "Creation Time"})

	for _, bucket := range result.Buckets {
		// Add a row for each bucket
		table.Append([]string{aws.ToString(bucket.Name), bucket.CreationDate.Format(time.RFC3339)})
	}

	// Optionally, customize the table
	table.SetBorder(true)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderColor(tablewriter.Color(tablewriter.BgBlackColor), tablewriter.Color(tablewriter.FgGreenColor))

	// Render the table
	table.Render()

	return nil
}
