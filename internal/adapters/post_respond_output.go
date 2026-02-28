package adapters

import (
	"fmt"
	"io"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

// RenderPostRespondResult renders post-respond confirmation and email preview.
func RenderPostRespondResult(w io.Writer, result domain.PostRespondResult) error {
	mode := "SEND"
	if result.DryRun {
		mode = "DRY RUN"
	}
	lines := []string{
		fmt.Sprintf("[%s] post respond", mode),
		fmt.Sprintf("post_id: %d", result.PostID),
		fmt.Sprintf("post_email: %s", result.PostEmail),
		fmt.Sprintf("reply_to: %s", result.ReplyTo),
		fmt.Sprintf("message_id: %d", result.MessageID),
		fmt.Sprintf("message_saved: %t", result.MessageSaved),
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
