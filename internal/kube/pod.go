package kube

import v1 "k8s.io/api/core/v1"

func IsReady(pod v1.Pod) bool {
	for _, c := range pod.Status.Conditions {
		if c.Type == v1.PodReady && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}
