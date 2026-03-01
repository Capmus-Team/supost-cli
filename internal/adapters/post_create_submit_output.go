package adapters

import (
	"fmt"
	"io"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// RenderPostCreateSubmitResult renders submit-mode confirmation.
func RenderPostCreateSubmitResult(w io.Writer, result domain.PostCreateSubmitResult) error {
	mode := "SUBMIT"
	if result.DryRun {
		mode = "DRY RUN"
	}
	lines := []string{
		fmt.Sprintf("[%s] post create", mode),
		fmt.Sprintf("post_id: %d", result.PostID),
		fmt.Sprintf("photo_count: %d", result.PhotoCount),
		fmt.Sprintf("publish_url: %s", result.PublishURL),
		fmt.Sprintf("email_to: %s", result.EmailTo),
		fmt.Sprintf("email_sent: %t", result.EmailSent),
		fmt.Sprintf("subject: %s", result.Subject),
		"",
		result.Body,
	}
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}
