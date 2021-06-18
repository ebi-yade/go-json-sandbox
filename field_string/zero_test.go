package main

import (
	"encoding/json"
	"testing"
)

type (
	Gettable interface {
		GetField() string
	}

	Naked struct {
		Field string
	}
	DefaultTag struct {
		Field string `json:"field"`
	}
	Omitempty struct {
		Field string `json:"field,omitempty"`
	}
)

func (s Naked) GetField() string      { return s.Field }
func (s DefaultTag) GetField() string { return s.Field }
func (s Omitempty) GetField() string  { return s.Field }

func TestZeroMarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		target   Gettable
		expected string
	}{
		{
			name:     "a naked field",
			target:   Naked{},
			expected: `{"Field":""}`,
		},
		{
			name:     "a default JSON tag",
			target:   DefaultTag{},
			expected: `{"field":""}`,
		},
		{
			name:     "a JSON tag with omitempty",
			target:   Omitempty{},
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

func TestZeroUnmarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		target   Gettable
		expected string
	}{
		// Context: to a struct with the naked field
		{
			name:     `omitted field to a struct with the naked field`,
			input:    `{}`,
			target:   &Naked{},
			expected: "",
		},
		{
			name:     `"" to a struct with the naked field`,
			input:    `{"Field":""}`,
			target:   &Naked{},
			expected: "",
		},
		{
			name:     `null to a struct with the naked field`,
			input:    `{"Field":null}`,
			target:   &Naked{},
			expected: "",
		},
		// Context: to a struct with the field JSON-tagged "field"
		{
			name:     `omitted field to a struct with the field JSON-tagged "field"`,
			input:    `{}`,
			target:   &DefaultTag{},
			expected: "",
		},
		{
			name:     `"" to a struct with the field JSON-tagged "field"`,
			input:    `{"field":""}`,
			target:   &DefaultTag{},
			expected: "",
		},
		{
			name:     `null to a struct with the field JSON-tagged "field"`,
			input:    `{"field":null}`,
			target:   &DefaultTag{},
			expected: "",
		},
		// Context: to a struct with the field JSON-tagged "field,omitempty"
		{
			name:     `omitted field to a struct with the field JSON-tagged "field,omitempty"`,
			input:    `{}`,
			target:   &Omitempty{},
			expected: "",
		},
		{
			name:     `"" to a struct with the field JSON-tagged "field,omitempty"`,
			input:    `{"field":""}`,
			target:   &Omitempty{},
			expected: "",
		},
		{
			name:     `null to a struct with the field JSON-tagged "field,omitempty"`,
			input:    `{"field":null}`,
			target:   &Omitempty{},
			expected: "",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(v.input), v.target); err != nil {
				t.Fatalf("failed to parse JSON: %s", err)
			}

			expected := v.expected
			result := v.target.GetField()
			if expected != result {
				t.Fatalf(`expected field value is "%s", but detected value is "%s"`, expected, result)
			}
		})
	}
}
