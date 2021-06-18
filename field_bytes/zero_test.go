package main

import (
	"encoding/json"
	"testing"
)

type (
	Gettable interface {
		GetField() []byte
	}

	Naked struct {
		Field []byte
	}
	DefaultTag struct {
		Field []byte `json:"field"`
	}
	Omitempty struct {
		Field []byte `json:"field,omitempty"`
	}
)

func (s Naked) GetField() []byte      { return s.Field }
func (s DefaultTag) GetField() []byte { return s.Field }
func (s Omitempty) GetField() []byte  { return s.Field }

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
			expected: `{"Field":null}`,
		},
		{
			name:     "a default JSON tag",
			target:   DefaultTag{},
			expected: `{"field":null}`,
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
