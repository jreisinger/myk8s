package main

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetSecrets(client kubernetes.Clientset, namespace string) (*v1.SecretList, error) {
	return client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
}
