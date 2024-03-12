package common

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

func IsNotKeyAccessDefined() bool {
	AWSAccessKeyID := viper.GetString("AWSAccessKeyID")
	AWSSecretAccessKey := viper.GetString("AWSSecretAccessKey")

	return AWSAccessKeyID == "" && AWSSecretAccessKey == ""
}
func StoreCredentials(accessKeyID, secretAccessKey, region string) {
	viper.Set("AWSAccessKeyID", accessKeyID)
	viper.Set("AWSSecretAccessKey", secretAccessKey)
	viper.Set("Region", region)

	var awsDir string
	var credentialsFilePath string
	var awsConfigFilePath string

	switch runtime.GOOS {
	case "windows":
		homeDir := os.Getenv("USERPROFILE")
		awsDir = filepath.Join(homeDir, ".aws")
	case "linux", "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			os.Exit(1)
		}
		awsDir = filepath.Join(homeDir, ".aws")
	default:
		fmt.Println("Unsupported operating system")
		os.Exit(1)
	}

	if _, err := os.Stat(awsDir); os.IsNotExist(err) {
		err := os.Mkdir(awsDir, 0700)
		if err != nil {
			fmt.Println("Error creating .aws directory:", err)
			os.Exit(1)
		}
	}

	credentialsFilePath = filepath.Join(awsDir, "credentials")
	credentialsFile, err := os.Create(credentialsFilePath)
	if err != nil {
		fmt.Println("Error creating credentials file:", err)
		os.Exit(1)
	}
	defer credentialsFile.Close()

	credentialsFile.WriteString(fmt.Sprintf("aws_access_key_id = %s\n", accessKeyID))
	credentialsFile.WriteString(fmt.Sprintf("aws_secret_access_key = %s\n", secretAccessKey))

	awsConfigFilePath = filepath.Join(awsDir, "config")
	awsconfigFile, err := os.Create(awsConfigFilePath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		os.Exit(1)
	}
	defer awsconfigFile.Close()

	awsconfigFile.WriteString(fmt.Sprintf("region = %s\n", region))

	if err := viper.WriteConfig(); err != nil {
		fmt.Println("Error updating the access key id, secret key, and region in the config file:", err)
		os.Exit(1)
	}
}

func IsFirstEncryption() bool {
	kekKey := viper.GetString("KEKkey")
	return kekKey == ""
}

func UpdateKEKKey(kek []byte) {

	kekString := hex.EncodeToString(kek)
	viper.Set("KEKkey", kekString)
	if err := viper.WriteConfig(); err != nil {
		fmt.Println("Error updating KEK key in the config file:", err)
		os.Exit(1)
	}
}

func HandleCredentials(accessKeyID, secretAccessKey string) (string, string, string, error) {
	KeyAccessDefined := IsNotKeyAccessDefined()
	region := viper.GetString("Region")

	if KeyAccessDefined {
		if accessKeyID == "" || secretAccessKey == "" {
			return "", "", "", fmt.Errorf("error: Access Key ID and Secret Access Key are required as for your login and password of the Cloud Storage, also the bucket name is required")
		}

		StoreCredentials(accessKeyID, secretAccessKey, region)
		return accessKeyID, secretAccessKey, region, nil
	}

	return viper.GetString("AWSAccessKeyID"), viper.GetString("AWSSecretAccessKey"), region, nil
}
