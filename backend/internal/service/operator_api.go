package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
)

type operatorAPIRepository interface {
	Dashboard(context.Context, time.Time, time.Time) (repos.OperatorDashboardRecord, error)
	GetOperatorTicket(context.Context, string) (repos.OperatorTicketDetailRecord, error)
	UpdateOperatorTicket(context.Context, string, int64, repos.OperatorTicketUpdate, time.Time) error
	CloseOperatorTicket(context.Context, string, int64, string, time.Time) (repos.OperatorCloseRecord, error)
	CanAccessOperatorAttachment(context.Context, string, int64) error
}

func (service *OperatorService) CanAccessAttachment(ctx context.Context, trackID string, attachmentID int64) error {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return err
	}
	if attachmentID <= 0 {
		return ErrAttachmentNotFound
	}
	repository, ok := service.repository.(operatorAPIRepository)
	if !ok {
		return errors.New("operator API repository is required")
	}
	return repository.CanAccessOperatorAttachment(ctx, normalized, attachmentID)
}

type OperatorDashboardResponse struct {
	NewCount      int `json:"new_count"`
	AssignedCount int `json:"assigned_count"`
	ReturnedCount int `json:"returned_count"`
	CrisisCount   int `json:"crisis_count"`
	OverdueCount  int `json:"overdue_count"`
}

type OperatorClarificationResponse struct {
	QuestionID int64  `json:"question_id"`
	Question   string `json:"question"`
	AnswerID   int64  `json:"answer_id"`
	Answer     string `json:"answer"`
}

type OperatorWorkerResponse struct {
	ID            int64  `json:"id"`
	FullName      string `json:"full_name"`
	ExpertGroup   string `json:"expert_group"`
	IsResponsible bool   `json:"is_responsible"`
}

type OperatorNoteResponse struct {
	ID        int64                     `json:"id"`
	Text      string                    `json:"text"`
	Author    ResponsibleWorkerResponse `json:"author"`
	CreatedAt time.Time                 `json:"created_at"`
}

type OperatorEventResponse struct {
	ID         int64                      `json:"id"`
	EventType  string                     `json:"event_type"`
	Actor      *ResponsibleWorkerResponse `json:"actor,omitempty"`
	FromStatus string                     `json:"from_status,omitempty"`
	ToStatus   string                     `json:"to_status,omitempty"`
	Reason     string                     `json:"reason,omitempty"`
	CreatedAt  time.Time                  `json:"created_at"`
}

type OperatorTicketDetailResponse struct {
	TrackID           string                          `json:"track_id"`
	Status            string                          `json:"status"`
	Priority          string                          `json:"priority"`
	CreatedAt         time.Time                       `json:"created_at"`
	ClosedAt          *time.Time                      `json:"closed_at,omitempty"`
	Category          CategoryResponse                `json:"category"`
	ApplicantType     string                          `json:"applicant_type"`
	Description       string                          `json:"description"`
	ReturnCount       int                             `json:"return_count"`
	ReturnReason      string                          `json:"return_reason,omitempty"`
	CrisisDetected    bool                            `json:"crisis_detected"`
	CrisisContact     string                          `json:"crisis_contact,omitempty"`
	Clarifications    []OperatorClarificationResponse `json:"clarifications"`
	Attachments       []AttachmentResponse            `json:"attachments"`
	Workers           []OperatorWorkerResponse        `json:"workers"`
	Notes             []OperatorNoteResponse          `json:"notes"`
	RecommendedGroups []ExpertGroupResponse           `json:"recommended_groups"`
	Events            []OperatorEventResponse         `json:"events"`
}

type OperatorTicketUpdateRequest struct {
	CategoryID *int64  `json:"category_id"`
	Priority   *string `json:"priority"`
	Status     *string `json:"status"`
	Reason     string  `json:"reason"`
}

type OperatorCloseResponse struct {
	Status   string    `json:"status"`
	ClosedAt time.Time `json:"closed_at"`
	Message  string    `json:"message"`
}

func (service *OperatorService) Dashboard(ctx context.Context) (OperatorDashboardResponse, error) {
	repository, ok := service.repository.(operatorAPIRepository)
	if !ok {
		return OperatorDashboardResponse{}, errors.New("operator API repository is required")
	}
	now := service.now().UTC()
	record, err := repository.Dashboard(ctx, now.Add(-service.newSLA), now.Add(-service.responseSLA))
	if err != nil {
		return OperatorDashboardResponse{}, err
	}
	return OperatorDashboardResponse{NewCount: record.NewCount, AssignedCount: record.AssignedCount, ReturnedCount: record.ReturnedCount, CrisisCount: record.CrisisCount, OverdueCount: record.OverdueCount}, nil
}

func (service *OperatorService) Ticket(ctx context.Context, trackID string) (OperatorTicketDetailResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return OperatorTicketDetailResponse{}, err
	}
	repository, ok := service.repository.(operatorAPIRepository)
	if !ok {
		return OperatorTicketDetailResponse{}, errors.New("operator API repository is required")
	}
	record, err := repository.GetOperatorTicket(ctx, normalized)
	if err != nil {
		return OperatorTicketDetailResponse{}, err
	}
	return operatorTicketDetailResponse(record), nil
}

