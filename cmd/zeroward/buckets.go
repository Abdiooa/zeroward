package zeroward

import (
	"fmt"

	listingbuckets "github.com/Abdiooa/zeroward/pkg/zeroward/buckets"
	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "List All Buckets(Folders) stored on the Cloud Storage",
	Long:  `This command is used to list all your buckets(folders) that are stored for you on the cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {

		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")

		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")

		KeyAccessDefined := common.IsNotKeyAccessDefined()
		region := viper.GetString("Region")

		if KeyAccessDefined {
			if accessKeyID == "" || secretAccessKey == "" {
				fmt.Println("Error: Access Key ID and Secret Access Key are required as for your login and password of the Cloud Storage.")
				return
			}

			common.StoreCredentials(accessKeyID, secretAccessKey, region)

			err := listingbuckets.ListBuckets(region, accessKeyID, secretAccessKey)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		} else {
			err := listingbuckets.ListBuckets(region, viper.GetString("AWSAccessKeyID"), viper.GetString("AWSSecretAccessKey"))

			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(bucketsCmd)

}
