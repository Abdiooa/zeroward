package common

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
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

func SetupS3Client(awsRegion, accessKeyId, accessKeySecret string) (*s3.Client, error) {
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

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	s3Cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion(awsRegion),
	)

	s3Cfg.HTTPClient = httpClient

	if err != nil {
		return nil, fmt.Errorf("cannot load the AWS configs: %s", err)
	}

	client := s3.NewFromConfig(s3Cfg)
	return client, nil
}

func RenderTable(header []string, rows [][]string) {
	// Create a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, row := range rows {
		// Add a row for each item
		table.Append(row)
	}

	// Optionally, customize the table
	table.SetBorder(true)
	table.SetCenterSeparator("│")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	if len(header) == 2 {
		table.SetHeaderColor(tablewriter.Color(tablewriter.BgBlackColor), tablewriter.Color(tablewriter.FgGreenColor))
	} else {
		table.SetHeaderColor(tablewriter.Color(tablewriter.BgBlackColor), tablewriter.Color(tablewriter.FgGreenColor), tablewriter.Color(tablewriter.FgGreenColor))
	}

	// Render the table
	table.Render()
}
