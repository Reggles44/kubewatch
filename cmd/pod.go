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

var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "Watch all pods in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		stopCh := make(chan struct{})
		defer close(stopCh)

		dataCh := make(chan [2]interface{})
		defer close(dataCh)

		model := display.New(
			dataCh,
			func(obj interface{}, now time.Time) []string {
				pod := obj.(*corev1.Pod)
				return []string{
					pod.Name,
					string(pod.Status.Phase),
					fmt.Sprintf("%-8v", now.Sub(pod.Status.StartTime.Time).Round(time.Second)),
					pod.Status.PodIP,
				}
			},
			func(obj interface{}) string {
				pod := obj.(*corev1.Pod)
				return pod.Name
			},
		)

		go kube.Watch(
			flagNamespace, dataCh, stopCh,
			func(factory informers.SharedInformerFactory) cache.SharedIndexInformer {
				return factory.Core().V1().Pods().Informer()
			})

		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}
