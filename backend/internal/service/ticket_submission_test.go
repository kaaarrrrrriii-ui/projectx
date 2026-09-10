package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
	"example.com/german/backend/internal/storage"
)

type submissionRepositoryStub struct {
	answerTexts []string
	answerErr   error
	createInput repos.CreateTicketRecord
	createErrs  []error
	createCalls int
	review      repos.CreatedReviewRecord
	reviewErr   error
}

func (*submissionRepositoryStub) ListPublicCategories(context.Context) ([]repos.PublicCategoryRecord, error) {
	return []repos.PublicCategoryRecord{}, nil
}

func (*submissionRepositoryStub) ListPublicQuestions(context.Context, int64) ([]repos.PublicQuestionRecord, error) {
	return []repos.PublicQuestionRecord{}, nil
}

func (stub *submissionRepositoryStub) SelectedAnswerTexts(context.Context, int64, []repos.TicketAnswerRecord) ([]string, error) {
	return stub.answerTexts, stub.answerErr
}

func (stub *submissionRepositoryStub) CreatePublicTicket(_ context.Context, input repos.CreateTicketRecord) (repos.CreatedTicketRecord, error) {
	stub.createInput = input
	call := stub.createCalls
	stub.createCalls++
	if call < len(stub.createErrs) && stub.createErrs[call] != nil {
		return repos.CreatedTicketRecord{}, stub.createErrs[call]
	}
	return repos.CreatedTicketRecord{TrackID: input.TrackID, Status: input.Status, CreatedAt: input.CreatedAt}, nil
}

func (stub *submissionRepositoryStub) CreateReview(context.Context, string, int, string) (repos.CreatedReviewRecord, error) {
	return stub.review, stub.reviewErr
}

func newSubmissionServiceForTest(t *testing.T, repository *submissionRepositoryStub) *TicketSubmissionService {
	t.Helper()
	localStorage, err := storage.NewLocal(t.TempDir())
	if err != nil {
		t.Fatalf("NewLocal() error = %v", err)
	}
	processor, err := NewImageProcessor(DefaultImageProcessorConfig())
	if err != nil {
		t.Fatalf("NewImageProcessor() error = %v", err)
	}
	preparer, err := NewAttachmentPreparer(processor, localStorage)
	if err != nil {
		t.Fatalf("NewAttachmentPreparer() error = %v", err)
	}
	service, err := NewTicketSubmissionService(repository, preparer, localStorage)
	if err != nil {
		t.Fatalf("NewTicketSubmissionService() error = %v", err)
	}
	service.now = func() time.Time { return time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC) }
	service.trackCode = func() (string, error) { return "ОТК-ABCD-2345", nil }
	return service
}

func TestCreateTicketDetectsCrisisInSelectedAnswer(t *testing.T) {
	t.Parallel()
	repository := &submissionRepositoryStub{answerTexts: []string{"Мне угрожают убить"}}
	service := newSubmissionServiceForTest(t, repository)

	response, err := service.Create(context.Background(), CreateTicketInput{
		ApplicantType: "student",
		CategoryID:    7,
		Answers:       []TicketAnswerInput{{QuestionID: 3, AnswerID: 9}},
		CrisisContact: " test@example.org ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if !response.CrisisDetected || len(response.CrisisContacts) == 0 {
		t.Fatalf("Create() response = %+v, want crisis contacts", response)
	}
	if repository.createInput.Priority != models.TicketPriorityUrgent {
		t.Fatalf("priority = %v, want urgent", repository.createInput.Priority)
	}
	if repository.createInput.ApplicantType != models.ApplicantTypeSchoolchild {
		t.Fatalf("applicant type = %v, want schoolchild", repository.createInput.ApplicantType)
	}
	if repository.createInput.CrisisContact != "test@example.org" || !repository.createInput.Crisis {
		t.Fatalf("crisis record = %+v", repository.createInput)
	}
}

func TestCreateTicketAllowsEmptyDescriptionAndDropsNonCrisisContact(t *testing.T) {
	t.Parallel()
	repository := &submissionRepositoryStub{}
	service := newSubmissionServiceForTest(t, repository)

	response, err := service.Create(context.Background(), CreateTicketInput{
		ApplicantType: "parent",
		CategoryID:    2,
		CrisisContact: "+7 999 000-00-00",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if response.CrisisDetected || len(response.CrisisContacts) != 0 {
		t.Fatalf("Create() response = %+v, want non-crisis", response)
	}
	if repository.createInput.Description != "" || repository.createInput.CrisisContact != "" {
		t.Fatalf("non-crisis input retained private text: %+v", repository.createInput)
	}
	if repository.createInput.Priority != models.TicketPriorityStandard {
		t.Fatalf("priority = %v, want standard", repository.createInput.Priority)
	}
}

func TestCreateTicketRetriesTrackCollision(t *testing.T) {
	t.Parallel()
	repository := &submissionRepositoryStub{createErrs: []error{repos.ErrTrackIDExists}}
	service := newSubmissionServiceForTest(t, repository)
	codes := []string{"ОТК-ABCD-2345", "ОТК-WXYZ-6789"}
	service.trackCode = func() (string, error) {
		code := codes[0]
		codes = codes[1:]
		return code, nil
	}

	response, err := service.Create(context.Background(), CreateTicketInput{ApplicantType: "teacher", CategoryID: 1})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if repository.createCalls != 2 || response.TrackID != "ОТК-WXYZ-6789" {
		t.Fatalf("calls = %d, response = %+v", repository.createCalls, response)
	}
}

func TestReviewValidationPrecedesRepository(t *testing.T) {
	t.Parallel()
	repository := &submissionRepositoryStub{reviewErr: errors.New("must not be called")}
	service := newSubmissionServiceForTest(t, repository)

	_, err := service.Review(context.Background(), "ОТК-ABCD-2345", 6, "")
	if !errors.Is(err, ErrInvalidRating) {
		t.Fatalf("Review() error = %v, want %v", err, ErrInvalidRating)
	}
}

func TestDetectCrisis(t *testing.T) {
	t.Parallel()
	for _, text := range []string{"Я не хочу жить", "Мне угрожают оружием", "Это сексуальное насилие"} {
		if !DetectCrisis(text) {
			t.Errorf("DetectCrisis(%q) = false", text)
		}
	}
	if DetectCrisis("Мне нужен совет по разговору с учителем") {
		t.Fatal("DetectCrisis() marked an ordinary request as crisis")
	}
}
