package kube

import (
	"context"
	"fmt"
	"log/slog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Rollout(clientset *kubernetes.Clientset, pod v1.Pod, namespace string) error {
	labels := pod.Labels

	if _, ok := labels["pod-template-hash"]; ok {
		delete(labels, "pod-template-hash")
	}

	// Query deployments with matching labels
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: labels,
		}),
	})
	if err != nil {
		panic(err.Error())
	}

	if len(deployments.Items) == 0 {
		return fmt.Errorf("no deployment found for the pod")
	}

	// Perform rollout restart on the first deployment found
	deploymentName := deployments.Items[0].Name
	slog.Info("Starting rolling out", "deployment", deploymentName)

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting deployment: %v", err)
	}

	retryOptions := retry.DefaultRetry
	retryErr := retry.RetryOnConflict(retryOptions, func() error {
		// Fetch the latest deployment state
		latestDeployment, err := clientset.AppsV1().Deployments(deployment.Namespace).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		// Update the annotations to trigger a rollout restart
		annotations := latestDeployment.Spec.Template.Annotations
		if annotations == nil {
			annotations = make(map[string]string)
		}
		annotations["kubectl.kubernetes.io/restartedAt"] = metav1.Now().String()
		latestDeployment.Spec.Template.SetAnnotations(annotations)

		_, updateErr := clientset.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), latestDeployment, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	return nil
}
