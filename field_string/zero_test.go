package main

import (
	"encoding/json"
	"testing"
)

func TestZeroMarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		target   interface{}
		expected string
	}{
		{
			name: "a naked field",
			target: struct {
				Field string
			}{Field: ""},
			expected: `{"Field":""}`,
		},
		{
			name: "a default JSON tag",
			target: struct {
				Field string `json:"field"`
			}{Field: ""},
			expected: `{"field":""}`,
		},
		{
			name: "a JSON tag with omitempty",
			target: struct {
				Field string `json:"field,omitempty"`
			}{Field: ""},
			expected: `{}`,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			payload, err := json.Marshal(v.target)
			if err != nil {
				t.Fatalf("failed to encode the struct to JSON: %s", err)
			}

			expected := v.expected
			result := string(payload)
			if expected != result {
				t.Fatalf("expected json is %s but generated one is %s", expected, result)
			}
		})
	}
}
