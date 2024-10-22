package get

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Deployments(client kubernetes.Clientset, namespace string) (*appsv1.DeploymentList, error) {
	return client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
}
