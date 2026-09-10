package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	_ "time/tzdata"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
)

var (
	ErrExpertAccessDenied    = repos.ErrExpertAccessDenied
	ErrWorkerRequestNotFound = repos.ErrWorkerRequestNotFound
	ErrNoteTextRequired      = errors.New("note text is required")
	ErrInvalidRequestType    = errors.New("invalid worker request type")
	ErrRequestReasonRequired = errors.New("worker request reason is required")
)

type expertRepository interface {
	Profile(context.Context, int64) (repos.ExpertProfileRecord, error)
	Dashboard(context.Context, int64) (repos.ExpertDashboardRecord, error)
	ListTickets(context.Context, repos.ExpertTicketFilter) (repos.ExpertTicketPageRecord, error)
	OpenTicket(context.Context, string, int64, time.Time) (models.TicketStatus, error)
	GetTicket(context.Context, string, int64) (repos.ExpertTicketDetailRecord, error)
	AddNote(context.Context, string, int64, string, string, time.Time) (repos.ExpertNoteRecord, error)
	AddAnswer(context.Context, string, int64, string, time.Time) (repos.ExpertMessageRecord, error)
	CreateWorkerRequest(context.Context, string, int64, string, string, time.Time) (repos.WorkerRequestRecord, error)
	ListWorkerRequests(context.Context, int64, string, int, int) (repos.WorkerRequestPageRecord, error)
	CanAccessTicket(context.Context, string, int64) error
	Analytics(context.Context, int64, time.Time, time.Time) (repos.AnalyticsRecord, error)
	ReportTickets(context.Context, int64, time.Time, time.Time) ([]repos.ReportTicketRecord, error)
}

type ExpertService struct {
	repository expertRepository
	location   *time.Location
	now        func() time.Time
}

type ExpertGroupProfileResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type ExpertProfileResponse struct {
	ID            int64                      `json:"id"`
	FullName      string                     `json:"full_name"`
	Role          string                     `json:"role"`
	ExpertGroup   ExpertGroupProfileResponse `json:"expert_group"`
	ActiveTickets int                        `json:"active_tickets"`
	MaxTickets    int                        `json:"max_tickets"`
	AverageRating float64                    `json:"avg_rating"`
}

type ExpertDashboardResponse struct {
	QueueCount      int `json:"queue_count"`
	InProgressCount int `json:"in_progress_count"`
	ReturnedCount   int `json:"returned_count"`
	UrgentCount     int `json:"urgent_count"`
	ProcessedCount  int `json:"processed_count"`
}

type ExpertTicketListRequest struct {
	Queue         string
	Search        string
	Priority      string
	CategoryID    string
	ApplicantType string
	Page          string
	Limit         string
}

type ExpertTicketListItemResponse struct {
	TrackID        string           `json:"track_id"`
	Category       CategoryResponse `json:"category"`
	Status         string           `json:"status"`
	ApplicantType  string           `json:"applicant_type"`
	Priority       string           `json:"priority"`
	CreatedAt      time.Time        `json:"created_at"`
	WaitingSeconds int64            `json:"waiting_seconds"`
	ReturnCount    int              `json:"return_count"`
	IsResponsible  bool             `json:"is_responsible"`
}

type ExpertTicketPageResponse struct {
	Items []ExpertTicketListItemResponse `json:"items"`
	Total int                            `json:"total"`
	Page  int                            `json:"page"`
	Limit int                            `json:"limit"`
}

type ExpertClarificationResponse struct {
	QuestionID int64  `json:"question_id"`
	Question   string `json:"question"`
	AnswerID   int64  `json:"answer_id"`
	Answer     string `json:"answer"`
}

type ExpertWorkerResponse struct {
	ID            int64  `json:"id"`
	FullName      string `json:"full_name"`
	ExpertGroup   string `json:"expert_group"`
	IsResponsible bool   `json:"is_responsible"`
}

