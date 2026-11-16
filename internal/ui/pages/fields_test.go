package pages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

func TestFieldsText(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		fields   []view.Field
		expected string
	}{
		{
			name:     "no fields make an empty string",
			fields:   []view.Field{},
			expected: "",
		},
		{
			name: "one field becomes one line",
			fields: []view.Field{
				{
					Name:  "SWAP",
					Value: "8.00 GB",
				},
			},
			expected: "SWAP: 8.00 GB",
		},
		{
			name: "every field becomes its own line",
			fields: []view.Field{
				{
					Name:  "Hostname",
					Value: "thinkpad",
				},
				{
					Name:  "OS",
					Value: "Ubuntu",
				},
			},
			expected: "Hostname: thinkpad\nOS: Ubuntu",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := fieldsText(tc.fields)
			assert.Equal(t, tc.expected, result)
		})
	}
}
