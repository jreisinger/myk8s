package get

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ReplicaSets(client kubernetes.Clientset, namespace string) (*appsv1.ReplicaSetList, error) {
	return client.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
}