type ExpertNoteResponse struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	Author    AuthUser  `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type ExpertTicketResponse struct {
	TrackID        string                        `json:"track_id"`
	Status         string                        `json:"status"`
	Priority       string                        `json:"priority"`
	CreatedAt      time.Time                     `json:"created_at"`
	Category       CategoryResponse              `json:"category"`
	ApplicantType  string                        `json:"applicant_type"`
	Description    string                        `json:"description"`
	ReturnCount    int                           `json:"return_count"`
	ReturnReason   string                        `json:"return_reason,omitempty"`
	Clarifications []ExpertClarificationResponse `json:"clarifications"`
	Attachments    []AttachmentResponse          `json:"attachments"`
	Messages       []MessageResponse             `json:"messages"`
	Workers        []ExpertWorkerResponse        `json:"workers"`
	Notes          []ExpertNoteResponse          `json:"notes"`
	AllowedActions []string                      `json:"allowed_actions"`
}

type ExpertOpenResponse struct {
	Status string `json:"status"`
}

type CreateExpertNoteResponse struct {
	Note ExpertNoteResponse `json:"note"`
}

type CreateExpertAnswerResponse struct {
	Message      MessageResponse `json:"message"`
	TicketStatus string          `json:"ticket_status"`
}

type CreateWorkerRequestResponse struct {
	Request WorkerRequestResponse `json:"request"`
}

type WorkerRequestResponse struct {
	ID            int64            `json:"id"`
	TrackID       string           `json:"track_id"`
	Category      CategoryResponse `json:"category,omitempty"`
	ApplicantType string           `json:"applicant_type,omitempty"`
	Priority      string           `json:"priority,omitempty"`
	RequestType   string           `json:"request_type"`
	Reason        string           `json:"reason"`
	Status        string           `json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
	CreatedBy     int64            `json:"created_by,omitempty"`
}

type WorkerRequestPageResponse struct {
	Items []WorkerRequestResponse `json:"items"`
	Total int                     `json:"total"`
	Page  int                     `json:"page"`
	Limit int                     `json:"limit"`
}

func NewExpertService(repository expertRepository) (*ExpertService, error) {
	if repository == nil {
		return nil, errors.New("expert repository is required")
	}
	location, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		return nil, fmt.Errorf("load expert analytics timezone: %w", err)
	}
	return &ExpertService{repository: repository, location: location, now: time.Now}, nil
}

func (service *ExpertService) Profile(ctx context.Context, workerID int64) (ExpertProfileResponse, error) {
	record, err := service.repository.Profile(ctx, workerID)
	if err != nil {
		return ExpertProfileResponse{}, err
	}
	return ExpertProfileResponse{
		ID: record.ID, FullName: record.FullName, Role: "expert",
		ExpertGroup:   ExpertGroupProfileResponse{ID: record.GroupID, Title: record.GroupTitle},
		ActiveTickets: record.ActiveTickets, MaxTickets: record.MaxTickets, AverageRating: record.AverageRating,
	}, nil
}

func (service *ExpertService) Dashboard(ctx context.Context, workerID int64) (ExpertDashboardResponse, error) {
	record, err := service.repository.Dashboard(ctx, workerID)
	if err != nil {
		return ExpertDashboardResponse{}, err
	}
	return ExpertDashboardResponse{
		QueueCount: record.QueueCount, InProgressCount: record.InProgressCount,
		ReturnedCount: record.ReturnedCount, UrgentCount: record.UrgentCount,
		ProcessedCount: record.ProcessedCount,
	}, nil
}

func (service *ExpertService) ListTickets(ctx context.Context, workerID int64, request ExpertTicketListRequest) (ExpertTicketPageResponse, error) {
	filter, page, err := parseExpertTicketFilter(workerID, request)
	if err != nil {
		return ExpertTicketPageResponse{}, err
	}
	record, err := service.repository.ListTickets(ctx, filter)
	if err != nil {
		return ExpertTicketPageResponse{}, err
	}
	response := ExpertTicketPageResponse{Items: []ExpertTicketListItemResponse{}, Total: record.Total, Page: page, Limit: filter.Limit}
	for _, item := range record.Items {
		response.Items = append(response.Items, ExpertTicketListItemResponse{
			TrackID: item.TrackID, Category: CategoryResponse{ID: item.CategoryID, Name: item.Category},
			Status: item.Status.String(), ApplicantType: item.ApplicantType.String(), Priority: item.Priority.String(),
			CreatedAt: item.CreatedAt.UTC(), WaitingSeconds: item.WaitingSeconds,
			ReturnCount: item.ReturnCount, IsResponsible: item.IsResponsible,
		})
	}
	return response, nil
}

