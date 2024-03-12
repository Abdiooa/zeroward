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

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload Command to upload Files on a cloud",
	Long:  `This Command is used to upload a secured(encrypted) user file on the cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")
		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")
		bcktName, _ := cmd.Flags().GetString("bucketName")
		passphrase, _ := cmd.Flags().GetString("passphrase")
		filePath, _ := cmd.Flags().GetString("filePath")
		objectkey, _ := cmd.Flags().GetString("objectkey")
		encrypt, _ := cmd.Flags().GetBool("encrypt")

		if encrypt {
			if err := handleEncryptionAndUpload(accessKeyID, secretAccessKey, bcktName, filePath, passphrase, objectkey); err != nil {
				fmt.Println("Error:", err)
				return
			}
		} else {
			if err := handleUpload(accessKeyID, secretAccessKey, bcktName, filePath, objectkey, nil); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
	},
}

func handleEncryptionAndUpload(accessKeyID, secretAccessKey, bcktName, filePath, passphrase, objectkey string) error {
	if common.IsFirstEncryption() {
		if passphrase == "" {
			return fmt.Errorf("error: Passphrase is required for the first encryption. Please provide a passphrase using the --passphrase flag")
		}
		kekKey, err := genekeys.GenerateKek(passphrase)
		if err != nil {
			return fmt.Errorf("error generating KEK key: %v", err)
		}
		common.UpdateKEKKey(kekKey)
	}

	kekk, err := hex.DecodeString(viper.GetString("KEKkey"))
	if err != nil {
		return fmt.Errorf("error decoding KEK key: %v", err)
	}

	dek, err := genekeys.GenerateDek()
	if err != nil {
		return fmt.Errorf("error generating DEK: %v", err)
	}

	if filePath != "" {
		if err := encryption.EncryptFile(filePath, dek); err != nil {
			return fmt.Errorf("error encrypting file: %v", err)
		}
	}

	encryptedDek, err := encryption.EncryptData(dek, kekk)
	if err != nil {
		return fmt.Errorf("error encrypting DEK: %v", err)
	}

	dekString := hex.EncodeToString(encryptedDek)

	metadata := map[string]string{
		DEKKeyMetadataKey: dekString,
	}

	return handleUpload(accessKeyID, secretAccessKey, bcktName, fmt.Sprintf("%s.enc", filePath), objectkey, metadata)
}

func handleUpload(accessKeyID, secretAccessKey, bcktName, filePath, objectkey string, metadata map[string]string) error {
	accessKeyID, secretAccessKey, region, err := common.HandleCredentials(accessKeyID, secretAccessKey)
	if err != nil {
		return err
	}

	return uploading.UploadFile(region, accessKeyID, secretAccessKey, bcktName, filePath, objectkey, metadata)
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().BoolP("encrypt", "e", false, "Encrypt file before uploading")
}
