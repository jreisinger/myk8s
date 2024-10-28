package graph

import (
	"log"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

func CreateDotFile(myServices []MyService, svcNames []string) {
	type Object struct{ Name string }
	objectHash := func(o Object) string {
		return o.Name
	}

	g := graph.New(objectHash)

	for _, svc := range myServices {
		for _, pod := range svc.Pods {
			for _, c := range pod.Containers {
				for _, e := range c.Env {
					if s := envVarReferencesService(e, myServices); s != nil {
						if contains(svcNames, s.Name) {
							_ = g.AddVertex(Object{Name: c.Name})
							_ = g.AddVertex(Object{Name: s.Name})
							_ = g.AddEdge(c.Name, s.Name)
						}
					}
				}
			}
		}
	}

	file, _ := os.Create("my-graph.gv")
	if err := draw.DOT(g, file); err != nil {
		log.Fatal(err)
	}
}

func contains(ss []string, s string) bool {
	for _, s1 := range ss {
		if s1 == s {
			return true
		}
	}
	return false
}