func parseExpertTicketFilter(workerID int64, request ExpertTicketListRequest) (repos.ExpertTicketFilter, int, error) {
	filter := repos.ExpertTicketFilter{WorkerID: workerID, Queue: strings.TrimSpace(request.Queue), Search: strings.TrimSpace(request.Search), Limit: 20}
	if workerID <= 0 || (filter.Queue != "queue" && filter.Queue != "assigned" && filter.Queue != "returned") {
		return repos.ExpertTicketFilter{}, 0, ErrInvalidQueue
	}
	page, err := positiveIntOrDefault(request.Page, 1)
	if err != nil {
		return repos.ExpertTicketFilter{}, 0, ErrInvalidFilter
	}
	filter.Limit, err = positiveIntOrDefault(request.Limit, 20)
	if err != nil || filter.Limit > 100 {
		return repos.ExpertTicketFilter{}, 0, ErrInvalidFilter
	}
	filter.Offset = (page - 1) * filter.Limit
	if request.Priority != "" {
		switch request.Priority {
		case "standard":
			filter.Priority = models.TicketPriorityStandard
		case "urgent":
			filter.Priority = models.TicketPriorityUrgent
		default:
			return repos.ExpertTicketFilter{}, 0, ErrInvalidFilter
		}
	}
	filter.CategoryID, err = optionalPositiveInt64(request.CategoryID)
	if err != nil {
		return repos.ExpertTicketFilter{}, 0, ErrInvalidFilter
	}
	if request.ApplicantType != "" {
		switch request.ApplicantType {
		case "schoolchild":
			filter.ApplicantType = models.ApplicantTypeSchoolchild
		case "parent":
			filter.ApplicantType = models.ApplicantTypeParent
		case "teacher":
			filter.ApplicantType = models.ApplicantTypeTeacher
		default:
			return repos.ExpertTicketFilter{}, 0, ErrInvalidFilter
		}
	}
	return filter, page, nil
}

func (service *ExpertService) OpenTicket(ctx context.Context, trackID string, workerID int64) (ExpertOpenResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return ExpertOpenResponse{}, err
	}
	status, err := service.repository.OpenTicket(ctx, normalized, workerID, service.now().UTC())
	if err != nil {
		return ExpertOpenResponse{}, err
	}
	return ExpertOpenResponse{Status: status.String()}, nil
}

func (service *ExpertService) Ticket(ctx context.Context, trackID string, workerID int64) (ExpertTicketResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return ExpertTicketResponse{}, err
	}
	record, err := service.repository.GetTicket(ctx, normalized, workerID)
	if err != nil {
		return ExpertTicketResponse{}, err
	}
	response := ExpertTicketResponse{
		TrackID: record.TrackID, Status: record.Status.String(), Priority: record.Priority.String(),
		CreatedAt: record.CreatedAt.UTC(), Category: CategoryResponse{ID: record.CategoryID, Name: record.Category},
		ApplicantType: record.ApplicantType.String(), Description: record.Description,
		ReturnCount: record.ReturnCount, Clarifications: []ExpertClarificationResponse{}, Attachments: []AttachmentResponse{},
		Messages: []MessageResponse{}, Workers: []ExpertWorkerResponse{}, Notes: []ExpertNoteResponse{},
		AllowedActions: expertAllowedActions(record.Status),
	}
	if record.ReturnReason.Valid {
		response.ReturnReason = record.ReturnReason.String
	}
	for _, item := range record.Clarifications {
		response.Clarifications = append(response.Clarifications, ExpertClarificationResponse{QuestionID: item.QuestionID, Question: item.Question, AnswerID: item.AnswerID, Answer: item.Answer})
	}
	for _, item := range record.Attachments {
		response.Attachments = append(response.Attachments, AttachmentResponse{ID: item.ID, Name: item.SafeName, MIMEType: item.MIMEType, SizeBytes: item.SizeBytes})
	}
	for _, item := range record.Messages {
		response.Messages = append(response.Messages, MessageResponse{ID: item.ID, Text: item.Text, Type: item.Type.String(), CreatedAt: item.CreatedAt.UTC(), Attachments: []AttachmentResponse{}})
	}
	for _, item := range record.Workers {
		response.Workers = append(response.Workers, ExpertWorkerResponse{ID: item.ID, FullName: item.FullName, ExpertGroup: item.GroupTitle, IsResponsible: item.IsResponsible})
	}
	for _, item := range record.Notes {
		response.Notes = append(response.Notes, ExpertNoteResponse{ID: item.ID, Text: item.Text, Author: AuthUser{ID: item.AuthorID, FullName: item.AuthorName, Role: "expert"}, CreatedAt: item.CreatedAt.UTC()})
	}
	return response, nil
}

