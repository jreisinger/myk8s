package main

import (
	"fmt"
	"path"
	"strings"

	"k8s.io/client-go/kubernetes"
)

type Object struct {
	Name string
	Kind string
}

func (o Object) String() string {
	return path.Join(o.Kind, o.Name)
}

type Graph []Object

func (objects Graph) String() string {
	var ss []string
	indentForKind := make(map[string]string)
	for i, object := range objects {
		indent, ok := indentForKind[object.Kind]
		if !ok {
			if i > 0 {
				indent = strings.Repeat(" ", 2*(i-1)) + "└─"
			}
			indentForKind[object.Kind] = indent
		}
		s := fmt.Sprintf("%s%s", indent, object)
		ss = append(ss, s)
	}

	// Remove empty objects.
	var out []string
	for _, s := range ss {
		if s != "" {
			out = append(out, s)
		}
	}

	return strings.Join(out, "\n")
}

type Graphs []Graph

func (graphs Graphs) String() string {
	// Make duplicate objects empty.
	seen := make(map[Object]bool)
	var ss []string
	for i, graph := range graphs {
		if seen[graph[0]] {
			graphs[i][0] = Object{}
		} else {
			seen[graph[0]] = true
		}
		ss = append(ss, graph.String())
	}

	return strings.Join(ss, "\n")
}

func GetGraphs(client kubernetes.Clientset, namespace, phase string) (Graphs, error) {
	var graphs Graphs

	pods, err := GetPods(client, namespace, phase)
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		var objects Graph

		// Above Pod
		for _, or := range pod.ObjectMeta.OwnerReferences {
			objects = append(objects, Object{Name: or.Name, Kind: or.Kind})
		}

		// Pod
		objects = append(objects, Object{Name: pod.Name, Kind: "Pod"})

		// Below Pod
		for _, v := range pod.Spec.Volumes {
			objects = append(objects, Object{Name: v.Name, Kind: "Volume"})
			if v.PersistentVolumeClaim != nil {
				objects = append(objects, Object{Name: v.PersistentVolumeClaim.ClaimName, Kind: "PersistentVolumeClaim"})
			}
		}

		graphs = append(graphs, objects)
	}

	return graphs, nil
}
