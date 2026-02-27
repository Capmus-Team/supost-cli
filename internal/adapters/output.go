// Package adapters handles external side effects and output rendering.
// See AGENTS.md §2.4, §5.5.
package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Render outputs data in the specified format.
// Default is JSON (the eventual consumer is a web frontend).
func Render(format string, data interface{}) error {
	return RenderTo(os.Stdout, format, data)
}

// RenderTo outputs data to a specific writer (testable).
func RenderTo(w io.Writer, format string, data interface{}) error {
	switch format {
	case "json", "":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case "text":
		_, err := fmt.Fprintf(w, "%v\n", data)
		return err
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
