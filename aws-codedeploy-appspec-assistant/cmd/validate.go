package cmd

import (
	"aws-codedeploy-appspec-assistant/pkg"
	"fmt"

	"github.com/spf13/cobra"
)

var filePath string
var computePlatform string

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a CodeDeploy AppSpec file",
	Long:  `Validate a CodeDeploy AppSpec file that is locally saved.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validate called")
		assistant.ValidateAppSpec(filePath, computePlatform)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVar(&filePath, "filePath", "", "FilePath of AppSpec file to validate")
	validateCmd.PersistentFlags().StringVar(&computePlatform, "computePlatform", "", "computePlatform of AppSpec file (server, lambda, ecs)")

	validateCmd.MarkFlagRequired("filePath")
	validateCmd.MarkFlagRequired("computePlatform")
}
