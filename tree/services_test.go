package tree

import (
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestEnvVarReferencesService(t *testing.T) {
	tests := []struct {
		name       string
		envVar     v1.EnvVar
		myServices []MyService
		expected   *MyService
	}{
		{
			name:   "no reference",
			envVar: v1.EnvVar{Name: "ENV_VAR", Value: "service3"},
			myServices: []MyService{
				{Name: "service1"},
				{Name: "service2"},
			},
			expected: nil,
		},
		{
			name:   "single reference",
			envVar: v1.EnvVar{Name: "ENV_VAR", Value: "service1"},
			myServices: []MyService{
				{Name: "service1"},
				{Name: "service2"},
			},
			expected: &MyService{Name: "service1"},
		},
		{
			name:   "multiple references, longest name",
			envVar: v1.EnvVar{Name: "ENV_VAR", Value: "service-name"},
			myServices: []MyService{
				{Name: "service"},
				{Name: "service-name"},
			},
			expected: &MyService{Name: "service-name"},
		},
		{
			name:       "no services",
			envVar:     v1.EnvVar{Name: "ENV_VAR", Value: "some-value"},
			myServices: []MyService{},
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := envVarReferencesService(tt.envVar, tt.myServices)
			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			} else if result != nil && tt.expected != nil && result.Name != tt.expected.Name {
				t.Errorf("expected %v, got %v", tt.expected.Name, result.Name)
			}
		})
	}
}
