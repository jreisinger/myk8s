package get

import (
	"context"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Deployments(client kubernetes.Clientset, namespace string) (*appsv1.DeploymentList, error) {
	return client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
}

func Pods(client kubernetes.Clientset, namespace string, phase string) (*corev1.PodList, error) {
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if phase != "" {
		var filteredPods []corev1.Pod
		for _, pod := range pods.Items {
			if strings.EqualFold(string(pod.Status.Phase), phase) {
				filteredPods = append(filteredPods, pod)
			}
		}
		pods.Items = filteredPods
	}

	return pods, nil
}

func ReplicaSets(client kubernetes.Clientset, namespace string) (*appsv1.ReplicaSetList, error) {
	return client.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
}

func Services(client kubernetes.Clientset, namespace string) (*corev1.ServiceList, error) {
	return client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
}
