package main

import (
	"encoding/json"
	"testing"
)

type (
	Gettable interface {
		GetField() map[string]interface{}
	}

	Naked struct {
		Field map[string]interface{}
	}
	DefaultTag struct {
		Field map[string]interface{} `json:"field"`
	}
	Omitempty struct {
		Field map[string]interface{} `json:"field,omitempty"`
	}
)

func (s Naked) GetField() map[string]interface{}      { return s.Field }
func (s DefaultTag) GetField() map[string]interface{} { return s.Field }
func (s Omitempty) GetField() map[string]interface{}  { return s.Field }

func TestZeroMarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		target   Gettable
		expected string
	}{
		// Context: of a naked field
		{
			name:     "nil of a naked field",
			target:   Naked{},
			expected: `{"Field":null}`,
		},
		{
			name:     "{} of a naked field",
			target:   Naked{Field: make(map[string]interface{})},
			expected: `{"Field":{}}`,
		},
		// Context: of a field JSON-tagged "field"
		{
			name:     `nil of a field JSON-tagged "field"`,
			target:   DefaultTag{},
			expected: `{"field":null}`,
		},
		{
			name:     `{} of a field JSON-tagged "field"`,
			target:   DefaultTag{Field: make(map[string]interface{})},
			expected: `{"field":{}}`,
		},
		// Context: of a field JSON-tagged "field,omitempty"
		{
			name:     `nil of a field JSON-tagged "field,omitempty"`,
			target:   Omitempty{},
			expected: `{}`,
		},
		{
			name:     `[] of a field JSON-tagged "field,omitempty"`,
			target:   Omitempty{Field: make(map[string]interface{})},
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
				t.Fatalf("isNil json is %s but generated one is %s", expected, result)
			}
		})
	}
}

func TestZeroUnmarshal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		input  string
		target Gettable
		isNil  bool
	}{
		// Context: to a struct with the naked field
		{
			name:   `omitted field to a struct with the naked field`,
			input:  `{}`,
			target: &Naked{},
			isNil:  true,
		},
		{
			name:   `"" to a struct with the naked field`,
			input:  `{"Field":{}}`,
			target: &Naked{},
			isNil:  false,
		},
		{
			name:   `null to a struct with the naked field`,
			input:  `{"Field":null}`,
			target: &Naked{},
			isNil:  true,
		},
		// Context: to a struct with the field JSON-tagged "field"
		{
			name:   `omitted field to a struct with the field JSON-tagged "field"`,
			input:  `{}`,
			target: &DefaultTag{},
			isNil:  true,
		},
		{
			name:   `"" to a struct with the field JSON-tagged "field"`,
			input:  `{"field":{}}`,
			target: &DefaultTag{},
			isNil:  false,
		},
		{
			name:   `null to a struct with the field JSON-tagged "field"`,
			input:  `{"field":null}`,
			target: &DefaultTag{},
			isNil:  true,
		},
		// Context: to a struct with the field JSON-tagged "field,omitempty"
		{
			name:   `omitted field to a struct with the field JSON-tagged "field,omitempty"`,
			input:  `{}`,
			target: &Omitempty{},
			isNil:  true,
		},
		{
			name:   `"" to a struct with the field JSON-tagged "field,omitempty"`,
			input:  `{"field":{}}`,
			target: &Omitempty{},
			isNil:  false,
		},
		{
			name:   `null to a struct with the field JSON-tagged "field,omitempty"`,
			input:  `{"field":null}`,
			target: &Omitempty{},
			isNil:  true,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(v.input), v.target); err != nil {
				t.Fatalf("failed to parse JSON: %s", err)
			}

			result := v.target.GetField()
			if (result == nil) != v.isNil {
				t.Fatalf(`isNil field value is "%v", but detected value is "%s"`, v.isNil, result)
			}
		})
	}
}
