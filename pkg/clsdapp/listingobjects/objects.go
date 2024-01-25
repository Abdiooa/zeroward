package listingobjects

import (
	"context"
	"fmt"
	"log"
	"os"

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

func ListObjects(awsRegion string, accessKeyId string, accessKeySecret string, bucketName string) (err error) {
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

	result, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Fatal(err)
	}

	// Create a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Object", "Size (Bytes)", "Last Modified"})

	for _, object := range result.Contents {
		// Add a row for each object
		table.Append([]string{
			aws.ToString(object.Key),
			fmt.Sprintf("%d", object.Size),
			object.LastModified.Format("2006-01-02 15:04:05 Monday"),
		})
	}

	// Optionally, customize the table
	table.SetBorder(true)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderColor(tablewriter.Color(tablewriter.BgBlackColor), tablewriter.Color(tablewriter.FgGreenColor), tablewriter.Color(tablewriter.FgGreenColor))

	// Render the table
	table.Render()

	return nil
}
