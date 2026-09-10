package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
	"example.com/german/backend/internal/storage"
)

const (
	MaxDescriptionCharacters      = 20_000
	MaxCustomTopicCharacters      = 120
	MaxCrisisContactCharacters    = 500
	MaxReviewTextCharacters       = 2_000
	MaxReturnReasonCharacters     = 2_000
	MaxApplicantMessageCharacters = 20_000
	maxTrackCodeAttempts          = 5
)

var (
	ErrCategoryNotFound    = repos.ErrCategoryNotFound
	ErrInvalidAnswer       = repos.ErrInvalidAnswer
	ErrReviewAlreadyExists = repos.ErrReviewAlreadyExists
	ErrReviewNotAllowed    = repos.ErrReviewNotAllowed
	ErrInvalidApplicant    = errors.New("invalid applicant type")
	ErrInvalidCategory     = errors.New("invalid category id")
	ErrInvalidSubmission   = errors.New("invalid ticket submission")
	ErrInvalidRating       = errors.New("rating must be between 1 and 5")
	ErrTextTooLong         = errors.New("text exceeds the allowed length")
)

type submissionRepository interface {
	ListPublicCategories(context.Context) ([]repos.PublicCategoryRecord, error)
	ListPublicQuestions(context.Context, int64) ([]repos.PublicQuestionRecord, error)
	SelectedAnswerTexts(context.Context, int64, []repos.TicketAnswerRecord) ([]string, error)
	CreatePublicTicket(context.Context, repos.CreateTicketRecord) (repos.CreatedTicketRecord, error)
	CreateReview(context.Context, string, int, string) (repos.CreatedReviewRecord, error)
}

type TicketSubmissionService struct {
	repository submissionRepository
	preparer   *AttachmentPreparer
	storage    storage.Storage
	limits     AttachmentLimits
	now        func() time.Time
	trackCode  func() (string, error)
}

type CategoriesResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

type PublicAnswerResponse struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}

type PublicQuestionResponse struct {
	ID      int64                  `json:"id"`
	Text    string                 `json:"text"`
	Answers []PublicAnswerResponse `json:"answers"`
}

type QuestionsResponse struct {
	Questions []PublicQuestionResponse `json:"questions"`
}

type TicketAnswerInput struct {
	QuestionID int64
	AnswerID   int64
}

type CreateTicketInput struct {
	ApplicantType string
	CategoryID    int64
	Description   string
	CustomTopic   string
	Answers       []TicketAnswerInput
	Attachments   []UploadedAttachment
	CrisisContact string
}

