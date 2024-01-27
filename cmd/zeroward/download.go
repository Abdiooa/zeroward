package zeroward

import (
	"fmt"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/Abdiooa/zeroward/pkg/zeroward/downloading"
	"github.com/Abdiooa/zeroward/pkg/zeroward/genekeys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Command to download Files from the cloud.",
	Long:  `This command is used to download User Objects(Files) from the cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {

		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")

		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")

		passphrase, _ := cmd.Flags().GetString("passphrase")

		bcktName, _ := cmd.Flags().GetString("bcktname")

		filePath, _ := cmd.Flags().GetString("filePath")

		objectkey, _ := cmd.Flags().GetString("objectkey")

		removeAfterDownload, _ := cmd.Flags().GetString("removeAfterDownload")

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

		KeyAccessDefined := common.IsNotKeyAccessDefined()
		region := viper.GetString("Region")
		if KeyAccessDefined {

			if accessKeyID == "" || secretAccessKey == "" || bcktName == "" {
				fmt.Println("Error: Access Key ID and Secret Access Key are required as for your login and password of the Yandex Cloud Storage, also the bucket name is required!")
				return
			}

			if filePath == "" || objectkey == "" {
				fmt.Println("Error: FilePath and ObjectKey are required!")
				return
			}

			common.StoreCredentials(accessKeyID, secretAccessKey, region)

			err := downloading.DownloadObject(region, accessKeyID, secretAccessKey, bcktName, filePath, objectkey, removeAfterDownload)

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

			if filePath == "" || objectkey == "" {
				fmt.Println("Error: FilePath and ObjectKey are required!")
				return
			} else {
				err := downloading.DownloadObject(region, viper.GetString("AWSAccessKeyID"), viper.GetString("AWSSecretAccessKey"), bcktName, filePath, objectkey, removeAfterDownload)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("removeAfterDownload", "r", "", "write (yes/y) for removing the file from the cloud storage after successful download, else don't specify anything")
}
