package downloading

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/Abdiooa/zeroward/pkg/zeroward/decryption"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/spf13/viper"
)

const DEKKeyMetadataKey = "dek-key" // Metadata key for the DEK key

func DownloadObject(awsRegion, accessKeyId, accessKeySecret, bucketName, localFilePath string, objectKey string) error {

	client, err := common.SetupS3Client(awsRegion, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	// Check if the object exists in the bucket
	_, err = client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				return fmt.Errorf("object not found: %s/%s", bucketName, objectKey)
			default:
				// Handle other errors
				return fmt.Errorf("error checking if the object exists: %v", err)
			}
		}
	}
	fileName := filepath.Base(objectKey)
	localFile := filepath.Join(localFilePath, fileName)
	decryptedFilePath := localFile[:len(localFile)-4]
	file, err := os.Create(decryptedFilePath)
	if err != nil {
		return fmt.Errorf("error creating local file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("error closing the local file: %v\n", closeErr)
		}
	}()

	// Download the object from S3
	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	if err != nil {
		return fmt.Errorf("error downloading the file from cloud storage: %v", err)
	}

	defer result.Body.Close()

	encryptedBody, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("error reading the body of the file: %v", err)
	}

	// Retrieve DEK key from metadata
	dekKeyString, ok := result.Metadata[DEKKeyMetadataKey]
	if !ok {
		return fmt.Errorf("dek key not found in metadata")
	}
	dekkeyEncrypted, err := hex.DecodeString(dekKeyString)
	if err != nil {
		return fmt.Errorf("error decoding DEK key: %v", err)
	}

	// Decrypt the DEK key with the KEK
	kekk := viper.GetString("KEKkey")
	kekBytes, err := hex.DecodeString(kekk)
	if err != nil {
		return fmt.Errorf("error decoding KEK key: %v", err)
	}

	dekkey, err := decryption.DecryptData(dekkeyEncrypted, kekBytes)
	if err != nil {
		return fmt.Errorf("error decrypting DEK key: %v", err)
	}

	// Decrypt the body with the DEK key
	body, err := decryption.DecryptData(encryptedBody, dekkey)
	if err != nil {
		return fmt.Errorf("error decrypting file body: %v", err)
	}

	_, err = file.Write(body)
	if err != nil {
		return fmt.Errorf("error writing the body on the local file: %v", err)
	}

	fmt.Printf("File '%s' downloaded successfully from S3://%s/%s to %s\n", objectKey, bucketName, objectKey, localFilePath)
	return nil

}
