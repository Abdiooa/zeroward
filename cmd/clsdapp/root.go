/*
Copyright © 2024 Abdi Omar martelluiz125@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package clsdapp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const configFileName = "config.json"
const clsdFolderName = "CLSD"

type Config struct {
	KEKkey string `json:"kekkey"`
}

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "CLSDAPP",
		Short: "Client-Side Data Encryption Cloud Security Application",
		Long: `CLSDAPP is a command-line application developed in Golang that implements a client-centric 
		approach to ensuring the security of data in cloud environments. It provides a robust solution 
		for encrypting client data locally before uploading it to a cloud storage server. The application 
		implements secure transmission protocols, ensuring data remains encrypted during transfer. 
		Once stored, it adheres to security policies set by the cloud service provider, including additional 
		encryption layers, access management, and continuous monitoring. CLSDAPP empowers users to actively 
		participate in securing their data, offering full control and confidence in maintaining the 
		confidentiality and integrity of their information.`,
		// Run: func(cmd *cobra.Command, args []string) {

		// },
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("filePath", "f", "", "Path of the file that you want to encrypt")
}

func initConfig() {

	if cfgFile != "" {

		viper.SetConfigFile(cfgFile)
	} else {

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		clsdFolderPath := filepath.Join(home, ".config", clsdFolderName)
		viper.AddConfigPath(clsdFolderPath)
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}

// CreateConfigFile attempts to create the config file and CLSD folder
func CreateConfigFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user's home directory: %v", err)
	}

	// Create CLSD folder path
	clsdFolderPath := filepath.Join(homeDir, ".config", clsdFolderName)

	if _, err := os.Stat(clsdFolderPath); os.IsNotExist(err) {
		// CLSD folder does not exist, create and generate KEK key
		err := os.Mkdir(clsdFolderPath, 0700) // Set read-write-execute for the owner only
		if err != nil {
			fmt.Printf("Error creating CLSD folder: %v", err)
		}

		configFilePath := filepath.Join(clsdFolderPath, configFileName)
		config := Config{
			KEKkey: "",
		}

		viper.SetDefault("KEKkey", config.KEKkey) // Set default value for KEKkey in viper

		// Save config using viper
		if err := viper.WriteConfigAs(configFilePath); err != nil {
			fmt.Printf("Error writing config file: %v", err)
		}
	}
	return nil
}
