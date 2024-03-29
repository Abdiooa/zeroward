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

const DEKKeyMetadataKey = "dek-key"

func DownloadObject(awsRegion, accessKeyId, accessKeySecret, bucketName, localFilePath string, objectKey string, removeAfterDownload bool) error {

	client, err := common.SetupS3Client(awsRegion, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

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

	dekKeyString, ok := result.Metadata[DEKKeyMetadataKey]
	if !ok {
		return fmt.Errorf("dek key not found in metadata")
	}
	dekkeyEncrypted, err := hex.DecodeString(dekKeyString)
	if err != nil {
		return fmt.Errorf("error decoding DEK key: %v", err)
	}

	kekk := viper.GetString("KEKkey")
	kekBytes, err := hex.DecodeString(kekk)
	if err != nil {
		return fmt.Errorf("error decoding KEK key: %v", err)
	}

	dekkey, err := decryption.DecryptData(dekkeyEncrypted, kekBytes)
	if err != nil {
		return fmt.Errorf("error decrypting DEK key: %v", err)
	}
	keySize := 60
	encryptedBody = encryptedBody[:len(encryptedBody)-keySize]
	body, err := decryption.DecryptFile(encryptedBody, dekkey)
	if err != nil {
		return fmt.Errorf("error decrypting file body: %v", err)
	}

	_, err = file.Write(body)
	if err != nil {
		return fmt.Errorf("error writing the body on the local file: %v", err)
	}
	outputFile := objectKey[:len(objectKey)-4]
	fmt.Printf("File '%s' downloaded successfully from S3://%s/%s to %s\n", outputFile, bucketName, objectKey, localFilePath)

	if removeAfterDownload {
		err := removeFileFromCloudStorage(client, bucketName, objectKey)
		if err != nil {
			return fmt.Errorf("error removing file from cloud storage: %v", err)
		}
		fmt.Printf("File '%s' removed from S3://%s/%s\n", objectKey, bucketName, objectKey)
	}

	return nil

}
func DownloadNormalObject(awsRegion, accessKeyId, accessKeySecret, bucketName, localFilePath string, objectKey string, removeAfterDownload bool) error {
	client, err := common.SetupS3Client(awsRegion, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

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

				return fmt.Errorf("error checking if the object exists: %v", err)
			}
		}
	}
	fileName := filepath.Base(objectKey)
	localFile := filepath.Join(localFilePath, fileName)

	file, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("error creating local file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("error closing the local file: %v\n", closeErr)
		}
	}()

	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	if err != nil {
		return fmt.Errorf("error downloading the file from cloud storage: %v", err)
	}

	defer result.Body.Close()

	Body, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("error reading the body of the file: %v", err)
	}
	_, err = file.Write(Body)
	if err != nil {
		return fmt.Errorf("error writing the body on the local file: %v", err)
	}

	fmt.Printf("File '%s' downloaded successfully from S3://%s/%s to %s\n", objectKey, bucketName, objectKey, localFilePath)

	if removeAfterDownload {
		err := removeFileFromCloudStorage(client, bucketName, objectKey)
		if err != nil {
			return fmt.Errorf("error removing file from cloud storage: %v", err)
		}
		fmt.Printf("File '%s' removed from S3://%s/%s\n", objectKey, bucketName, objectKey)
	}

	return nil
}
func removeFileFromCloudStorage(client *s3.Client, bucketName, objectKey string) error {
	_, err := client.DeleteObject(
		context.TODO(),
		&s3.DeleteObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
		},
	)
	if err != nil {
		return fmt.Errorf("error deleting file from cloud storage: %v", err)
	}

	return nil
}
