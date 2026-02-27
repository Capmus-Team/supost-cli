package repository

import "testing"

func TestClampRecentLimit(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{name: "negative", input: -1, want: maxRecentActivePosts},
		{name: "zero", input: 0, want: maxRecentActivePosts},
		{name: "within range", input: 12, want: 12},
		{name: "too high", input: 999, want: maxRecentActivePosts},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clampRecentLimit(tt.input)
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestParseRecentActivePosts(t *testing.T) {
	raw := "130\x1fperson@stanford.edu\x1fBike for sale\x1f1\x1f1700000000\x1f1700000001\x1f100.00\x1ft\x1ff"

	posts, err := parseRecentActivePosts(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("got %d posts, want 1", len(posts))
	}

	post := posts[0]
	if post.ID != 130 {
		t.Fatalf("id got %d, want %d", post.ID, 130)
	}
	if post.Email != "person@stanford.edu" {
		t.Fatalf("email got %q", post.Email)
	}
	if post.Name != "Bike for sale" {
		t.Fatalf("name got %q", post.Name)
	}
	if !post.HasPrice {
		t.Fatalf("expected has_price=true")
	}
	if post.HasImage {
		t.Fatalf("expected has_image=false")
	}
	if post.Price != 100 {
		t.Fatalf("price got %v, want 100", post.Price)
	}
	if post.TimePosted != 1700000000 {
		t.Fatalf("time_posted got %d", post.TimePosted)
	}
	if post.TimePostedAt.IsZero() {
		t.Fatalf("expected non-zero time_posted_at")
	}
}
