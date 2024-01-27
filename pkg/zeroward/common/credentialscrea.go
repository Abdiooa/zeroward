package common

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	// Create the .aws directory if it doesn't exist
	awsDir := filepath.Join(homeDir, ".aws")
	if _, err := os.Stat(awsDir); os.IsNotExist(err) {
		err := os.Mkdir(awsDir, 0700)
		if err != nil {
			fmt.Println("Error creating .aws directory:", err)
			os.Exit(1)
		}
	}

	// Write credentials to credentials file
	credentialsFilePath := filepath.Join(awsDir, "credentials")
	credentialsFile, err := os.Create(credentialsFilePath)
	if err != nil {
		fmt.Println("Error creating credentials file:", err)
		os.Exit(1)
	}
	defer credentialsFile.Close()

	// credentialsFile.WriteString(fmt.Sprintf("[default]\n"))
	credentialsFile.WriteString(fmt.Sprintf("aws_access_key_id = %s\n", accessKeyID))
	credentialsFile.WriteString(fmt.Sprintf("aws_secret_access_key = %s\n", secretAccessKey))

	// Write region to config file
	awsConfigFilePath := filepath.Join(awsDir, "config")
	awsconfigFile, err := os.Create(awsConfigFilePath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		os.Exit(1)
	}
	defer awsconfigFile.Close()

	// awsconfigFile.WriteString(fmt.Sprintf("[default]\n"))
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