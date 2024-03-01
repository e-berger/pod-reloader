package process

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/e-berger/pod-reloader/internal/imageref"
	"github.com/e-berger/pod-reloader/internal/kube"
	"github.com/e-berger/pod-reloader/internal/registry"
	"github.com/e-berger/pod-reloader/internal/utils"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Process struct {
	Client   *kubernetes.Clientset
	Registry registry.IRegistry
}

func NewProcess(registry registry.IRegistry) (*Process, error) {
	clientset, err := getClient()
	if err != nil {
		return nil, err
	}
	return &Process{
		Client:   clientset,
		Registry: registry,
	}, nil
}

func (p *Process) Tick() error {
	slog.Debug("Starting loop")
	namespaces, err := kube.ListNamespaces(p.Client)
	if err != nil {
		return fmt.Errorf("error finding namespace %v", err)
	}

	var pods *v1.PodList
	for _, namespace := range namespaces.Items {
		selected, err := kube.NamespaceIsSelected(namespace)
		if err != nil {
			slog.Error("Error during namespace scan", "error", err)
			break
		}
		if selected {
			slog.Info("Select namespace", "selected", namespace.GetName())
			pods, err = p.Client.CoreV1().Pods(namespace.GetName()).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				slog.Error("Listing pod", "error", err)
				break
			}
			for _, pod := range pods.Items {
				slog.Info("Selected pod", "pod", pod.GetName())
				if pod.GetAnnotations() != nil && pod.GetAnnotations()["pod-reloader/ignore"] == "true" {
					break
				}
				if kube.IsReady(pod) {
					listImage := imageref.ExtractImages(pod)
					for _, image := range listImage {
						slog.Debug("Image", "image", image)
						digest, err := p.Registry.RetreiveImage(image)
						if err != nil {
							continue
						}
						if digest != image.Digest {
							slog.Info("Reload image", "image", image)
							err = kube.Rollout(p.Client, pod, namespace.GetName())
							if err != nil {
								slog.Error("Error during rollout", "error", err)
							}
							slog.Info("Deployment rollout restarted successfully")
						}
					}
				}
			}
		}
	}
	slog.Debug("End loop")
	return nil
}

func getClient() (*kubernetes.Clientset, error) {
	kubeConfigPath := utils.FindKubeconfig()
	slog.Info("findKubeconfig", "kubeconfig", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("error getting kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting kubernetes config: %v", err)
	}
	return clientset, nil
}