func operatorTicketDetailResponse(record repos.OperatorTicketDetailRecord) OperatorTicketDetailResponse {
	result := OperatorTicketDetailResponse{TrackID: record.TrackID, Status: record.Status.String(), Priority: record.Priority.String(), CreatedAt: record.CreatedAt.UTC(), Category: CategoryResponse{ID: record.CategoryID, Name: record.Category}, ApplicantType: record.ApplicantType.String(), Description: record.Description, ReturnCount: record.ReturnCount, ReturnReason: record.ReturnReason.String, CrisisDetected: record.CrisisDetected, CrisisContact: record.CrisisContact.String, Clarifications: []OperatorClarificationResponse{}, Attachments: []AttachmentResponse{}, Workers: []OperatorWorkerResponse{}, Notes: []OperatorNoteResponse{}, RecommendedGroups: []ExpertGroupResponse{}, Events: []OperatorEventResponse{}}
	if record.ClosedAt.Valid {
		value := record.ClosedAt.Time.UTC()
		result.ClosedAt = &value
	}
	for _, item := range record.Clarifications {
		result.Clarifications = append(result.Clarifications, OperatorClarificationResponse{QuestionID: item.QuestionID, Question: item.Question, AnswerID: item.AnswerID, Answer: item.Answer})
	}
	for _, item := range record.Attachments {
		result.Attachments = append(result.Attachments, AttachmentResponse{ID: item.ID, Name: item.SafeName, MIMEType: item.MIMEType, SizeBytes: item.SizeBytes})
	}
	for _, item := range record.Workers {
		result.Workers = append(result.Workers, OperatorWorkerResponse{ID: item.ID, FullName: item.FullName, ExpertGroup: item.GroupTitle, IsResponsible: item.IsResponsible})
	}
	for _, item := range record.Notes {
		result.Notes = append(result.Notes, OperatorNoteResponse{ID: item.ID, Text: item.Text, Author: ResponsibleWorkerResponse{ID: item.AuthorID, FullName: item.AuthorName}, CreatedAt: item.CreatedAt.UTC()})
	}
	for _, item := range record.RecommendedGroups {
		result.RecommendedGroups = append(result.RecommendedGroups, ExpertGroupResponse{ID: item.ID, Title: item.Title})
	}
	for _, item := range record.Events {
		event := OperatorEventResponse{ID: item.ID, EventType: item.EventType, Reason: item.Reason.String, CreatedAt: item.CreatedAt.UTC()}
		if item.ActorID.Valid {
			event.Actor = &ResponsibleWorkerResponse{ID: item.ActorID.Int64, FullName: item.ActorName.String}
		}
		if item.FromStatus.Valid {
			event.FromStatus = models.TicketStatus(item.FromStatus.Int64).String()
		}
		if item.ToStatus.Valid {
			event.ToStatus = models.TicketStatus(item.ToStatus.Int64).String()
		}
		result.Events = append(result.Events, event)
	}
	return result
}

func (service *OperatorService) UpdateTicket(ctx context.Context, trackID string, actorID int64, request OperatorTicketUpdateRequest) (OperatorTicketDetailResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return OperatorTicketDetailResponse{}, err
	}
	if actorID <= 0 {
		return OperatorTicketDetailResponse{}, ErrInvalidWorkerID
	}
	if request.CategoryID == nil && request.Priority == nil && request.Status == nil {
		return OperatorTicketDetailResponse{}, ErrInvalidFilter
	}
	if request.CategoryID != nil && *request.CategoryID <= 0 {
		return OperatorTicketDetailResponse{}, ErrInvalidFilter
	}
	reason := strings.TrimSpace(request.Reason)
	if utf8.RuneCountInString(reason) > 40 {
		return OperatorTicketDetailResponse{}, ErrInvalidFilter
	}
	update := repos.OperatorTicketUpdate{CategoryID: request.CategoryID, Reason: reason}
	if request.Priority != nil {
		value := ticketPriorityFromString(*request.Priority)
		if value == 0 {
			return OperatorTicketDetailResponse{}, ErrInvalidFilter
		}
		update.Priority = &value
	}
	if request.Status != nil {
		value := statusFromString(strings.TrimSpace(*request.Status))
		if value == 0 {
			return OperatorTicketDetailResponse{}, ErrInvalidFilter
		}
		update.Status = &value
	}
	repository, ok := service.repository.(operatorAPIRepository)
	if !ok {
		return OperatorTicketDetailResponse{}, errors.New("operator API repository is required")
	}
	if err := repository.UpdateOperatorTicket(ctx, normalized, actorID, update, service.now().UTC()); err != nil {
		return OperatorTicketDetailResponse{}, err
	}
	record, err := repository.GetOperatorTicket(ctx, normalized)
	if err != nil {
		return OperatorTicketDetailResponse{}, err
	}
	return operatorTicketDetailResponse(record), nil
}

func (service *OperatorService) CloseTicket(ctx context.Context, trackID string, actorID int64, message string) (OperatorCloseResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return OperatorCloseResponse{}, err
	}
	if actorID <= 0 {
		return OperatorCloseResponse{}, ErrInvalidWorkerID
	}
	message = strings.TrimSpace(message)
	if message == "" {
		return OperatorCloseResponse{}, ErrOperatorMessageRequired
	}
	repository, ok := service.repository.(operatorAPIRepository)
	if !ok {
		return OperatorCloseResponse{}, errors.New("operator API repository is required")
	}
	record, err := repository.CloseOperatorTicket(ctx, normalized, actorID, message, service.now().UTC())
	if err != nil {
		return OperatorCloseResponse{}, err
	}
	return OperatorCloseResponse{Status: record.Status.String(), ClosedAt: record.ClosedAt.UTC(), Message: record.Message}, nil
}
