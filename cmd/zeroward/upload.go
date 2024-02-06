package zeroward

import (
	"encoding/hex"
	"fmt"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/Abdiooa/zeroward/pkg/zeroward/encryption"
	"github.com/Abdiooa/zeroward/pkg/zeroward/genekeys"
	"github.com/Abdiooa/zeroward/pkg/zeroward/uploading"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const DEKKeyMetadataKey = "dek-key" // Metadata key for the DEK key

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload Command to upload Files on a cloud",
	Long:  `This Command is used to upload a secured(encrypted) user file on the cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {

		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")

		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")

		bcktName, _ := cmd.Flags().GetString("bcktname")

		passphrase, _ := cmd.Flags().GetString("passphrase")

		filePath, _ := cmd.Flags().GetString("filePath")

		objectkey, _ := cmd.Flags().GetString("objectkey")

		firstEncryption := common.IsFirstEncryption()

		if firstEncryption {
			if passphrase == "" {
				fmt.Println("Error: Passphrase is required for the first encryption. Please provide a passphrase using the --passphrase flag.")
				return
			}

			kekKey, err := genekeys.GenerateKek(passphrase)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			common.UpdateKEKKey(kekKey)
		}

		kekk := viper.GetString("KEKkey")

		kekBytes, err := hex.DecodeString(kekk)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		dek, err := genekeys.GenerateDek()
		cobra.CheckErr(err)

		if filePath != "" {
			if err := encryption.EncryptFile(filePath, dek); err != nil {
				fmt.Println("Error encrypting File:", err)
				return
			}
		}
		encryptedDek, err := encryption.EncryptData(dek, kekBytes)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		dekString := hex.EncodeToString(encryptedDek)

		// Create metadata map and add DEK key
		metadata := map[string]string{
			DEKKeyMetadataKey: dekString,
		}

		KeyAccessDefined := common.IsNotKeyAccessDefined()
		region := viper.GetString("Region")

		if KeyAccessDefined {
			if accessKeyID == "" || secretAccessKey == "" || bcktName == "" {
				fmt.Println("Error: Access Key ID and Secret Access Key are required as for your login and password of the Cloud Storage, also the bucket name is required!")
				return
			}
			if filePath == "" {
				fmt.Println("Error: FilePath are required!")
				return
			}
			common.StoreCredentials(accessKeyID, secretAccessKey, region)

			outputFilePath := fmt.Sprintf("%s.enc", filePath)
			err := uploading.UploadFile(region, accessKeyID, secretAccessKey, bcktName, outputFilePath, objectkey, metadata)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

		if bcktName == "" {
			fmt.Println("Error:  the bucket name is required!")
			return
		}

		if !KeyAccessDefined {
			if bcktName != "" {

				if filePath == "" {
					fmt.Println("Error: FilePath is required!")
					return
				} else {
					outputFilePath := fmt.Sprintf("%s.enc", filePath)
					err := uploading.UploadFile(region, viper.GetString("AWSAccessKeyID"), viper.GetString("AWSSecretAccessKey"), bcktName, outputFilePath, objectkey, metadata)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
