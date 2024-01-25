package clsdapp

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/encryption"
	"github.com/Abdiooa/CLSDAPP/pkg/clsdapp/genekeys"
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

		firstEncryption := isFirstEncryption()
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
			updateKEKKey(kekKey)
		}
		// Use the existing KEKKey
		kekk := viper.GetString("KEKkey")

		kekBytes, err := hex.DecodeString(kekk)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// copy(kekBytes[:], viper.GetString("KEKkey"))
		// kekKey := kekBytes[:]

		dek, err := genekeys.GenerateDek()
		cobra.CheckErr(err)
		if filePath != "" {
			if err := encryption.Encrypt(filePath, dek); err != nil {
				fmt.Println("Error encrypting File:", err)
				return
			}
		}
		if err := encryption.Encrypt(dek, kekBytes); err != nil {
			fmt.Println("Error encrypting DEK:", err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringP("passphrase", "p", "", "Passphrase for encryption required for the first encryption")

	// encryptCmd.Flags().StringP("filePath", "f", "", "Path of the file that you want to encrypt")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func isFirstEncryption() bool {
	kekKey := viper.GetString("KEKkey")
	return kekKey == ""
}

func updateKEKKey(kek []byte) {

	kekString := hex.EncodeToString(kek)
	viper.Set("KEKkey", kekString)
	if err := viper.WriteConfig(); err != nil {
		fmt.Println("Error updating KEK key in the config file:", err)
		os.Exit(1)
	}
}
