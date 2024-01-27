package zeroward

import (
	"encoding/hex"
	"fmt"

	"github.com/Abdiooa/zeroward/pkg/zeroward/common"
	"github.com/Abdiooa/zeroward/pkg/zeroward/encryption"
	"github.com/Abdiooa/zeroward/pkg/zeroward/genekeys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file Locally",
	Long:  `This command is used to perform client-side data encryption, it uses an advanced algorithms to encrypt user files`,
	Run: func(cmd *cobra.Command, args []string) {

		passphrase, _ := cmd.Flags().GetString("passphrase")

		filePath, _ := cmd.Flags().GetString("filePath")

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
		// Use the existing KEKKey
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
		if err := encryption.EncryptKey(dek, kekBytes, filePath); err != nil {
			fmt.Println("Error encrypting DEK:", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}
