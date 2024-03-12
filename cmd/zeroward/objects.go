package zeroward

import (
	"fmt"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	listingobjects "github.com/Abdiooa/zeroward/pkg/zeroward/objects"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ObjectsCmd = &cobra.Command{
	Use:   "objects",
	Short: "List All Objects(Files/Images/Docs) stored on the Cloud Storage",
	Long:  `This command is used to list all your Objects (Files/Images/Docs) that are stored for you on the cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {

		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")

		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")

		bcktName, _ := cmd.Flags().GetString("bucketName")

		KeyAccessDefined := common.IsNotKeyAccessDefined()
		region := viper.GetString("Region")

		if KeyAccessDefined {
			if accessKeyID == "" || secretAccessKey == "" || bcktName == "" {
				fmt.Println("Error: Access Key ID and Secret Access Key are required as for your login and password of the Cloud Storage, also the bucket name is required!")
				return
			}

			common.StoreCredentials(accessKeyID, secretAccessKey, region)

			err := listingobjects.ListObjects(region, accessKeyID, secretAccessKey, bcktName)

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
				err := listingobjects.ListObjects(region, viper.GetString("AWSAccessKeyID"), viper.GetString("AWSSecretAccessKey"), bcktName)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ObjectsCmd)
}
