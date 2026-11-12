package pages

import (
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

func fieldsText(fields []view.Field) string {
	lines := make([]string, 0, len(fields))

	for _, field := range fields {
		lines = append(lines, field.Name+": "+field.Value)
	}

	return strings.Join(lines, "\n")
}