func expertAllowedActions(status models.TicketStatus) []string {
	actions := []string{"add_note", "create_worker_request", "download_attachment"}
	switch status {
	case models.TicketStatusAssigned:
		return append(actions, "open", "send_answer")
	case models.TicketStatusInProgress, models.TicketStatusNeedsClarification:
		return append(actions, "send_answer")
	default:
		return actions
	}
}

func (service *ExpertService) AddNote(ctx context.Context, trackID string, user AuthUser, text string) (CreateExpertNoteResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return CreateExpertNoteResponse{}, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return CreateExpertNoteResponse{}, ErrNoteTextRequired
	}
	record, err := service.repository.AddNote(ctx, normalized, user.ID, user.FullName, text, service.now().UTC())
	if err != nil {
		return CreateExpertNoteResponse{}, err
	}
	return CreateExpertNoteResponse{Note: ExpertNoteResponse{ID: record.ID, Text: record.Text, Author: user, CreatedAt: record.CreatedAt.UTC()}}, nil
}

func (service *ExpertService) AddAnswer(ctx context.Context, trackID string, workerID int64, text string) (CreateExpertAnswerResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return CreateExpertAnswerResponse{}, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return CreateExpertAnswerResponse{}, ErrMessageTextRequired
	}
	record, err := service.repository.AddAnswer(ctx, normalized, workerID, text, service.now().UTC())
	if err != nil {
		return CreateExpertAnswerResponse{}, err
	}
	return CreateExpertAnswerResponse{
		Message:      MessageResponse{ID: record.ID, Text: record.Text, Type: record.Type.String(), CreatedAt: record.CreatedAt.UTC(), Attachments: []AttachmentResponse{}},
		TicketStatus: record.TicketStatus.String(),
	}, nil
}

func (service *ExpertService) CreateWorkerRequest(ctx context.Context, trackID string, workerID int64, requestType, reason string) (CreateWorkerRequestResponse, error) {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return CreateWorkerRequestResponse{}, err
	}
	requestType = strings.TrimSpace(requestType)
	if requestType != "add_coworker" && requestType != "replace_responsible" {
		return CreateWorkerRequestResponse{}, ErrInvalidRequestType
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return CreateWorkerRequestResponse{}, ErrRequestReasonRequired
	}
	record, err := service.repository.CreateWorkerRequest(ctx, normalized, workerID, requestType, reason, service.now().UTC())
	if err != nil {
		return CreateWorkerRequestResponse{}, err
	}
	return CreateWorkerRequestResponse{Request: workerRequestResponse(record)}, nil
}

func (service *ExpertService) ListWorkerRequests(ctx context.Context, workerID int64, status, pageValue, limitValue string) (WorkerRequestPageResponse, error) {
	status = strings.TrimSpace(status)
	if status != "" && status != "sent" && status != "completed" {
		return WorkerRequestPageResponse{}, ErrInvalidFilter
	}
	page, err := positiveIntOrDefault(pageValue, 1)
	if err != nil {
		return WorkerRequestPageResponse{}, ErrInvalidFilter
	}
	limit, err := positiveIntOrDefault(limitValue, 20)
	if err != nil || limit > 100 {
		return WorkerRequestPageResponse{}, ErrInvalidFilter
	}
	record, err := service.repository.ListWorkerRequests(ctx, workerID, status, limit, (page-1)*limit)
	if err != nil {
		return WorkerRequestPageResponse{}, err
	}
	response := WorkerRequestPageResponse{Items: []WorkerRequestResponse{}, Total: record.Total, Page: page, Limit: limit}
	for _, item := range record.Items {
		response.Items = append(response.Items, workerRequestResponse(item))
	}
	return response, nil
}

