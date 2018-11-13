package cmd

import (
	"fmt"
	"os"

	"github.com/mschurenko/tfgen/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version will be set at build time
var Version string

var cfgFile string
var tfgenConf = ".tfgen.yml"
var s3Config = make(map[string]string)
var environments []string
var stackRx string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tfgen",
	Short:   "Generate Terraform templates",
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "override config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.

func parseConfig() {
	// s3_backend
	bMap, ok := viper.Get("s3_backend").(map[string]interface{})
	if !ok {
		fmt.Println("Error: is s3_backend a map?")
		os.Exit(1)
	}

	for k, v := range bMap {
		s, ok := v.(string)
		if !ok {
			fmt.Println("Error: all values for s3_backend must be strings")
			os.Exit(1)
		}
		s3Config[k] = s
	}

	// environments
	eSlice, ok := viper.Get("environments").([]interface{})
	if !ok {
		fmt.Println("Error: is environments a list?")
		os.Exit(1)
	}
	for _, environment := range eSlice {
		if s, ok := environment.(string); ok {
			environments = append(environments, s)
		} else {
			fmt.Println("Error: environments must be strings")
			os.Exit(1)
		}
	}

	// stack_regexp
	stackRx, ok = viper.Get("stack_regexp").(string)
	if !ok {
		fmt.Println("Error: is stack_regexp a string?")
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		/*
			Search config in configPath with name ".tfgen" (without extension).
			Note: The first added path has precedence
		*/
		tfDirPath, err := utils.FindTfGenPath(tfgenConf)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		viper.AddConfigPath(tfDirPath)

		// Add config path from absolute path
		viper.SetConfigName(".tfgen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// retrieve values from config
	parseConfig()
}
