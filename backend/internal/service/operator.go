package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
)

var (
	ErrWorkerNotFound        = repos.ErrWorkerNotFound
	ErrWorkerUnavailable     = repos.ErrWorkerUnavailable
	ErrWorkerAlreadyAssigned = repos.ErrWorkerAlreadyAssigned
	ErrWorkerNotAssigned     = repos.ErrWorkerNotAssigned
	ErrResponsibleWorker     = repos.ErrResponsibleWorker
	ErrResponsibleRequired   = repos.ErrResponsibleRequired
	ErrInvalidWorkerID       = errors.New("invalid worker id")
	ErrInvalidQueue          = errors.New("invalid ticket queue")
	ErrInvalidFilter         = errors.New("invalid ticket filter")
	ErrInvalidReasonCode     = errors.New("invalid rejection reason code")
	ErrRejectionMessage      = errors.New("rejection message is required")
	ErrInvalidDateRange      = errors.New("invalid date range")
	ErrInvalidReportFormat   = errors.New("invalid report format")
	ErrWorkerRequestComplete = errors.New("worker request could not be completed")
)

type operatorWorkerRequestRepository interface {
	ListOperatorWorkerRequests(context.Context, string, int, int) (repos.WorkerRequestPageRecord, error)
	CompleteWorkerRequest(context.Context, int64, int64, int64, time.Time) (repos.WorkerRequestRecord, error)
}

type operatorRepository interface {
	EligibleWorkers(context.Context, string, string, int64) (repos.EligibleWorkersRecord, error)
	SetResponsibleWorker(context.Context, string, int64, int64, time.Time) (repos.AssignmentRecord, error)
	AddWorker(context.Context, string, int64, int64, time.Time) (repos.AssignmentRecord, error)
	RemoveWorker(context.Context, string, int64, int64, time.Time) (repos.AssignmentRecord, error)
	Reject(context.Context, string, int64, string, string, time.Time) (repos.RejectionRecord, error)
	ListTickets(context.Context, repos.OperatorTicketFilter) (repos.OperatorTicketPageRecord, error)
	Analytics(context.Context, time.Time, time.Time) (repos.AnalyticsRecord, error)
	ReportTickets(context.Context, time.Time, time.Time) ([]repos.ReportTicketRecord, error)
}

type OperatorService struct {
	repository operatorRepository
	location   *time.Location
	now        func() time.Time
}

type ExpertGroupResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type EligibleWorkerResponse struct {
	ID            int64               `json:"id"`
	FullName      string              `json:"full_name"`
	ExpertGroup   ExpertGroupResponse `json:"expert_group"`
	ActiveTickets int                 `json:"active_tickets"`
	MaxTickets    int                 `json:"max_tickets"`
	Recommended   bool                `json:"recommended"`
}

type EligibleWorkersResponse struct {
	RecommendedGroup  *ExpertGroupResponse     `json:"recommended_group"`
	RecommendedGroups []ExpertGroupResponse    `json:"recommended_groups"`
	Workers           []EligibleWorkerResponse `json:"workers"`
}

type AssignmentResponse struct {
	WorkerID int64  `json:"worker_id"`
	Status   string `json:"status"`
}

type RejectionResponse struct {
	Status            string    `json:"status"`
	ClosedAt          time.Time `json:"closed_at"`
	ResolutionMessage string    `json:"resolution_message"`
}

type TicketListRequest struct {
	Queue         string
	Search        string
	Priority      string
	Status        string
	CategoryID    string
	ApplicantType string
	Page          string
	Limit         string
}

type ResponsibleWorkerResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}

type OperatorTicketResponse struct {
	TrackID             string                     `json:"track_id"`
	Category            CategoryResponse           `json:"category"`
	Status              string                     `json:"status"`
	ApplicantType       string                     `json:"applicant_type"`
	Priority            string                     `json:"priority"`
	CreatedAt           time.Time                  `json:"created_at"`
	WaitingSeconds      int64                      `json:"waiting_seconds"`
	Responsible         *ResponsibleWorkerResponse `json:"responsible,omitempty"`
	ReturnCount         int                        `json:"return_count"`
	ReturnReason        string                     `json:"return_reason,omitempty"`
	ReturnedAt          *time.Time                 `json:"returned_at,omitempty"`
	PreviousResponsible *ResponsibleWorkerResponse `json:"previous_responsible,omitempty"`
}

