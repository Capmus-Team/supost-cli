package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

type mockPostCreateSubmitRepo struct {
	categories    []domain.Category
	subcategories []domain.Subcategory
	submission    domain.PostCreateSubmission
	persisted     domain.PostCreatePersisted
	savedPhotos   []domain.PostCreateSavedPhoto
	createCalled  bool
	saveCalled    bool
}

func (m *mockPostCreateSubmitRepo) ListCategories(_ context.Context) ([]domain.Category, error) {
	return m.categories, nil
}

func (m *mockPostCreateSubmitRepo) ListSubcategories(_ context.Context) ([]domain.Subcategory, error) {
	return m.subcategories, nil
}

func (m *mockPostCreateSubmitRepo) CreatePendingPost(_ context.Context, submission domain.PostCreateSubmission) (domain.PostCreatePersisted, error) {
	m.createCalled = true
	m.submission = submission
	if m.persisted.PostID == 0 {
		m.persisted.PostID = 130031999
	}
	if m.persisted.AccessToken == "" {
		m.persisted.AccessToken = submission.AccessToken
	}
	if m.persisted.PostedAt.IsZero() {
		m.persisted.PostedAt = submission.PostedAt
	}
	return m.persisted, nil
}

func (m *mockPostCreateSubmitRepo) SavePostPhotos(_ context.Context, photos []domain.PostCreateSavedPhoto) error {
	m.saveCalled = true
	m.savedPhotos = append([]domain.PostCreateSavedPhoto(nil), photos...)
	return nil
}

type mockPublishSender struct {
	last domain.PublishEmailMessage
	sent bool
}

func (m *mockPublishSender) SendPublishEmail(_ context.Context, msg domain.PublishEmailMessage) error {
	m.last = msg
	m.sent = true
	return nil
}

type mockPostCreatePhotoUploader struct {
	uploads []domain.PostCreatePhotoUpload
}

func (m *mockPostCreatePhotoUploader) UploadPostPhoto(_ context.Context, postID int64, photo domain.PostCreatePhotoUpload) (domain.PostCreateSavedPhoto, error) {
	m.uploads = append(m.uploads, photo)
	return domain.PostCreateSavedPhoto{
		PostID:      postID,
		S3Key:       "v2/posts/130031999/photo.jpg",
		TickerS3Key: "",
		Position:    photo.Position,
	}, nil
}

func TestPostCreateService_Submit_DryRunDoesNotPersistOrSend(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)
	sender := &mockPublishSender{}

	result, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         100,
		PriceProvided: true,
	}, true, "https://supost.com", "response@mg.supost.com", sender, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected dry run result")
	}
	if repo.createCalled {
		t.Fatalf("expected no repository insert in dry run")
	}
	if sender.sent {
		t.Fatalf("expected no email send in dry run")
	}
	if !strings.Contains(result.PublishURL, "/post/publish/") {
		t.Fatalf("missing publish URL in result: %q", result.PublishURL)
	}
}

func TestPostCreateService_Submit_PersistsAndSends(t *testing.T) {
	now := time.Date(2026, time.February, 26, 17, 32, 0, 0, time.UTC)
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
		persisted: domain.PostCreatePersisted{
			PostID:      130031999,
			AccessToken: "abcdef",
			PostedAt:    now,
		},
	}
	sender := &mockPublishSender{}
	svc := NewPostCreateService(repo)

	result, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		IP:            "203.0.113.10",
		Price:         100,
		PriceProvided: true,
	}, false, "https://supost.com", "response@mg.supost.com", sender, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !repo.createCalled {
		t.Fatalf("expected repository insert to be called")
	}
	if !sender.sent {
		t.Fatalf("expected publish email send")
	}
	if repo.submission.IP != "203.0.113.10" {
		t.Fatalf("expected IP to be forwarded to repository, got %q", repo.submission.IP)
	}
	if result.PostID != 130031999 {
		t.Fatalf("unexpected post id %d", result.PostID)
	}
	if sender.last.To != "wientjes@alumni.stanford.edu" {
		t.Fatalf("unexpected recipient %q", sender.last.To)
	}
	if !strings.Contains(sender.last.Text, "https://supost.com/post/publish/abcdef") {
		t.Fatalf("missing publish URL in email body")
	}
}

func TestPostCreateService_Submit_InvalidEmailRejected(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "user@gmail.com",
		Price:         100,
		PriceProvided: true,
	}, false, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "Stanford email") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "1 error prohibited this post from being saved") {
		t.Fatalf("expected formatted validation header, got %v", err)
	}
}

func TestPostCreateService_Submit_InvalidIPRejected(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		IP:            "not-an-ip",
		Price:         100,
		PriceProvided: true,
	}, false, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "IP must be a valid IPv4 or IPv6 address.") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_StanfordSubdomainAccepted(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@cs.stanford.edu",
		Price:         100,
		PriceProvided: true,
	}, true, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err != nil {
		t.Fatalf("expected cs.stanford.edu to pass validation, got %v", err)
	}
}

