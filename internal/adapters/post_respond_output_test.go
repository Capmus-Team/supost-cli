package adapters

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestRenderPostRespondResult(t *testing.T) {
	var out bytes.Buffer
	result := domain.PostRespondResult{
		DryRun:       true,
		PostID:       130031908,
		PostEmail:    "wientjes@alumni.stanford.edu",
		ReplyTo:      "gwientjes@gmail.com",
		MessageSaved: false,
		EmailSent:    false,
		Subject:      "SUpost - gwientjes@gmail.com response: Looking for a buddy to go to the movies",
		Body:         "Reply to: gwientjes@gmail.com",
	}

	if err := RenderPostRespondResult(&out, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	plain := out.String()
	for _, needle := range []string{
		"[DRY RUN] post respond",
		"post_id: 130031908",
		"post_email: wientjes@alumni.stanford.edu",
		"reply_to: gwientjes@gmail.com",
		"subject: SUpost - gwientjes@gmail.com response: Looking for a buddy to go to the movies",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in output", needle)
		}
	}
}
