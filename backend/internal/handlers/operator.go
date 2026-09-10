package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strconv"

	"example.com/german/backend/internal/service"
)

const maxOperatorRequestBytes = 64 * 1024

type operatorService interface {
	EligibleWorkers(context.Context, string, string, string) (service.EligibleWorkersResponse, error)
	SetResponsibleWorker(context.Context, string, int64, int64) (service.AssignmentResponse, error)
	AddWorker(context.Context, string, int64, int64) (service.AssignmentResponse, error)
	RemoveWorker(context.Context, string, int64, int64) (service.AssignmentResponse, error)
	Reject(context.Context, string, int64, string, string) (service.RejectionResponse, error)
	ListTickets(context.Context, service.TicketListRequest) (service.OperatorTicketPageResponse, error)
	Analytics(context.Context, string, string) (service.AnalyticsResponse, error)
	Report(context.Context, string, string, string) (service.GeneratedReport, error)
}

type operatorWorkerRequestService interface {
	ListWorkerRequests(context.Context, string, string, string) (service.WorkerRequestPageResponse, error)
	CompleteWorkerRequest(context.Context, int64, int64, int64) (service.CompleteWorkerRequestResponse, error)
}

type OperatorHandler struct {
	service operatorService
	auth    authService
}

type workerRequest struct {
	WorkerID int64 `json:"worker_id"`
}

type rejectRequest struct {
	ReasonCode string `json:"reason_code"`
	Message    string `json:"message"`
}

func NewOperatorHandler(operatorService operatorService, auth authService) (*OperatorHandler, error) {
	if operatorService == nil || auth == nil {
		return nil, errors.New("operator and auth services are required")
	}
	return &OperatorHandler{service: operatorService, auth: auth}, nil
}

func (handler *OperatorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/operator/tickets", handler.listTickets)
	mux.HandleFunc("GET /api/operator/tickets/{track_id}/eligible-workers", handler.eligibleWorkers)
	mux.HandleFunc("PUT /api/operator/tickets/{track_id}/responsible-worker", handler.setResponsibleWorker)
	mux.HandleFunc("POST /api/operator/tickets/{track_id}/workers", handler.addWorker)
	mux.HandleFunc("DELETE /api/operator/tickets/{track_id}/workers/{worker_id}", handler.removeWorker)
	mux.HandleFunc("POST /api/operator/tickets/{track_id}/reject", handler.rejectTicket)
	mux.HandleFunc("GET /api/operator/analytics", handler.analytics)
	mux.HandleFunc("GET /api/operator/reports", handler.report)
	mux.HandleFunc("GET /api/operator/worker-requests", handler.listWorkerRequests)
	mux.HandleFunc("POST /api/operator/worker-requests/{request_id}/complete", handler.completeWorkerRequest)
}

func (handler *OperatorHandler) listWorkerRequests(w http.ResponseWriter, request *http.Request) {
	if _, ok := handler.operatorUser(w, request); !ok {
		return
	}
	requestService, ok := handler.service.(operatorWorkerRequestService)
	if !ok {
		writeAPIError(w, http.StatusInternalServerError, "internal_error", http.StatusText(http.StatusInternalServerError))
		return
	}
	query := request.URL.Query()
	response, err := requestService.ListWorkerRequests(request.Context(), query.Get("status"), query.Get("page"), query.Get("limit"))
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) completeWorkerRequest(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.operatorUser(w, request)
	if !ok {
		return
	}
	requestID, err := strconv.ParseInt(request.PathValue("request_id"), 10, 64)
	if err != nil || requestID <= 0 {
		writeAPIError(w, http.StatusNotFound, "worker_request_not_found", "worker request not found")
		return
	}
	payload, ok := readWorkerRequest(w, request)
	if !ok {
		return
	}
	requestService, ok := handler.service.(operatorWorkerRequestService)
	if !ok {
		writeAPIError(w, http.StatusInternalServerError, "internal_error", http.StatusText(http.StatusInternalServerError))
		return
	}
	response, err := requestService.CompleteWorkerRequest(request.Context(), requestID, payload.WorkerID, user.ID)
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) eligibleWorkers(w http.ResponseWriter, request *http.Request) {
	if _, ok := handler.operatorUser(w, request); !ok {
		return
	}
	response, err := handler.service.EligibleWorkers(
		request.Context(),
		request.PathValue("track_id"),
		request.URL.Query().Get("name"),
		request.URL.Query().Get("expert_group_id"),
	)
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) setResponsibleWorker(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.operatorUser(w, request)
	if !ok {
		return
	}
	payload, ok := readWorkerRequest(w, request)
	if !ok {
		return
	}
	response, err := handler.service.SetResponsibleWorker(request.Context(), request.PathValue("track_id"), payload.WorkerID, user.ID)
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) addWorker(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.operatorUser(w, request)
	if !ok {
		return
	}
	payload, ok := readWorkerRequest(w, request)
	if !ok {
		return
	}
	response, err := handler.service.AddWorker(request.Context(), request.PathValue("track_id"), payload.WorkerID, user.ID)
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, response)
}

