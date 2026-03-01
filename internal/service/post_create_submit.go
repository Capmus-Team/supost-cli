package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/Capmus-Team/supost-cli/internal/domain"
)

const (
	defaultSupostBaseURL = "https://supost.com"
	publishSafetyURL     = "https://supost.com/safety"
)

var exactStanfordEmailDomains = map[string]struct{}{
	"stanford.edu":           {},
	"alumni.stanford.edu":    {},
	"stanfordalumni.org":     {},
	"stanfordchildrens.org":  {},
	"stanfordhealthcare.org": {},
	"stanfordmed.org":        {},
	"lpch.org":               {},
}

// PostCreateEmailSender defines publish-email side effects where consumed.
type PostCreateEmailSender interface {
	SendPublishEmail(ctx context.Context, msg domain.PublishEmailMessage) error
}

// PostCreatePhotoUploader defines photo upload side effects where consumed.
type PostCreatePhotoUploader interface {
	UploadPostPhoto(ctx context.Context, postID int64, photo domain.PostCreatePhotoUpload) (domain.PostCreateSavedPhoto, error)
}

// Submit creates a post, persists it, and sends a publish-link email.
func (s *PostCreateService) Submit(
	ctx context.Context,
	input domain.PostCreateSubmission,
	dryRun bool,
	baseURL string,
	fromEmail string,
	sender PostCreateEmailSender,
	photoUploader PostCreatePhotoUploader,
) (domain.PostCreateSubmitResult, error) {
	normalized, err := s.validateSubmissionInput(ctx, input)
	if err != nil {
		return domain.PostCreateSubmitResult{}, err
	}

	if normalized.AccessToken == "" {
		token, err := generateAccessTokenHex(32)
		if err != nil {
			return domain.PostCreateSubmitResult{}, fmt.Errorf("generating access token: %w", err)
		}
		normalized.AccessToken = token
	}
	if normalized.PostedAt.IsZero() {
		normalized.PostedAt = time.Now()
	}

	publishURL := buildPublishURL(baseURL, normalized.AccessToken)
	subject, body := buildPublishEmailContent(normalized.Name, publishURL, normalized.PostedAt)

	result := domain.PostCreateSubmitResult{
		DryRun:      dryRun,
		AccessToken: normalized.AccessToken,
		PublishURL:  publishURL,
		PostedAt:    normalized.PostedAt,
		EmailTo:     normalized.Email,
		EmailSent:   false,
		PhotoCount:  len(normalized.Photos),
		Subject:     subject,
		Body:        body,
	}

	if dryRun {
		return result, nil
	}
	if sender == nil {
		return domain.PostCreateSubmitResult{}, fmt.Errorf("email sender is required")
	}

	persisted, err := s.repo.CreatePendingPost(ctx, normalized)
	if err != nil {
		return domain.PostCreateSubmitResult{}, err
	}

	if persisted.PostID > 0 {
		result.PostID = persisted.PostID
	}
	if strings.TrimSpace(persisted.AccessToken) != "" {
		result.AccessToken = strings.TrimSpace(persisted.AccessToken)
		result.PublishURL = buildPublishURL(baseURL, result.AccessToken)
		result.Subject, result.Body = buildPublishEmailContent(normalized.Name, result.PublishURL, result.PostedAt)
	}
	if !persisted.PostedAt.IsZero() {
		result.PostedAt = persisted.PostedAt
		result.Subject, result.Body = buildPublishEmailContent(normalized.Name, result.PublishURL, result.PostedAt)
	}

	if len(normalized.Photos) > 0 {
		if photoUploader == nil {
			return domain.PostCreateSubmitResult{}, fmt.Errorf("photo uploader is required")
		}
		if result.PostID <= 0 {
			return domain.PostCreateSubmitResult{}, fmt.Errorf("invalid persisted post id for photo upload")
		}

		savedPhotos := make([]domain.PostCreateSavedPhoto, 0, len(normalized.Photos))
		for _, photo := range normalized.Photos {
			savedPhoto, err := photoUploader.UploadPostPhoto(ctx, result.PostID, photo)
			if err != nil {
				return domain.PostCreateSubmitResult{}, fmt.Errorf("uploading photo at position %d: %w", photo.Position, err)
			}
			savedPhoto.PostID = result.PostID
			savedPhoto.Position = photo.Position
			savedPhotos = append(savedPhotos, savedPhoto)
		}

		if err := s.repo.SavePostPhotos(ctx, savedPhotos); err != nil {
			return domain.PostCreateSubmitResult{}, fmt.Errorf("saving post photos: %w", err)
		}

		result.PhotoS3Keys = make([]string, 0, len(savedPhotos))
		for _, savedPhoto := range savedPhotos {
			result.PhotoS3Keys = append(result.PhotoS3Keys, savedPhoto.S3Key)
		}
	}

	msg := domain.PublishEmailMessage{
		From:    strings.TrimSpace(fromEmail),
		To:      result.EmailTo,
		Subject: result.Subject,
		Text:    result.Body,
	}
	if err := sender.SendPublishEmail(ctx, msg); err != nil {
		return domain.PostCreateSubmitResult{}, err
	}
	result.EmailSent = true
	return result, nil
}

