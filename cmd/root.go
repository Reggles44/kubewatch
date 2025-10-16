package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/reggles44/kubewatch/internal/display"
	"github.com/reggles44/kubewatch/internal/kube"
	"github.com/reggles44/kubewatch/internal/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagNamespace string

var rootCmd = &cobra.Command{
	Use:       "kubewatch",
	Short:     "Watch a kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		dataCh := make(chan kube.WatchEvent)
		defer close(dataCh)

		res, ok := resources.Get(args[0])
		if !ok {
			log.Fatal("args not valid")
		}

		disp := display.New(dataCh, res)

		go kube.Watch(res.Scheme, flagNamespace, dataCh)

		p := tea.NewProgram(disp, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
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

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
