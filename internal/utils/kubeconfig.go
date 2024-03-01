package utils

import (
	"os"
	"path/filepath"
)

func FindKubeconfig() string {
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