func TestPostCreateService_Submit_PriceRequiredForForSale(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
	}, true, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "Price is required for this category.") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_PriceForbiddenForPersonals(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 8, Name: "personals/dating", ShortName: "personals"}},
		subcategories: []domain.Subcategory{
			{ID: 130, CategoryID: 8, Name: "friendship"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    8,
		SubcategoryID: 130,
		Name:          "Missed connection",
		Body:          "Saw you at Coupa.",
		Email:         "wientjes@cs.stanford.edu",
		Price:         10,
		PriceProvided: true,
	}, true, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "Price is not allowed for this category.") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_PriceForbiddenBeforeSubcategoryMismatch(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{
			{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
			{ID: 8, Name: "personals/dating", ShortName: "personals"},
		},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
			{ID: 130, CategoryID: 8, Name: "friendship"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    8,
		SubcategoryID: 14,
		Name:          "Yellow Orange bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         105,
		PriceProvided: true,
	}, true, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "Price is not allowed for this category.") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_SubcategoryMismatchWithoutPrice(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{
			{ID: 5, Name: "for sale/wanted", ShortName: "for sale"},
			{ID: 8, Name: "personals/dating", ShortName: "personals"},
		},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
			{ID: 130, CategoryID: 8, Name: "friendship"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    8,
		SubcategoryID: 14,
		Name:          "Yellow Orange bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
	}, true, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err == nil {
		t.Fatalf("expected subcategory mismatch error")
	}
	if !strings.Contains(err.Error(), "subcategory 14 not found in category 8") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_PersonalsNoPriceAllowed(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 8, Name: "personals/dating", ShortName: "personals"}},
		subcategories: []domain.Subcategory{
			{ID: 130, CategoryID: 8, Name: "friendship"},
		},
		persisted: domain.PostCreatePersisted{
			PostID:      130032100,
			AccessToken: "tok_personals",
		},
	}
	svc := NewPostCreateService(repo)
	sender := &mockPublishSender{}

	result, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    8,
		SubcategoryID: 130,
		Name:          "Missed connection",
		Body:          "Saw you at Coupa.",
		Email:         "wientjes@cs.stanford.edu",
	}, false, "https://supost.com", "response@mg.supost.com", sender, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.EmailSent || !repo.createCalled {
		t.Fatalf("expected personals submit to send+persist")
	}
	if repo.submission.PriceProvided {
		t.Fatalf("expected personals submit with no price")
	}
}

func TestPostCreateService_Submit_WithPhotosUploadsAndPersistsMetadata(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
		persisted: domain.PostCreatePersisted{
			PostID:      130031999,
			AccessToken: "abcdef",
			PostedAt:    time.Now(),
		},
	}
	sender := &mockPublishSender{}
	uploader := &mockPostCreatePhotoUploader{}
	svc := NewPostCreateService(repo)

	result, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         100,
		PriceProvided: true,
		Photos: []domain.PostCreatePhotoUpload{
			{FileName: "photo-1.jpg", ContentType: "image/jpeg", Content: []byte("image-1")},
			{FileName: "photo-2.png", ContentType: "image/png", Content: []byte("image-2")},
		},
	}, false, "https://supost.com", "response@mg.supost.com", sender, uploader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(uploader.uploads) != 2 {
		t.Fatalf("expected 2 uploads, got %d", len(uploader.uploads))
	}
	if !repo.saveCalled {
		t.Fatalf("expected photo metadata to be persisted")
	}
	if len(repo.savedPhotos) != 2 {
		t.Fatalf("expected 2 saved photo rows, got %d", len(repo.savedPhotos))
	}
	if result.PhotoCount != 2 {
		t.Fatalf("expected photo_count=2, got %d", result.PhotoCount)
	}
}

func TestPostCreateService_Submit_TooManyPhotosRejected(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         100,
		PriceProvided: true,
		Photos: []domain.PostCreatePhotoUpload{
			{FileName: "1.jpg", Content: []byte("1")},
			{FileName: "2.jpg", Content: []byte("2")},
			{FileName: "3.jpg", Content: []byte("3")},
			{FileName: "4.jpg", Content: []byte("4")},
			{FileName: "5.jpg", Content: []byte("5")},
		},
	}, true, "https://supost.com", "response@mg.supost.com", &mockPublishSender{}, nil)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "At most 4 photos are allowed.") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_RejectsMissingPhotoUploaderWhenPhotosProvided(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
		persisted: domain.PostCreatePersisted{
			PostID:      130031999,
			AccessToken: "abcdef",
			PostedAt:    time.Now(),
		},
	}
	sender := &mockPublishSender{}
	svc := NewPostCreateService(repo)

	_, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         100,
		PriceProvided: true,
		Photos: []domain.PostCreatePhotoUpload{
			{FileName: "photo-1.jpg", ContentType: "image/jpeg", Content: []byte("image-1")},
		},
	}, false, "https://supost.com", "response@mg.supost.com", sender, nil)
	if err == nil {
		t.Fatalf("expected uploader required error")
	}
	if !strings.Contains(err.Error(), "photo uploader is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostCreateService_Submit_DryRunAllowsPhotosWithoutUploader(t *testing.T) {
	repo := &mockPostCreateSubmitRepo{
		categories: []domain.Category{{ID: 5, Name: "for sale/wanted", ShortName: "for sale"}},
		subcategories: []domain.Subcategory{
			{ID: 14, CategoryID: 5, Name: "furniture"},
		},
	}
	sender := &mockPublishSender{}
	svc := NewPostCreateService(repo)

	result, err := svc.Submit(context.Background(), domain.PostCreateSubmission{
		CategoryID:    5,
		SubcategoryID: 14,
		Name:          "Red bike for sale",
		Body:          "Pick up on campus.",
		Email:         "wientjes@alumni.stanford.edu",
		Price:         100,
		PriceProvided: true,
		Photos: []domain.PostCreatePhotoUpload{
			{FileName: "photo-1.jpg", ContentType: "image/jpeg", Content: []byte("image-1")},
		},
	}, true, "https://supost.com", "response@mg.supost.com", sender, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.PhotoCount != 1 {
		t.Fatalf("expected photo_count=1, got %d", result.PhotoCount)
	}
}
