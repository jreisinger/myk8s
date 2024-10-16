package services

import (
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestToMySvc(t *testing.T) {
	svc := v1.Service{}
	newName := "foo"
	mySvc := ToMySvc(svc, "", newName)
	if mySvc.Metadata.Name != newName {
		t.Errorf("got %s but want %s", mySvc.Metadata.Name, newName)
	}
}
