package main

import (
	"context"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(client kubernetes.Clientset, namespace string, phase string) (*v1.PodList, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if phase != "" {
		var filteredPods []v1.Pod
		for _, pod := range pods.Items {
			if strings.EqualFold(string(pod.Status.Phase), phase) {
				filteredPods = append(filteredPods, pod)
			}
		}
		pods.Items = filteredPods
	}

	return pods, nil
}