type OperatorTicketPageResponse struct {
	Items []OperatorTicketResponse `json:"items"`
	Total int                      `json:"total"`
	Page  int                      `json:"page"`
	Limit int                      `json:"limit"`
}

type DistributionResponse struct {
	ID      int64   `json:"id,omitempty"`
	Value   string  `json:"value"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type AnalyticsResponse struct {
	DateFrom                    string                 `json:"date_from"`
	DateTo                      string                 `json:"date_to"`
	TotalTickets                int                    `json:"total_tickets"`
	Categories                  []DistributionResponse `json:"categories"`
	ApplicantTypes              []DistributionResponse `json:"applicant_types"`
	Statuses                    []DistributionResponse `json:"statuses"`
	AverageAcceptanceSeconds    *float64               `json:"average_acceptance_seconds"`
	AverageFirstResponseSeconds *float64               `json:"average_first_response_seconds"`
	AverageCloseSeconds         *float64               `json:"average_close_seconds"`
	ExpertLoadPercent           float64                `json:"expert_load_percent"`
	UrgentSharePercent          float64                `json:"urgent_share_percent"`
	ReturnSharePercent          float64                `json:"return_share_percent"`
}

type GeneratedReport struct {
	Name        string
	ContentType string
	Data        []byte
}

type CompleteWorkerRequestResponse struct {
	RequestID int64  `json:"request_id"`
	Status    string `json:"status"`
	WorkerID  int64  `json:"worker_id"`
}

func NewOperatorService(repository operatorRepository) (*OperatorService, error) {
	if repository == nil {
		return nil, errors.New("operator repository is required")
	}
	location, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		return nil, fmt.Errorf("load analytics timezone: %w", err)
	}
	return &OperatorService{repository: repository, location: location, now: time.Now}, nil
}

func (service *OperatorService) EligibleWorkers(ctx context.Context, trackID, name, groupIDValue string) (EligibleWorkersResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return EligibleWorkersResponse{}, err
	}
	groupID, err := optionalPositiveInt64(groupIDValue)
	if err != nil {
		return EligibleWorkersResponse{}, ErrInvalidFilter
	}
	record, err := service.repository.EligibleWorkers(ctx, trackID, strings.TrimSpace(name), groupID)
	if err != nil {
		return EligibleWorkersResponse{}, err
	}
	result := EligibleWorkersResponse{RecommendedGroups: make([]ExpertGroupResponse, 0, len(record.RecommendedGroups)), Workers: make([]EligibleWorkerResponse, 0, len(record.Workers))}
	if record.RecommendedGroup != nil {
		result.RecommendedGroup = &ExpertGroupResponse{ID: record.RecommendedGroup.ID, Title: record.RecommendedGroup.Title}
	}
	for _, group := range record.RecommendedGroups {
		result.RecommendedGroups = append(result.RecommendedGroups, ExpertGroupResponse{ID: group.ID, Title: group.Title})
	}
	for _, worker := range record.Workers {
		result.Workers = append(result.Workers, EligibleWorkerResponse{
			ID: worker.ID, FullName: worker.FullName,
			ExpertGroup:   ExpertGroupResponse{ID: worker.GroupID, Title: worker.GroupTitle},
			ActiveTickets: worker.ActiveTickets, MaxTickets: worker.MaxTickets,
			Recommended: worker.Recommended,
		})
	}
	return result, nil
}

func (service *OperatorService) SetResponsibleWorker(ctx context.Context, trackID string, workerID, actorID int64) (AssignmentResponse, error) {
	return service.changeWorker(ctx, trackID, workerID, actorID, service.repository.SetResponsibleWorker)
}

func (service *OperatorService) AddWorker(ctx context.Context, trackID string, workerID, actorID int64) (AssignmentResponse, error) {
	return service.changeWorker(ctx, trackID, workerID, actorID, service.repository.AddWorker)
}

func (service *OperatorService) RemoveWorker(ctx context.Context, trackID string, workerID, actorID int64) (AssignmentResponse, error) {
	return service.changeWorker(ctx, trackID, workerID, actorID, service.repository.RemoveWorker)
}

func (service *OperatorService) changeWorker(ctx context.Context, trackID string, workerID, actorID int64, operation func(context.Context, string, int64, int64, time.Time) (repos.AssignmentRecord, error)) (AssignmentResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return AssignmentResponse{}, err
	}
	if workerID <= 0 || actorID <= 0 {
		return AssignmentResponse{}, ErrInvalidWorkerID
	}
	record, err := operation(ctx, trackID, workerID, actorID, service.now().UTC())
	if err != nil {
		return AssignmentResponse{}, err
	}
	return AssignmentResponse{WorkerID: record.WorkerID, Status: record.Status.String()}, nil
}

func (service *OperatorService) Reject(ctx context.Context, trackID string, actorID int64, reasonCode, message string) (RejectionResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return RejectionResponse{}, err
	}
	reasonCode = strings.TrimSpace(reasonCode)
	switch reasonCode {
	case "no_expert_help", "spam", "out_of_scope":
	default:
		return RejectionResponse{}, ErrInvalidReasonCode
	}
	message = strings.TrimSpace(message)
	if message == "" {
		return RejectionResponse{}, ErrRejectionMessage
	}
	record, err := service.repository.Reject(ctx, trackID, actorID, reasonCode, message, service.now().UTC())
	if err != nil {
		return RejectionResponse{}, err
	}
	return RejectionResponse{Status: record.Status.String(), ClosedAt: record.ClosedAt.UTC(), ResolutionMessage: record.Message}, nil
}

func (service *OperatorService) ListTickets(ctx context.Context, request TicketListRequest) (OperatorTicketPageResponse, error) {
	filter, page, err := parseTicketFilter(request)
	if err != nil {
		return OperatorTicketPageResponse{}, err
	}
	record, err := service.repository.ListTickets(ctx, filter)
	if err != nil {
		return OperatorTicketPageResponse{}, err
	}
	result := OperatorTicketPageResponse{Items: make([]OperatorTicketResponse, 0, len(record.Items)), Total: record.Total, Page: page, Limit: filter.Limit}
	for _, item := range record.Items {
		response := OperatorTicketResponse{
			TrackID: item.TrackID, Category: CategoryResponse{ID: item.CategoryID, Name: item.Category},
			Status: item.Status.String(), ApplicantType: item.ApplicantType.String(), Priority: item.Priority.String(),
			CreatedAt: item.CreatedAt.UTC(), WaitingSeconds: item.WaitingSeconds,
			ReturnCount: item.ReturnCount, ReturnReason: item.ReturnReason.String,
		}
		if item.ResponsibleID.Valid {
			response.Responsible = &ResponsibleWorkerResponse{ID: item.ResponsibleID.Int64, FullName: item.ResponsibleName.String}
		}
		if item.ReturnedAt.Valid {
			returnedAt := item.ReturnedAt.Time.UTC()
			response.ReturnedAt = &returnedAt
		}
		if item.PreviousWorkerID.Valid {
			response.PreviousResponsible = &ResponsibleWorkerResponse{ID: item.PreviousWorkerID.Int64, FullName: item.PreviousWorkerName.String}
		}
		result.Items = append(result.Items, response)
	}
	return result, nil
}

func (service *OperatorService) Analytics(ctx context.Context, dateFrom, dateTo string) (AnalyticsResponse, error) {
	start, end, err := service.dateRange(dateFrom, dateTo)
	if err != nil {
		return AnalyticsResponse{}, err
	}
	record, err := service.repository.Analytics(ctx, start, end)
	if err != nil {
		return AnalyticsResponse{}, err
	}
	result := AnalyticsResponse{
		DateFrom: dateFrom, DateTo: dateTo, TotalTickets: record.Total,
		Categories:                  make([]DistributionResponse, 0, len(record.Categories)),
		ApplicantTypes:              make([]DistributionResponse, 0, len(record.ApplicantTypes)),
		Statuses:                    make([]DistributionResponse, 0, len(record.Statuses)),
		AverageAcceptanceSeconds:    record.DurationSecondsToAccept,
		AverageFirstResponseSeconds: record.DurationSecondsToResponse,
		AverageCloseSeconds:         record.DurationSecondsToClose,
		ExpertLoadPercent:           record.ExpertLoadPercent,
		UrgentSharePercent:          percent(record.UrgentCount, record.Total),
		ReturnSharePercent:          percent(record.ReturnedCount, record.Total),
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
	return result, nil
}

func (service *OperatorService) Report(ctx context.Context, dateFrom, dateTo, format string) (GeneratedReport, error) {
	start, end, err := service.dateRange(dateFrom, dateTo)
	if err != nil {
		return GeneratedReport{}, err
	}
	records, err := service.repository.ReportTickets(ctx, start, end)
	if err != nil {
		return GeneratedReport{}, err
	}
	format = strings.ToLower(strings.TrimSpace(format))
	name := "tickets-" + dateFrom + "-" + dateTo
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

func (service *OperatorService) ListWorkerRequests(ctx context.Context, status, pageValue, limitValue string) (WorkerRequestPageResponse, error) {
	repository, ok := service.repository.(operatorWorkerRequestRepository)
	if !ok {
		return WorkerRequestPageResponse{}, ErrWorkerRequestComplete
	}
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
	record, err := repository.ListOperatorWorkerRequests(ctx, status, limit, (page-1)*limit)
	if err != nil {
		return WorkerRequestPageResponse{}, err
	}
	response := WorkerRequestPageResponse{Items: []WorkerRequestResponse{}, Total: record.Total, Page: page, Limit: limit}
	for _, item := range record.Items {
		response.Items = append(response.Items, workerRequestResponse(item))
	}
	return response, nil
}

func (service *OperatorService) CompleteWorkerRequest(ctx context.Context, requestID, workerID, actorID int64) (CompleteWorkerRequestResponse, error) {
	if requestID <= 0 || workerID <= 0 || actorID <= 0 {
		return CompleteWorkerRequestResponse{}, ErrInvalidWorkerID
	}
	repository, ok := service.repository.(operatorWorkerRequestRepository)
	if !ok {
		return CompleteWorkerRequestResponse{}, ErrWorkerRequestComplete
	}
	_, err := repository.CompleteWorkerRequest(ctx, requestID, workerID, actorID, service.now().UTC())
	if err != nil {
		return CompleteWorkerRequestResponse{}, err
	}
	return CompleteWorkerRequestResponse{RequestID: requestID, Status: "completed", WorkerID: workerID}, nil
}

func (service *OperatorService) dateRange(dateFrom, dateTo string) (time.Time, time.Time, error) {
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

func parseTicketFilter(request TicketListRequest) (repos.OperatorTicketFilter, int, error) {
	filter := repos.OperatorTicketFilter{Queue: strings.TrimSpace(request.Queue), Search: strings.TrimSpace(request.Search), Limit: 20}
	if filter.Queue != "assigned" && filter.Queue != "returned" {
		return repos.OperatorTicketFilter{}, 0, ErrInvalidQueue
	}
	page, err := positiveIntOrDefault(request.Page, 1)
	if err != nil {
		return repos.OperatorTicketFilter{}, 0, ErrInvalidFilter
	}
	filter.Limit, err = positiveIntOrDefault(request.Limit, 20)
	if err != nil || filter.Limit > 100 {
		return repos.OperatorTicketFilter{}, 0, ErrInvalidFilter
	}
	filter.Offset = (page - 1) * filter.Limit
	if request.Priority != "" {
		switch request.Priority {
		case "standard":
			filter.Priority = models.TicketPriorityStandard
		case "urgent":
			filter.Priority = models.TicketPriorityUrgent
		default:
			return repos.OperatorTicketFilter{}, 0, ErrInvalidFilter
		}
	}
	if request.Status != "" {
		filter.Status = statusFromString(request.Status)
		if filter.Status == 0 {
			return repos.OperatorTicketFilter{}, 0, ErrInvalidFilter
		}
	}
	filter.CategoryID, err = optionalPositiveInt64(request.CategoryID)
	if err != nil {
		return repos.OperatorTicketFilter{}, 0, ErrInvalidFilter
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
			return repos.OperatorTicketFilter{}, 0, ErrInvalidFilter
		}
	}
	return filter, page, nil
}

func statusFromString(value string) models.TicketStatus {
	for status := models.TicketStatusNew; status <= models.TicketStatusClosedWithoutAnswer; status++ {
		if status.String() == value {
			return status
		}
	}
	return 0
}

func optionalPositiveInt64(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return 0, ErrInvalidFilter
	}
	return parsed, nil
}

func positiveIntOrDefault(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, ErrInvalidFilter
	}
	return parsed, nil
}

func percent(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) * 100 / float64(total)
}
