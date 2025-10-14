package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/reggles44/kubewatch/internal/display"
	"github.com/reggles44/kubewatch/internal/kube"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/spf13/cobra"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Watch all deployments in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		stopCh := make(chan struct{})
		defer close(stopCh)

		dataCh := make(chan [2]interface{})
		defer close(dataCh)

		model := display.New(
			dataCh,
			func(obj interface{}, now time.Time) []string {
				deployment := obj.(*appsv1.Deployment)
				return []string{
					deployment.Name,
					fmt.Sprintf("%v/%v", deployment.Status.AvailableReplicas, deployment.Status.Replicas),
				}
			},
			func(obj interface{}) string {
				deployment := obj.(*appsv1.Deployment)
				return deployment.Name
			},
		)

		go kube.Watch(
			flagNamespace, dataCh, stopCh,
			func(factory informers.SharedInformerFactory) cache.SharedIndexInformer {
				return factory.Apps().V1().Deployments().Informer()
			})

		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}
