package kube

import (
	"context"
	"fmt"

	"log/slog"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListNamespaces(client *kubernetes.Clientset) (*v1.NamespaceList, error) {
	slog.Debug("Get Kubernetes Namespaces")
	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return namespaces, nil
}

func NamespaceIsSelected(n v1.Namespace) (bool, error) {
	if n.GetAnnotations() != nil && n.GetAnnotations()["pod-reloader/auto"] == "true" {
		return true, nil
	}
	return false, nil
}
