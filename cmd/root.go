package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagNamespace string

var rootCmd = &cobra.Command{
	Use:   "kubewatch",
	Short: "Watch a kubernetes cluster",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(0)
	}
}

func init() {
	// cobra.OnInitialize(config.InitConfig)

	// RootCmd Flags
	rootCmd.Flags().StringVar(&flagNamespace, "namespace", "", "kubernetes namespace")

	rootCmd.AddCommand(deploymentCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(podCmd)

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
