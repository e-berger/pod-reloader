package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func findKubeconfig() string {
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath != "" {
		return kubeConfigPath
	}

	userHomeDir, err := os.UserHomeDir()
	if err == nil && userHomeDir != "" {
		kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
		return kubeConfigPath
	}

	return ""
}

func ListNamespaces(client kubernetes.Interface) (*v1.NamespaceList, error) {
	slog.Debug("Get Kubernetes Namespaces")
	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting namespaces: %v", err)
	}
	return namespaces, nil
}

func getClient() (*kubernetes.Clientset, error) {

	kubeConfigPath := findKubeconfig()
	slog.Info("findKubeconfig", "kubeconfig", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("Error getting kubernetes config", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Error getting kubernetes config", err)
	}
	return clientset, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	clientset, err := getClient()
	if err != nil {
		slog.Error("Error finding kubeconfig", "error", err)
		panic(err)
	}
	listNamespace, err := ListNamespaces(clientset)
	if err != nil {
		slog.Error("Error finding namespace", "error", err)
		panic(err)
	}
	slog.Debug("List", "namespace", listNamespace)
}
