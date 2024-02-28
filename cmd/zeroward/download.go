package zeroward

import (
	"fmt"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/Abdiooa/zeroward/pkg/zeroward/downloading"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Command to download Files from the cloud.",
	Long:  `This command is used to download User Objects(Files) from the cloud storage, also gives you an readable files if the files were encrypted.`,
	Run: func(cmd *cobra.Command, args []string) {
		accessKeyID, _ := cmd.Flags().GetString("accessKeyID")
		secretAccessKey, _ := cmd.Flags().GetString("secretAccessKey")
		bcktName, _ := cmd.Flags().GetString("bucketName")
		filePath, _ := cmd.Flags().GetString("filePath")
		objectkey, _ := cmd.Flags().GetString("objectkey")
		removeAfterDownload, _ := cmd.Flags().GetBool("remove")
		decryptWhileDownloading, _ := cmd.Flags().GetBool("decrypt")

		if decryptWhileDownloading {
			if err := handleDownload(accessKeyID, secretAccessKey, bcktName, filePath, objectkey, removeAfterDownload); err != nil {
				fmt.Println("Error:", err)
				return
			}
		} else {
			accessKeyID, secretAccessKey, region, err := common.HandleCredentials(accessKeyID, secretAccessKey)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if err := downloading.DownloadNormalObject(region, accessKeyID, secretAccessKey, bcktName, filePath, objectkey, removeAfterDownload); err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
	},
}

func handleDownload(accessKeyID, secretAccessKey, bcktName, filePath, objectkey string, removeAfterDownload bool) error {
	accessKeyID, secretAccessKey, region, err := common.HandleCredentials(accessKeyID, secretAccessKey)
	if err != nil {
		return err
	}

	return downloading.DownloadObject(region, accessKeyID, secretAccessKey, bcktName, filePath, objectkey, removeAfterDownload)
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("remove", "r", false, "Remove the file from the cloud storage after successful download")
	downloadCmd.Flags().BoolP("decrypt", "d", false, "Decrypt file while downloading")
}
