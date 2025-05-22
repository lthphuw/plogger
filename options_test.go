package plogger

import (
	"errors"
	"reflect"
	"testing"
)

// testConfig is a test struct for the generic type T.
type testConfig struct {
	Name  string
	Value int
}

// mockLister implements Lister[T] for testing.
type mockLister[T any] struct {
	setters []Setter[T]
}

func (m *mockLister[T]) List() []Setter[T] {
	return m.setters
}

func TestNewOptions(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		opts        []Lister[testConfig]
		expected    *testConfig
		expectedErr string
	}{
		{
			name:        "No options, applies defaults",
			opts:        nil,
			expected:    &testConfig{},
			expectedErr: "",
		},
		{
			name: "Nil options, skips and applies defaults",
			opts: []Lister[testConfig]{
				nil,
				&mockLister[testConfig]{setters: nil},
			},
			expected:    &testConfig{},
			expectedErr: "",
		},
		{
			name: "Valid options, applies setters",
			opts: []Lister[testConfig]{
				&mockLister[testConfig]{
					setters: []Setter[testConfig]{
						func(cfg *testConfig) error {
							cfg.Name = "setter1"
							return nil
						},
						func(cfg *testConfig) error {
							cfg.Value = 100
							return nil
						},
					},
				},
			},
			expected:    &testConfig{Name: "setter1", Value: 100},
			expectedErr: "",
		},
		{
			name: "Setter returns error, propagates error",
			opts: []Lister[testConfig]{
				&mockLister[testConfig]{
					setters: []Setter[testConfig]{
						func(cfg *testConfig) error {
							return errors.New("setter error")
						},
					},
				},
			},
			expected:    nil,
			expectedErr: "setter error",
		},
		{
			name: "Multiple listers with mixed setters",
			opts: []Lister[testConfig]{
				&mockLister[testConfig]{
					setters: []Setter[testConfig]{
						func(cfg *testConfig) error {
							cfg.Name = "setter1"
							return nil
						},
						nil,
					},
				},
				&mockLister[testConfig]{
					setters: []Setter[testConfig]{
						func(cfg *testConfig) error {
							cfg.Value = 200
							return nil
						},
					},
				},
			},
			expected:    &testConfig{Name: "setter1", Value: 200},
			expectedErr: "",
		},
		{
			name: "Nil setter in list, skips nil",
			opts: []Lister[testConfig]{
				&mockLister[testConfig]{
					setters: []Setter[testConfig]{
						nil,
						func(cfg *testConfig) error {
							cfg.Value = 300
							return nil
						},
						nil,
					},
				},
			},

			expected:    &testConfig{Value: 300},
			expectedErr: "",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			cfg, err := newOptions[testConfig](tt.opts...)

			// Verify
			if tt.expectedErr == "" {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error %q, got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("Expected error %q, got %q", tt.expectedErr, err.Error())
				}
			}

			if !reflect.DeepEqual(cfg, tt.expected) {
				t.Errorf("Expected config %v, got %v", tt.expected, cfg)
			}
		})
	}
}
