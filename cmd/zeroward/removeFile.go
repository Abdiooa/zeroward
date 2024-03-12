package zeroward

import (
	"fmt"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/Abdiooa/zeroward/pkg/zeroward/remove"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove Command to delete a file from cloud storage",
	Long:  `This Command is used to delete a file from the cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {

		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")

		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")

		bcktName, _ := cmd.Flags().GetString("bucketName")

		objectkey, _ := cmd.Flags().GetString("objectkey")

		KeyAccessDefined := common.IsNotKeyAccessDefined()
		region := viper.GetString("Region")

		if KeyAccessDefined {
			if accessKeyID == "" || secretAccessKey == "" || bcktName == "" {
				fmt.Println("Error: Access Key ID and Secret Access Key are required as for your login and password of the Yandex Cloud Storage, also the bucket name is required!")
				return
			}

			if objectkey == "" {
				fmt.Println("Error: ObjectKey is required!")
				return
			}

			common.StoreCredentials(accessKeyID, secretAccessKey, region)
		}

		if bcktName == "" {
			fmt.Println("Error:  the bucket name is required!")
			return
		}

		if !KeyAccessDefined {
			if bcktName != "" {
				if objectkey == "" {
					fmt.Println("Error: ObjectKey is required!")
					return
				} else {

					client, err := common.SetupS3Client(region, viper.GetString("AWSAccessKeyID"), viper.GetString("AWSSecretAccessKey"))
					if err != nil {
						fmt.Println("Error setting up S3 client:", err)
						return
					}
					err = remove.Remove(client, bcktName, objectkey)
					if err != nil {
						fmt.Println("Error removing file from cloud storage:", err)
						return
					}

					fmt.Printf("File '%s' removed successfully from S3://%s/%s\n", objectkey, bcktName, objectkey)

				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
