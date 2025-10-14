package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/reggles44/kubewatch/internal/display"
	"github.com/reggles44/kubewatch/internal/kube"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/spf13/cobra"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Watch all nodes in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		stopCh := make(chan struct{})
		defer close(stopCh)

		dataCh := make(chan [2]interface{})
		defer close(dataCh)

		model := display.New(
			dataCh,
			func(obj interface{}, now time.Time) []string {
				node := obj.(*corev1.Node)
				params := []string{
					node.Name,
					fmt.Sprintf("%v Pods", node.Status.Allocatable.Pods()),
					fmt.Sprintf("%v Cores", node.Status.Allocatable.Cpu()),
					fmt.Sprintf("%v Mem", node.Status.Allocatable.Memory()),
				}
				return params
			},
			func(obj interface{}) string {
				node := obj.(*corev1.Node)
				return node.Name
			},
		)

		go kube.Watch(
			flagNamespace, dataCh, stopCh,
			func(factory informers.SharedInformerFactory) cache.SharedIndexInformer {
				return factory.Core().V1().Nodes().Informer()
			})

		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}
