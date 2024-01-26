/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package clsdapp

import (
	"fmt"

	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/common"
	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/listingobjects"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ObjectsCmd represents the Objects command
var ObjectsCmd = &cobra.Command{
	Use:   "Objects",
	Short: "List All Objects(Files/Images/Docs) stored on the Yandex Cloud Storage",
	Long:  `This command is used to list all your Objects (Files/Images/Docs) that are stored for you on the yandex cloud storage.`,
	Run: func(cmd *cobra.Command, args []string) {

		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")

		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")

		bcktName, _ := cmd.Flags().GetString("bcktname")

		KeyAccessDefined := common.IsNotKeyAccessDefined()
		region := viper.GetString("Region")

		if KeyAccessDefined {
			if accessKeyID == "" || secretAccessKey == "" || bcktName == "" {
				fmt.Println("Error: Access Key ID and Secret Access Key are required as for your login and password of the Yandex Cloud Storage, also the bucket name is required!")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ObjectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ObjectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