type CrisisContactResponse struct {
	Title       string `json:"title"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
	URL         string `json:"url,omitempty"`
}

type CreateTicketResponse struct {
	TrackID        string                  `json:"track_id"`
	Status         string                  `json:"status"`
	CreatedAt      time.Time               `json:"created_at"`
	CrisisDetected bool                    `json:"crisis_detected"`
	CrisisContacts []CrisisContactResponse `json:"crisis_contacts"`
}

type CreateReviewResponse struct {
	ID     int64  `json:"id"`
	Rating int    `json:"rating"`
	Text   string `json:"text"`
}

func NewTicketSubmissionService(repository submissionRepository, preparer *AttachmentPreparer, attachmentStorage storage.Storage) (*TicketSubmissionService, error) {
	if repository == nil || preparer == nil || attachmentStorage == nil {
		return nil, errors.New("submission repository, attachment preparer and storage are required")
	}
	return &TicketSubmissionService{
		repository: repository,
		preparer:   preparer,
		storage:    attachmentStorage,
		limits:     DefaultAttachmentLimits(),
		now:        time.Now,
		trackCode:  GenerateTrackCode,
	}, nil
}

func (service *TicketSubmissionService) Categories(ctx context.Context) (CategoriesResponse, error) {
	records, err := service.repository.ListPublicCategories(ctx)
	if err != nil {
		return CategoriesResponse{}, err
	}
	response := CategoriesResponse{Categories: make([]CategoryResponse, 0, len(records))}
	for _, record := range records {
		response.Categories = append(response.Categories, CategoryResponse{ID: record.ID, Name: record.Name})
	}
	return response, nil
}

func (service *TicketSubmissionService) Questions(ctx context.Context, categoryID int64) (QuestionsResponse, error) {
	if categoryID <= 0 {
		return QuestionsResponse{}, ErrInvalidCategory
	}
	records, err := service.repository.ListPublicQuestions(ctx, categoryID)
	if err != nil {
		return QuestionsResponse{}, err
	}
	response := QuestionsResponse{Questions: make([]PublicQuestionResponse, 0, len(records))}
	for _, record := range records {
		question := PublicQuestionResponse{ID: record.ID, Text: record.Text, Answers: make([]PublicAnswerResponse, 0, len(record.Answers))}
		for _, answer := range record.Answers {
			question.Answers = append(question.Answers, PublicAnswerResponse{ID: answer.ID, Text: answer.Text})
		}
		response.Questions = append(response.Questions, question)
	}
	return response, nil
}

func (service *TicketSubmissionService) Create(ctx context.Context, input CreateTicketInput) (CreateTicketResponse, error) {
	applicantType, err := parseApplicantType(input.ApplicantType)
	if err != nil {
		return CreateTicketResponse{}, err
	}
	if input.CategoryID <= 0 {
		return CreateTicketResponse{}, ErrInvalidCategory
	}

	description := normalizeMultilineText(input.Description)
	customTopic := normalizeSingleLineText(input.CustomTopic)
	crisisContact := normalizeSingleLineText(input.CrisisContact)
	if utf8.RuneCountInString(description) > MaxDescriptionCharacters ||
		utf8.RuneCountInString(customTopic) > MaxCustomTopicCharacters ||
		utf8.RuneCountInString(crisisContact) > MaxCrisisContactCharacters {
		return CreateTicketResponse{}, ErrTextTooLong
	}

	answers := make([]repos.TicketAnswerRecord, 0, len(input.Answers))
	seenQuestions := make(map[int64]struct{}, len(input.Answers))
	for _, answer := range input.Answers {
		if answer.QuestionID <= 0 || answer.AnswerID <= 0 {
			return CreateTicketResponse{}, ErrInvalidAnswer
		}
		if _, exists := seenQuestions[answer.QuestionID]; exists {
			return CreateTicketResponse{}, ErrInvalidAnswer
		}
		seenQuestions[answer.QuestionID] = struct{}{}
		answers = append(answers, repos.TicketAnswerRecord{QuestionID: answer.QuestionID, AnswerID: answer.AnswerID})
	}
	answerTexts, err := service.repository.SelectedAnswerTexts(ctx, input.CategoryID, answers)
	if err != nil {
		return CreateTicketResponse{}, err
	}

	if customTopic != "" {
		if description == "" {
			description = "Тема заявителя: " + customTopic
		} else {
			description = "Тема заявителя: " + customTopic + "\n\n" + description
		}
	}
	crisisText := strings.Join(append([]string{customTopic, description}, answerTexts...), "\n")
	crisis := DetectCrisis(crisisText)
	if !crisis {
		crisisContact = ""
	}

	candidates := make([]AttachmentCandidate, 0, len(input.Attachments))
	for _, attachment := range input.Attachments {
		candidates = append(candidates, AttachmentCandidate{OriginalSizeBytes: attachment.Size})
	}
	if err := ValidateAttachmentBatch(service.limits, nil, candidates); err != nil {
		return CreateTicketResponse{}, err
	}

	prepared := make([]PreparedAttachment, 0, len(input.Attachments))
	cleanup := func() {
		for _, attachment := range prepared {
			_ = service.storage.Delete(context.Background(), attachment.StorageKey)
		}
	}
	for _, upload := range input.Attachments {
		if upload.Source == nil {
			cleanup()
			return CreateTicketResponse{}, ErrInvalidSubmission
		}
		attachment, err := service.preparer.Prepare(ctx, upload.Source)
		if err != nil {
			cleanup()
			return CreateTicketResponse{}, err
		}
		prepared = append(prepared, attachment)
	}

	attachments := make([]repos.NewAttachmentRecord, 0, len(prepared))
	for _, attachment := range prepared {
		attachments = append(attachments, repos.NewAttachmentRecord{
			StorageKey: attachment.StorageKey, SafeName: attachment.SafeName, MIMEType: attachment.MIMEType,
			OriginalSizeBytes: attachment.OriginalSizeBytes, SizeBytes: attachment.SizeBytes,
			Width: attachment.Width, Height: attachment.Height,
		})
	}

	priority := models.TicketPriorityStandard
	if crisis {
		priority = models.TicketPriorityUrgent
	}
	createdAt := service.now().UTC()
	var created repos.CreatedTicketRecord
	for attempt := 0; attempt < maxTrackCodeAttempts; attempt++ {
		trackID, err := service.trackCode()
		if err != nil {
			cleanup()
			return CreateTicketResponse{}, fmt.Errorf("generate ticket track id: %w", err)
		}
		created, err = service.repository.CreatePublicTicket(ctx, repos.CreateTicketRecord{
			TrackID: trackID, ApplicantType: applicantType, CategoryID: input.CategoryID,
			Priority: priority, Status: models.TicketStatusNew, Description: description,
			Answers: answers, Attachments: attachments, CrisisContact: crisisContact,
			Crisis: crisis, CreatedAt: createdAt,
		})
		if err == nil {
			break
		}
		if !errors.Is(err, repos.ErrTrackIDExists) {
			cleanup()
			return CreateTicketResponse{}, err
		}
	}
	if created.TrackID == "" {
		cleanup()
		return CreateTicketResponse{}, errors.New("could not generate a unique ticket track id")
	}

	response := CreateTicketResponse{
		TrackID: created.TrackID, Status: created.Status.String(), CreatedAt: created.CreatedAt.UTC(),
		CrisisDetected: crisis, CrisisContacts: []CrisisContactResponse{},
	}
	if crisis {
		response.CrisisContacts = russianCrisisContacts()
	}
	return response, nil
}

func (service *TicketSubmissionService) Review(ctx context.Context, trackID string, rating int, text string) (CreateReviewResponse, error) {
	normalizedTrackID, err := normalizeTrackID(trackID)
	if err != nil {
		return CreateReviewResponse{}, err
	}
	text = normalizeMultilineText(text)
	if rating < 1 || rating > 5 {
		return CreateReviewResponse{}, ErrInvalidRating
	}
	if utf8.RuneCountInString(text) > MaxReviewTextCharacters {
		return CreateReviewResponse{}, ErrTextTooLong
	}
	record, err := service.repository.CreateReview(ctx, normalizedTrackID, rating, text)
	if err != nil {
		return CreateReviewResponse{}, err
	}
	return CreateReviewResponse{ID: record.ID, Rating: record.Rating, Text: record.Text}, nil
}

func parseApplicantType(value string) (models.ApplicantType, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "student", "schoolchild":
		return models.ApplicantTypeSchoolchild, nil
	case "parent":
		return models.ApplicantTypeParent, nil
	case "teacher":
		return models.ApplicantTypeTeacher, nil
	default:
		return 0, ErrInvalidApplicant
	}
}

func normalizeSingleLineText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func normalizeMultilineText(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	return strings.TrimSpace(value)
}

func russianCrisisContacts() []CrisisContactResponse {
	return []CrisisContactResponse{
		{Title: "Единый номер экстренных служб", Phone: "112", Description: "Если опасность непосредственная — позвоните прямо сейчас."},
		{Title: "Всероссийский детский телефон доверия", Phone: "124", Description: "Круглосуточно, бесплатно и анонимно с мобильного телефона.", URL: "https://telefon-doveria.ru/"},
		{Title: "Всероссийский детский телефон доверия", Phone: "8 800 2000-122", Description: "Для детей, подростков и родителей; круглосуточно, бесплатно и анонимно.", URL: "https://telefon-doveria.ru/"},
		{Title: "Экстренная психологическая помощь МЧС России", Phone: "+7 (495) 989-50-50", Description: "Круглосуточная психологическая помощь; междугородняя связь оплачивается по тарифу оператора.", URL: "https://psi.mchs.gov.ru/"},
	}
}
