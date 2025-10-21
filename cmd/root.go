package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/reggles44/kubewatch/pkg/display"
	"github.com/reggles44/kubewatch/pkg/kube"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/watch"
)

func NewCmd() (*cobra.Command, error) {
	configFlags := kube.NewConfigFlags()

	cmd := &cobra.Command{
		Use:   "kubewatch",
		Short: "Watch a kubernetes cluster",
		Run: func(cmd *cobra.Command, args []string) {
			// Get Resource
			r, err := kube.GetResource(args, configFlags)
			if err != nil {
				panic(err)
			}

			// Make Channel
			dataCh := make(chan watch.Event)
			defer close(dataCh)

			// Start watching
			go kube.WatchRuntimeObject(r, dataCh)

			d := display.New(dataCh)

			p := tea.NewProgram(d, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				log.Fatal(err)
			}
		},
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	return cmd, nil
}
