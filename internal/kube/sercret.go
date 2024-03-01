package kube

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetSecret(client *kubernetes.Clientset, namespace string, secretname string) (*v1.Secret, error) {
	secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret, nil
}
