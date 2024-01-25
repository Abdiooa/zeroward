package clsdapp

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/decryption"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Command to decrypt user encrypted file",
	Long:  `This command is used  to decrypt user encrypted files, so he can read and have the file `,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("filePath")

		dekkeyPath, _ := cmd.Flags().GetString("dekkey")

		kekk := viper.GetString("KEKkey")

		kekBytes, err := hex.DecodeString(kekk)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		dekkey, err := decryption.DecryptKey(dekkeyPath, kekBytes)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if filePath != "" {
			if err := decryption.DecryptFile(filePath, dekkey); err != nil {
				fmt.Println("Error encrypting File:", err)
				return
			}
		}
		if err := os.Remove(dekkeyPath); err != nil {
			fmt.Println("Error:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringP("dekkey", "k", "", "DEK Key to decrypt the file please!")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