func workerRequestResponse(record repos.WorkerRequestRecord) WorkerRequestResponse {
	return WorkerRequestResponse{
		ID: record.ID, TrackID: record.TrackID, Category: CategoryResponse{ID: record.CategoryID, Name: record.Category},
		ApplicantType: record.ApplicantType.String(), Priority: record.Priority.String(), RequestType: record.RequestType,
		Reason: record.Reason, Status: record.Status, CreatedAt: record.CreatedAt.UTC(),
		CreatedBy: record.CreatedBy,
	}
}

func (service *ExpertService) CanAccessTicket(ctx context.Context, trackID string, workerID int64) error {
	normalized, err := normalizeTrackID(trackID)
	if err != nil {
		return err
	}
	return service.repository.CanAccessTicket(ctx, normalized, workerID)
}

func (service *ExpertService) Analytics(ctx context.Context, workerID int64, dateFrom, dateTo string) (AnalyticsResponse, error) {
	start, end, err := service.dateRange(dateFrom, dateTo)
	if err != nil {
		return AnalyticsResponse{}, err
	}
	record, err := service.repository.Analytics(ctx, workerID, start, end)
	if err != nil {
		return AnalyticsResponse{}, err
	}
	return analyticsResponse(record, dateFrom, dateTo), nil
}

func analyticsResponse(record repos.AnalyticsRecord, dateFrom, dateTo string) AnalyticsResponse {
	result := AnalyticsResponse{
		DateFrom: dateFrom, DateTo: dateTo, TotalTickets: record.Total,
		Categories: []DistributionResponse{}, ApplicantTypes: []DistributionResponse{}, Statuses: []DistributionResponse{},
		AverageAcceptanceSeconds:    record.DurationSecondsToAccept,
		AverageFirstResponseSeconds: record.DurationSecondsToResponse,
		AverageCloseSeconds:         record.DurationSecondsToClose,
		ExpertLoadPercent:           record.ExpertLoadPercent,
		UrgentSharePercent:          percent(record.UrgentCount, record.Total), ReturnSharePercent: percent(record.ReturnedCount, record.Total),
	}
	for _, item := range record.Categories {
		result.Categories = append(result.Categories, DistributionResponse{ID: item.ID, Value: item.Name, Count: item.Count, Percent: percent(item.Count, record.Total)})
	}
	for _, item := range record.ApplicantTypes {
		result.ApplicantTypes = append(result.ApplicantTypes, DistributionResponse{Value: models.ApplicantType(item.Value).String(), Count: item.Count, Percent: percent(item.Count, record.Total)})
	}
	for _, item := range record.Statuses {
		result.Statuses = append(result.Statuses, DistributionResponse{Value: models.TicketStatus(item.Value).String(), Count: item.Count, Percent: percent(item.Count, record.Total)})
	}
	return result
}

func (service *ExpertService) Report(ctx context.Context, workerID int64, dateFrom, dateTo, format string) (GeneratedReport, error) {
	start, end, err := service.dateRange(dateFrom, dateTo)
	if err != nil {
		return GeneratedReport{}, err
	}
	records, err := service.repository.ReportTickets(ctx, workerID, start, end)
	if err != nil {
		return GeneratedReport{}, err
	}
	format = strings.ToLower(strings.TrimSpace(format))
	name := "expert-tickets-" + dateFrom + "-" + dateTo
	switch format {
	case "csv":
		data, err := buildCSVReport(records)
		return GeneratedReport{Name: name + ".csv", ContentType: "text/csv; charset=utf-8", Data: data}, err
	case "xlsx":
		data, err := buildXLSXReport(records)
		return GeneratedReport{Name: name + ".xlsx", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Data: data}, err
	default:
		return GeneratedReport{}, ErrInvalidReportFormat
	}
}

func (service *ExpertService) dateRange(dateFrom, dateTo string) (time.Time, time.Time, error) {
	start, err := time.ParseInLocation(time.DateOnly, dateFrom, service.location)
	if err != nil {
		return time.Time{}, time.Time{}, ErrInvalidDateRange
	}
	last, err := time.ParseInLocation(time.DateOnly, dateTo, service.location)
	if err != nil || last.Before(start) {
		return time.Time{}, time.Time{}, ErrInvalidDateRange
	}
	return start.UTC(), last.AddDate(0, 0, 1).UTC(), nil
}
