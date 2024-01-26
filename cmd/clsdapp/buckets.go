/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package clsdapp

import (
	"fmt"

	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/common"
	listingbuckets "github.com/Abdiooa/CLSDAPP/pkg/clsdapp/listingbuckets"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bucketsCmd represents the buckets command
var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "List All Buckets(Folders) stored on the Yandex Cloud Storage",
	Long:  `This command is used to list all your buckets(folders) that are stored for you on the yandex cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {

		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")

		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")

		KeyAccessDefined := common.IsNotKeyAccessDefined()
		region := viper.GetString("Region")

		if KeyAccessDefined {
			if accessKeyID == "" || secretAccessKey == "" {
				fmt.Println("Error: Access Key ID and Secret Access Key are required as for your login and password of the Yandex Cloud Storage.")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bucketsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bucketsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
