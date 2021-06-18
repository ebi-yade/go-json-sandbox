package main

import (
	"encoding/json"
	"testing"
)

type (
	Gettable interface {
		GetField() json.RawMessage
	}

	Naked struct {
		Field json.RawMessage
	}
	DefaultTag struct {
		Field json.RawMessage `json:"field"`
	}
	Omitempty struct {
		Field json.RawMessage `json:"field,omitempty"`
	}
)

func (s Naked) GetField() json.RawMessage      { return s.Field }
func (s DefaultTag) GetField() json.RawMessage { return s.Field }
func (s Omitempty) GetField() json.RawMessage  { return s.Field }

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
			name:     "[] of a naked field",
			target:   Naked{Field: json.RawMessage("")},
			expected: `{"Field":""}`,
		},
		// Context: of a field JSON-tagged "field"
		{
			name:     `nil of a field JSON-tagged "field"`,
			target:   DefaultTag{},
			expected: `{"field":null}`,
		},
		{
			name:     `[] of a field JSON-tagged "field"`,
			target:   Naked{Field: json.RawMessage("")},
			expected: `{"Field":""}`,
		},
		// Context: of a field JSON-tagged "field,omitempty"
		{
			name:     `nil of a field JSON-tagged "field,omitempty"`,
			target:   Omitempty{},
			expected: `{}`,
		},
		{
			name:     `[] of a field JSON-tagged "field,omitempty"`,
			target:   Omitempty{Field: json.RawMessage("")},
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
		expected json.RawMessage
	}{
		// Context: to a struct with the naked field
		{
			name:     `omitted field to a struct with the naked field`,
			input:    `{}`,
			target:   &Naked{},
			expected: nil,
		},
		{
			name:     `"" to a struct with the naked field`,
			input:    `{"Field":""}`,
			target:   &Naked{},
			expected: json.RawMessage(""),
		},
		{
			name:     `null to a struct with the naked field`,
			input:    `{"Field":null}`,
			target:   &Naked{},
			expected: nil,
		},
		// Context: to a struct with the field JSON-tagged "field"
		{
			name:     `omitted field to a struct with the field JSON-tagged "field"`,
			input:    `{}`,
			target:   &DefaultTag{},
			expected: nil,
		},
		{
			name:     `"" to a struct with the field JSON-tagged "field"`,
			input:    `{"field":""}`,
			target:   &DefaultTag{},
			expected: json.RawMessage(""),
		},
		{
			name:     `null to a struct with the field JSON-tagged "field"`,
			input:    `{"field":null}`,
			target:   &DefaultTag{},
			expected: nil,
		},
		// Context: to a struct with the field JSON-tagged "field,omitempty"
		{
			name:     `omitted field to a struct with the field JSON-tagged "field,omitempty"`,
			input:    `{}`,
			target:   &Omitempty{},
			expected: nil,
		},
		{
			name:     `"" to a struct with the field JSON-tagged "field,omitempty"`,
			input:    `{"field":""}`,
			target:   &Omitempty{},
			expected: json.RawMessage(""),
		},
		{
			name:     `null to a struct with the field JSON-tagged "field,omitempty"`,
			input:    `{"field":null}`,
			target:   &Omitempty{},
			expected: nil,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			if err := json.Unmarshal(json.RawMessage(v.input), v.target); err != nil {
				t.Fatalf("failed to parse JSON: %s", err)
			}

			expected := strOrNil(v.expected)
			result := strOrNil(v.target.GetField())

			if expected != result {
				t.Fatalf(`expected field value is "%s", but detected value is "%s"`, expected, result)
			}
		})
	}
}

func strOrNil(b json.RawMessage) string {
	if b == nil {
		return "<nil>"
	}

	return string(b)
}