func (handler *OperatorHandler) removeWorker(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.operatorUser(w, request)
	if !ok {
		return
	}
	workerID, err := strconv.ParseInt(request.PathValue("worker_id"), 10, 64)
	if err != nil || workerID <= 0 {
		writeAPIError(w, http.StatusBadRequest, "invalid_worker_id", "worker_id must be a positive integer")
		return
	}
	response, err := handler.service.RemoveWorker(request.Context(), request.PathValue("track_id"), workerID, user.ID)
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) rejectTicket(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.operatorUser(w, request)
	if !ok {
		return
	}
	request.Body = http.MaxBytesReader(w, request.Body, maxOperatorRequestBytes)
	var payload rejectRequest
	if err := decodeJSONRequest(request, &payload); err != nil {
		writeJSONRequestError(w, err)
		return
	}
	response, err := handler.service.Reject(request.Context(), request.PathValue("track_id"), user.ID, payload.ReasonCode, payload.Message)
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) listTickets(w http.ResponseWriter, request *http.Request) {
	if _, ok := handler.operatorUser(w, request); !ok {
		return
	}
	query := request.URL.Query()
	response, err := handler.service.ListTickets(request.Context(), service.TicketListRequest{
		Queue: query.Get("queue"), Search: query.Get("search"), Priority: query.Get("priority"),
		Status: query.Get("status"), CategoryID: query.Get("category_id"),
		ApplicantType: query.Get("applicant_type"), Page: query.Get("page"), Limit: query.Get("limit"),
	})
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) analytics(w http.ResponseWriter, request *http.Request) {
	if _, ok := handler.operatorUser(w, request); !ok {
		return
	}
	query := request.URL.Query()
	response, err := handler.service.Analytics(request.Context(), query.Get("date_from"), query.Get("date_to"))
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *OperatorHandler) report(w http.ResponseWriter, request *http.Request) {
	if _, ok := handler.operatorUser(w, request); !ok {
		return
	}
	query := request.URL.Query()
	report, err := handler.service.Report(request.Context(), query.Get("date_from"), query.Get("date_to"), query.Get("format"))
	if err != nil {
		writeOperatorError(w, err)
		return
	}
	w.Header().Set("Content-Type", report.ContentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": report.Name}))
	w.Header().Set("Content-Length", strconv.Itoa(len(report.Data)))
	w.Header().Set("Cache-Control", "private, no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(report.Data)
}

func (handler *OperatorHandler) operatorUser(w http.ResponseWriter, request *http.Request) (service.AuthUser, bool) {
	user, err := handler.auth.Authenticate(request.Context(), request.Header.Get("Authorization"))
	if err != nil {
		writeAuthError(w, err)
		return service.AuthUser{}, false
	}
	if err := service.RequireOperator(user); err != nil {
		writeAuthError(w, err)
		return service.AuthUser{}, false
	}
	return user, true
}

func readWorkerRequest(w http.ResponseWriter, request *http.Request) (workerRequest, bool) {
	request.Body = http.MaxBytesReader(w, request.Body, maxOperatorRequestBytes)
	var payload workerRequest
	if err := decodeJSONRequest(request, &payload); err != nil {
		writeJSONRequestError(w, err)
		return workerRequest{}, false
	}
	if payload.WorkerID <= 0 {
		writeAPIError(w, http.StatusBadRequest, "invalid_worker_id", "worker_id must be a positive integer")
		return workerRequest{}, false
	}
	return payload, true
}

func decodeJSONRequest(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	return ensureJSONEnd(decoder)
}

func writeOperatorError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTrackID), errors.Is(err, service.ErrTicketNotFound):
		writeAPIError(w, http.StatusNotFound, "ticket_not_found", "ticket not found")
	case errors.Is(err, service.ErrInvalidWorkerID), errors.Is(err, service.ErrInvalidFilter),
		errors.Is(err, service.ErrInvalidQueue), errors.Is(err, service.ErrInvalidDateRange),
		errors.Is(err, service.ErrInvalidReportFormat), errors.Is(err, service.ErrInvalidReasonCode),
		errors.Is(err, service.ErrRejectionMessage):
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, service.ErrWorkerNotFound):
		writeAPIError(w, http.StatusNotFound, "worker_not_found", "worker not found")
	case errors.Is(err, service.ErrWorkerRequestNotFound):
		writeAPIError(w, http.StatusNotFound, "worker_request_not_found", "worker request not found")
	case errors.Is(err, service.ErrWorkerUnavailable):
		writeAPIError(w, http.StatusConflict, "worker_unavailable", "worker has reached the active ticket limit")
	case errors.Is(err, service.ErrWorkerAlreadyAssigned):
		writeAPIError(w, http.StatusConflict, "worker_already_assigned", "worker is already assigned to the ticket")
	case errors.Is(err, service.ErrWorkerNotAssigned):
		writeAPIError(w, http.StatusNotFound, "worker_not_assigned", "worker is not assigned to the ticket")
	case errors.Is(err, service.ErrResponsibleWorker):
		writeAPIError(w, http.StatusConflict, "responsible_worker", "responsible worker must be replaced, not removed")
	case errors.Is(err, service.ErrResponsibleRequired):
		writeAPIError(w, http.StatusConflict, "responsible_worker_required", "assign a responsible worker before adding co-workers")
	case errors.Is(err, service.ErrInvalidTransition):
		writeAPIError(w, http.StatusConflict, "invalid_status_transition", "operation is not allowed for the current ticket status")
	default:
		writeAPIError(w, http.StatusInternalServerError, "internal_error", http.StatusText(http.StatusInternalServerError))
	}
}
