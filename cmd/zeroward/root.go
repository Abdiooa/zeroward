package zeroward

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFileName        = "config.json"
	clsdFolderName        = "zeroward"
	awsFolderName         = ".aws"
	credentialsFileName   = "credentials.json"
	defaultAWSRegion      = "auto"
	defaultAWSAccessKeyID = ""
	defaultAWSSecretKey   = ""
)

type Config struct {
	KEKkey             string `json:"kekkey"`
	Region             string `json:"region"`
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
}

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "zeroward",
		Short: "Zero-Knowledge-Encryption Cloud Storage Security Application",
		Long: `zeroward is a zero-knowledge-ecryption command-line application that secures client data at all stage(locally,
		 during transmission to the cloud, and at rest). It provides a robust solution 
		for encrypting client data locally before uploading it to a cloud storage server. The application 
		implements secure transmission protocols, ensuring data remains encrypted during transfer. 
		Once stored, it adheres to security policies set by the cloud service provider, including additional 
		encryption layers, access management, and continuous monitoring. zeroward empowers users to actively 
		participate in securing their data, offering full control and confidence in maintaining the 
		confidentiality and integrity of their information.`,
		// Run: func(cmd *cobra.Command, args []string) {
		// 	err := createConfigFile()
		// 	if err != nil {
		// 		fmt.Fprintln(os.Stderr, "Error initializing configuration files:", err)
		// 		os.Exit(1)
		// 	}
		// 	fmt.Println("Configuration files initialized successfully.")
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
	rootCmd.PersistentFlags().StringP("accessKeyID", "i", "", "Access Key Id as your Login Key")
	rootCmd.PersistentFlags().StringP("secretAccessKey", "s", "", "Secret Access Key as your Password")
	rootCmd.PersistentFlags().StringP("bucketName", "b", "", "Bucket name out of all your existing buckets!")
	rootCmd.PersistentFlags().StringP("objectkey", "o", "", "objectkey refers to the unique identifier or name of the object(file) with a bucket, it is the path or where in the bucket the file should be stored.")
	rootCmd.PersistentFlags().StringP("passphrase", "p", "", "Passphrase for encryption required for the first encryption")

}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Get user's home directory based on OS
		var home string
		switch runtime.GOOS {
		case "windows":
			home = os.Getenv("USERPROFILE")
		default:
			homeDir, err := os.UserHomeDir()
			cobra.CheckErr(err)
			home = homeDir
		}

		clsdFolderPath := filepath.Join(home, ".config", clsdFolderName)
		viper.AddConfigPath(clsdFolderPath)
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}
}

func CreateConfigFile() error {
	// Get user's home directory based on OS
	var homeDir string
	switch runtime.GOOS {
	case "windows":
		homeDir = os.Getenv("USERPROFILE")
	default:
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user's home directory: %v", err)
		}
	}

	clsdFolderPath := filepath.Join(homeDir, ".config", clsdFolderName)

	if _, err := os.Stat(clsdFolderPath); os.IsNotExist(err) {
		err := os.Mkdir(clsdFolderPath, 0700)
		if err != nil {
			fmt.Printf("Error creating CLSD folder: %v", err)
		}

		configFilePath := filepath.Join(clsdFolderPath, configFileName)
		config := Config{
			KEKkey:             "",
			Region:             defaultAWSRegion,
			AWSAccessKeyID:     defaultAWSAccessKeyID,
			AWSSecretAccessKey: defaultAWSSecretKey,
		}

		viper.SetDefault("KEKkey", config.KEKkey)
		viper.SetDefault("Region", config.Region)
		viper.SetDefault("AWSAccessKeyID", config.AWSAccessKeyID)
		viper.SetDefault("AWSSecretAccessKey", config.AWSSecretAccessKey)

		// Save config using viper
		if err := viper.WriteConfigAs(configFilePath); err != nil {
			fmt.Printf("Error writing config file: %v", err)
		}
	}
	return nil
}
