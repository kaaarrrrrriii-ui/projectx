package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
)

var (
	ErrTicketNotFound         = repos.ErrTicketNotFound
	ErrAttachmentNotFound     = repos.ErrAttachmentNotFound
	ErrInvalidTransition      = repos.ErrInvalidTransition
	ErrReturnLimitReached     = repos.ErrReturnLimitReached
	ErrAttachmentLimitReached = repos.ErrAttachmentLimitReached
	ErrInvalidTrackID         = errors.New("invalid track id")
	ErrMessageTextRequired    = errors.New("message text is required")
)

const maxTicketReturns = 2

var trackIDPattern = regexp.MustCompile(`^ОТК-[23456789ABCDEFGHJKMNPQRSTUVWXYZ]{4}-[23456789ABCDEFGHJKMNPQRSTUVWXYZ]{4}$`)

type ticketRepository interface {
	GetStatusByTrack(context.Context, string) (repos.TicketStatusRecord, error)
	GetChatByTrack(context.Context, string) (repos.TicketChatRecord, error)
	Complete(context.Context, string, time.Time) (repos.CompleteTicketRecord, error)
	Return(context.Context, string, string, time.Time) (repos.ReturnTicketRecord, error)
}

type TicketService struct {
	repository ticketRepository
	now        func() time.Time
}

type CategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TicketStatusResponse struct {
	TrackID     string           `json:"track_id"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	Category    CategoryResponse `json:"category"`
	CanOpenChat bool             `json:"can_open_chat"`
	CanComplete bool             `json:"can_complete"`
	CanReturn   bool             `json:"can_return"`
}

type SpecialistResponse struct {
	Label       string `json:"label"`
	ExpertGroup string `json:"expert_group"`
}

type MessageResponse struct {
	ID          int64                `json:"id"`
	Text        string               `json:"text"`
	Type        string               `json:"type"`
	CreatedAt   time.Time            `json:"created_at"`
	Attachments []AttachmentResponse `json:"attachments"`
}

type AttachmentResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	MIMEType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

type TicketChatResponse struct {
	TrackID     string               `json:"track_id"`
	Status      string               `json:"status"`
	Specialist  *SpecialistResponse  `json:"specialist"`
	Messages    []MessageResponse    `json:"messages"`
	Attachments []AttachmentResponse `json:"attachments"`
}

type CompleteTicketResponse struct {
	Status   string    `json:"status"`
	ClosedAt time.Time `json:"closed_at"`
}

type ReturnTicketResponse struct {
	Status      string `json:"status"`
	ReturnCount int    `json:"return_count"`
}

func NewTicketService(repository ticketRepository) (*TicketService, error) {
	if repository == nil {
		return nil, errors.New("ticket repository is required")
	}
	return &TicketService{repository: repository, now: time.Now}, nil
}

func (service *TicketService) Status(ctx context.Context, trackID string) (TicketStatusResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return TicketStatusResponse{}, err
	}
	record, err := service.repository.GetStatusByTrack(ctx, normalized)
	if err != nil {
		return TicketStatusResponse{}, err
	}
	return TicketStatusResponse{
		TrackID:     record.TrackID,
		Status:      record.Status.String(),
		CreatedAt:   record.CreatedAt.UTC(),
		Category:    CategoryResponse{ID: record.CategoryID, Name: record.Category},
		CanOpenChat: canOpenChat(record.Status),
		CanComplete: record.Status.IsApplicantCompletable(),
		CanReturn:   record.Status == models.TicketStatusAnswerReady && record.ReturnCount < maxTicketReturns,
	}, nil
}

func (service *TicketService) Chat(ctx context.Context, trackID string) (TicketChatResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return TicketChatResponse{}, err
	}
	record, err := service.repository.GetChatByTrack(ctx, normalized)
	if err != nil {
		return TicketChatResponse{}, err
	}

	response := TicketChatResponse{
		TrackID:     record.TrackID,
		Status:      record.Status.String(),
		Messages:    make([]MessageResponse, 0, len(record.Messages)),
		Attachments: make([]AttachmentResponse, 0, len(record.Attachments)),
	}
	if record.Specialist != nil {
		response.Specialist = &SpecialistResponse{
			Label:       record.Specialist.Label,
			ExpertGroup: record.Specialist.ExpertGroup,
		}
	}
	for _, message := range record.Messages {
		response.Messages = append(response.Messages, MessageResponse{
			ID:          message.ID,
			Text:        message.Text,
			Type:        message.Type.String(),
			CreatedAt:   message.CreatedAt.UTC(),
			Attachments: make([]AttachmentResponse, 0),
		})
	}
	for _, attachment := range record.Attachments {
		response.Attachments = append(response.Attachments, AttachmentResponse{
			ID:        attachment.ID,
			Name:      attachment.SafeName,
			MIMEType:  attachment.MIMEType,
			SizeBytes: attachment.SizeBytes,
		})
	}
	return response, nil
}

func (service *TicketService) Complete(ctx context.Context, trackID string) (CompleteTicketResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return CompleteTicketResponse{}, err
	}
	record, err := service.repository.Complete(ctx, normalized, service.now().UTC())
	if err != nil {
		return CompleteTicketResponse{}, err
	}
	return CompleteTicketResponse{Status: record.Status.String(), ClosedAt: record.ClosedAt.UTC()}, nil
}

func (service *TicketService) Return(ctx context.Context, trackID, reason string) (ReturnTicketResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return ReturnTicketResponse{}, err
	}
	record, err := service.repository.Return(ctx, normalized, strings.TrimSpace(reason), service.now().UTC())
	if err != nil {
		return ReturnTicketResponse{}, err
	}
	return ReturnTicketResponse{Status: record.Status.String(), ReturnCount: record.ReturnCount}, nil
}

func normalizeTrackID(trackID string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(trackID))
	if !trackIDPattern.MatchString(normalized) {
		return "", ErrInvalidTrackID
	}
	return normalized, nil
}

func canOpenChat(status models.TicketStatus) bool {
	switch status {
	case models.TicketStatusAssigned,
		models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification,
		models.TicketStatusAnswerReady,
		models.TicketStatusReturned,
		models.TicketStatusCompleted:
		return true
	default:
		return false
	}
}