func (s *PostCreateService) validateSubmissionInput(ctx context.Context, input domain.PostCreateSubmission) (domain.PostCreateSubmission, error) {
	normalized := input
	normalized.Name = strings.TrimSpace(input.Name)
	normalized.Body = strings.TrimSpace(input.Body)
	normalized.Email = strings.ToLower(strings.TrimSpace(input.Email))
	normalized.IP = strings.TrimSpace(input.IP)
	normalized.Photos = normalizePostCreatePhotos(input.Photos)

	problems := make([]string, 0, 12)
	if normalized.CategoryID <= 0 {
		problems = append(problems, "category is required")
	}
	if normalized.SubcategoryID <= 0 {
		problems = append(problems, "subcategory is required")
	}
	if normalized.Name == "" {
		problems = append(problems, "name is required")
	}
	if normalized.Body == "" {
		problems = append(problems, "body is required")
	}
	if normalized.Email == "" {
		problems = append(problems, "Email is required.")
	} else if !isStanfordEmail(normalized.Email) {
		problems = append(problems, "Email must be a Stanford email (e.g., @stanford.edu, @cs.stanford.edu).")
	}
	if normalized.IP != "" {
		if _, err := netip.ParseAddr(normalized.IP); err != nil {
			problems = append(problems, "IP must be a valid IPv4 or IPv6 address.")
		}
	}
	if domain.CategoryPriceRequired(normalized.CategoryID) {
		if !normalized.PriceProvided {
			problems = append(problems, "Price is required for this category.")
		} else if normalized.Price < 0 {
			problems = append(problems, "Price must be non-negative.")
		}
	} else if normalized.PriceProvided {
		problems = append(problems, "Price is not allowed for this category.")
	}
	if len(normalized.Photos) > 4 {
		problems = append(problems, "At most 4 photos are allowed.")
	}
	for _, photo := range normalized.Photos {
		if len(photo.Content) == 0 {
			problems = append(problems, fmt.Sprintf("Photo at position %d is empty.", photo.Position))
		}
	}

	if len(problems) > 0 {
		return domain.PostCreateSubmission{}, fmt.Errorf("%s", formatPostCreateValidationErrors(problems))
	}

	page, err := s.BuildPage(ctx, normalized.CategoryID, normalized.SubcategoryID)
	if err != nil {
		return domain.PostCreateSubmission{}, err
	}
	if page.Stage != domain.PostCreateStageForm {
		return domain.PostCreateSubmission{}, fmt.Errorf("invalid category/subcategory combination")
	}
	return normalized, nil
}

func normalizePostCreatePhotos(photos []domain.PostCreatePhotoUpload) []domain.PostCreatePhotoUpload {
	if len(photos) == 0 {
		return nil
	}

	normalized := make([]domain.PostCreatePhotoUpload, 0, len(photos))
	for idx, photo := range photos {
		normalized = append(normalized, domain.PostCreatePhotoUpload{
			FileName:    strings.TrimSpace(photo.FileName),
			ContentType: strings.TrimSpace(photo.ContentType),
			Content:     append([]byte(nil), photo.Content...),
			Position:    idx,
		})
	}
	return normalized
}

func formatPostCreateValidationErrors(problems []string) string {
	count := len(problems)
	header := fmt.Sprintf("%d errors prohibited this post from being saved", count)
	if count == 1 {
		header = "1 error prohibited this post from being saved"
	}
	return header + "\nThere were problems with the following fields:\n\n" + strings.Join(problems, "\n")
}

func isStanfordEmail(email string) bool {
	at := strings.LastIndex(email, "@")
	if at <= 0 || at == len(email)-1 {
		return false
	}
	domainPart := strings.ToLower(strings.TrimSpace(email[at+1:]))
	if domainPart == "" {
		return false
	}
	if _, ok := exactStanfordEmailDomains[domainPart]; ok {
		return true
	}
	return strings.HasSuffix(domainPart, ".stanford.edu")
}

func generateAccessTokenHex(numBytes int) (string, error) {
	if numBytes <= 0 {
		return "", fmt.Errorf("numBytes must be positive")
	}
	buf := make([]byte, numBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func buildPublishURL(baseURL, accessToken string) string {
	root := strings.TrimSpace(baseURL)
	if root == "" {
		root = defaultSupostBaseURL
	}
	root = strings.TrimRight(root, "/")
	return root + "/post/publish/" + strings.TrimSpace(accessToken)
}

func buildPublishEmailContent(postName, publishURL string, postedAt time.Time) (string, string) {
	title := strings.TrimSpace(postName)
	if title == "" {
		title = "(untitled post)"
	}
	if postedAt.IsZero() {
		postedAt = time.Now()
	}
	dateLine := postedAt.Format("Mon, Jan 2, 2006 03:04 PM")
	subject := "SUpost - Publish your post! " + title
	body := strings.Join([]string{
		"Publish your post by pressing:",
		"",
		publishURL,
		"",
		title,
		"",
		"Posted on: " + dateLine + " -- Stanford University",
		"",
		"Do not send electronic payments to sellers: " + publishSafetyURL,
	}, "\n")
	return subject, body
}
