// Package trouble finds and prints troublesome pods.
package trouble

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type Pod struct {
	Name              string
	StatusPhase       v1.PodPhase
	ContainerStatuses []v1.ContainerStatus
}

func GetPods(podList *v1.PodList) []Pod {
	var myPods []Pod
	for _, pod := range podList.Items {
		var podIsTroublesome bool

		// Pod ...
		if pod.Status.Phase != v1.PodRunning && pod.Status.Phase != v1.PodSucceeded {
			// ... not running and not successfully finished.
			podIsTroublesome = true
		}

		// Container ...
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Running != nil && cs.Ready {
				// ... running and ready.
			} else if cs.State.Terminated != nil && pod.Status.Phase == v1.PodSucceeded {
				// ... successfully finished.
			} else {
				podIsTroublesome = true
			}
		}

		if podIsTroublesome {
			var myPod Pod
			myPod.Name = pod.Name
			myPod.StatusPhase = pod.Status.Phase
			myPod.ContainerStatuses = pod.Status.ContainerStatuses
			myPods = append(myPods, myPod)
		}
	}

	return myPods
}

func PrintPods(pods []Pod) {
	for _, pod := range pods {
		fmt.Printf("%s (%s)\n", pod.Name, pod.StatusPhase)
		for _, cs := range pod.ContainerStatuses {
			switch {
			case cs.State.Waiting != nil:
				fmt.Printf("  └─%s (%s): %s\n", cs.Name, cs.State.Waiting.Reason, cs.State.Waiting.Message)
			case cs.State.Terminated != nil:
				fmt.Printf("  └─%s (%s): %s\n", cs.Name, cs.State.Terminated.Reason, cs.State.Terminated.Message)
			case cs.State.Running != nil:
				if cs.Ready {
					fmt.Printf("  └─%s (%s): %s\n", cs.Name, "Running", "ready")
				} else {
					fmt.Printf("  └─%s (%s): %s\n", cs.Name, "Running", "not ready")
				}
			default:
				fmt.Printf("  └─%s (%s)\n", cs.Name, "Unknown status")
			}
		}
	}
}
