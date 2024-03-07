package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func FindKubeconfig() (*rest.Config, error) {
	kubeConfigPath := os.Getenv("KUBECONFIG")
	var config *rest.Config
	if kubeConfigPath == "" {
		if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
			slog.Info("Using in-cluster config")
			config, err = rest.InClusterConfig()
			if err != nil {
				return nil, err
			}
			return config, nil
		}

		userHomeDir, err := os.UserHomeDir()
		if err == nil && userHomeDir != "" {
			slog.Info("Using home dir config")
			kubeConfigPath = filepath.Join(userHomeDir, ".kube", "config")
		}
	}

	if kubeConfigPath != "" {
		kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("error getting kubernetes config: %v", err)
		}
		return kubeConfig, nil
	}

	return nil, fmt.Errorf("error getting kubernetes config")
}
