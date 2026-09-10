package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
)

type stubTicketRepository struct {
	statusRecord   repos.TicketStatusRecord
	statusErr      error
	statusTrackID  string
	chatRecord     repos.TicketChatRecord
	chatErr        error
	completeRecord repos.CompleteTicketRecord
	completeErr    error
	completeAt     time.Time
	returnRecord   repos.ReturnTicketRecord
	returnErr      error
	returnReason   string
}

func (repository *stubTicketRepository) GetStatusByTrack(_ context.Context, trackID string) (repos.TicketStatusRecord, error) {
	repository.statusTrackID = trackID
	return repository.statusRecord, repository.statusErr
}

func (repository *stubTicketRepository) GetChatByTrack(context.Context, string) (repos.TicketChatRecord, error) {
	return repository.chatRecord, repository.chatErr
}

func (repository *stubTicketRepository) Complete(_ context.Context, _ string, closedAt time.Time) (repos.CompleteTicketRecord, error) {
	repository.completeAt = closedAt
	return repository.completeRecord, repository.completeErr
}

func (repository *stubTicketRepository) Return(_ context.Context, _, reason string, _ time.Time) (repos.ReturnTicketRecord, error) {
	repository.returnReason = reason
	return repository.returnRecord, repository.returnErr
}

func TestTicketServiceStatusNormalizesTrackAndExposesActions(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 9, 9, 14, 48, 0, 0, time.FixedZone("test", 5*60*60))
	repository := &stubTicketRepository{statusRecord: repos.TicketStatusRecord{
		TrackID:     "ОТК-ABCD-2345",
		Status:      models.TicketStatusAnswerReady,
		CreatedAt:   createdAt,
		CategoryID:  2,
		Category:    "Кибербуллинг",
		ReturnCount: 1,
	}}
	ticketService, err := NewTicketService(repository)
	if err != nil {
		t.Fatalf("NewTicketService() error = %v", err)
	}

	response, err := ticketService.Status(context.Background(), "  отк-abcd-2345 ")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if repository.statusTrackID != "ОТК-ABCD-2345" {
		t.Fatalf("repository track id = %q", repository.statusTrackID)
	}
	if response.Status != "answer_ready" || !response.CanOpenChat || !response.CanComplete || !response.CanReturn {
		t.Fatalf("Status() response = %+v", response)
	}
	if response.CreatedAt.Location() != time.UTC {
		t.Fatalf("created_at location = %v, want UTC", response.CreatedAt.Location())
	}
}

func TestTicketServiceStatusHidesReturnAtLimit(t *testing.T) {
	t.Parallel()

	repository := &stubTicketRepository{statusRecord: repos.TicketStatusRecord{
		TrackID:     "ОТК-ABCD-2345",
		Status:      models.TicketStatusAnswerReady,
		ReturnCount: 2,
	}}
	ticketService, _ := NewTicketService(repository)
	response, err := ticketService.Status(context.Background(), "ОТК-ABCD-2345")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if response.CanReturn {
		t.Fatal("CanReturn = true at the return limit")
	}
}

func TestTicketServiceRejectsMalformedTrackBeforeRepository(t *testing.T) {
	t.Parallel()

	repository := &stubTicketRepository{}
	ticketService, _ := NewTicketService(repository)
	_, err := ticketService.Status(context.Background(), "not-a-track")
	if !errors.Is(err, ErrInvalidTrackID) {
		t.Fatalf("Status() error = %v, want %v", err, ErrInvalidTrackID)
	}
	if repository.statusTrackID != "" {
		t.Fatal("repository called for malformed track id")
	}
}

func TestTicketServiceReturnTrimsOptionalReason(t *testing.T) {
	t.Parallel()

	repository := &stubTicketRepository{returnRecord: repos.ReturnTicketRecord{
		Status:      models.TicketStatusReturned,
		ReturnCount: 1,
	}}
	ticketService, _ := NewTicketService(repository)
	response, err := ticketService.Return(context.Background(), "ОТК-ABCD-2345", "  Не хватило деталей  ")
	if err != nil {
		t.Fatalf("Return() error = %v", err)
	}
	if repository.returnReason != "Не хватило деталей" {
		t.Fatalf("repository reason = %q", repository.returnReason)
	}
	if response.Status != "returned" || response.ReturnCount != 1 {
		t.Fatalf("Return() response = %+v", response)
	}
}

func TestTicketServiceLimitsReturnReasonLength(t *testing.T) {
	t.Parallel()
	repository := &stubTicketRepository{}
	ticketService, _ := NewTicketService(repository)

	_, err := ticketService.Return(context.Background(), "ОТК-ABCD-2345", strings.Repeat("я", MaxReturnReasonCharacters+1))
	if !errors.Is(err, ErrTextTooLong) {
		t.Fatalf("Return() error = %v, want %v", err, ErrTextTooLong)
	}
	if repository.returnReason != "" {
		t.Fatal("repository called for an oversized return reason")
	}
}

func TestTicketStatusApplicantCompletionRules(t *testing.T) {
	t.Parallel()

	if !models.TicketStatusAnswerReady.IsApplicantCompletable() {
		t.Fatal("answer_ready must be applicant-completable")
	}
	for _, status := range []models.TicketStatus{
		models.TicketStatusNew,
		models.TicketStatusAssigned,
		models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification,
		models.TicketStatusReturned,
		models.TicketStatusCompleted,
		models.TicketStatusRejected,
		models.TicketStatusClosedWithoutAnswer,
	} {
		if status.IsApplicantCompletable() {
			t.Fatalf("status %s must not be applicant-completable", status.String())
		}
	}
}
