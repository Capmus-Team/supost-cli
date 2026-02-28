package adapters

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

func TestRenderPostCreateSubmitResult(t *testing.T) {
	var out bytes.Buffer
	result := domain.PostCreateSubmitResult{
		DryRun:     true,
		PostID:     130031999,
		PublishURL: "https://supost.com/post/publish/token",
		EmailTo:    "wientjes@alumni.stanford.edu",
		EmailSent:  false,
		Subject:    "SUpost - Publish your post! Red bike for sale",
		Body:       "Publish your post by pressing:",
	}

	if err := RenderPostCreateSubmitResult(&out, result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	plain := out.String()
	for _, needle := range []string{
		"[DRY RUN] post create",
		"post_id: 130031999",
		"publish_url: https://supost.com/post/publish/token",
		"email_to: wientjes@alumni.stanford.edu",
		"subject: SUpost - Publish your post! Red bike for sale",
	} {
		if !strings.Contains(plain, needle) {
			t.Fatalf("missing %q in output", needle)
		}
	}
}
