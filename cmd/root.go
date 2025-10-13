package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/reggles44/kubewatch/internal/display"
	"github.com/reggles44/kubewatch/internal/kube"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
)

var flagNamespace string

var rootCmd = &cobra.Command{
	Use:   "kubewatch <resource>",
	Short: "Watch the status of different kubernetes resources",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		model := display.New(
			kube.PodList(flagNamespace),
			func(pod *v1.Pod, now time.Time) []string {
				return []string{
					pod.Name,
					string(pod.Status.Phase),
					fmt.Sprintf("%-8v", now.Sub(pod.Status.StartTime.Time).Round(time.Second)),
					pod.Status.PodIP,
				}
			},
			func(pod *v1.Pod) string { return pod.Name },
		)

		kube.WatchPods(model, []string{""})

		p := tea.NewProgram(model, tea.WithAltScreen())
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
