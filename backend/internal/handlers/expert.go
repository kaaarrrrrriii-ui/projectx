package handlers

import (
	"context"
	"errors"
	"io"
	"mime"
	"net/http"
	"strconv"

	"example.com/german/backend/internal/service"
)

const maxExpertRequestBytes = 64 * 1024

type expertService interface {
	Profile(context.Context, int64) (service.ExpertProfileResponse, error)
	Dashboard(context.Context, int64) (service.ExpertDashboardResponse, error)
	ListTickets(context.Context, int64, service.ExpertTicketListRequest) (service.ExpertTicketPageResponse, error)
	OpenTicket(context.Context, string, int64) (service.ExpertOpenResponse, error)
	Ticket(context.Context, string, int64) (service.ExpertTicketResponse, error)
	AddNote(context.Context, string, service.AuthUser, string) (service.CreateExpertNoteResponse, error)
	AddAnswer(context.Context, string, int64, string) (service.CreateExpertAnswerResponse, error)
	CreateWorkerRequest(context.Context, string, int64, string, string) (service.CreateWorkerRequestResponse, error)
	ListWorkerRequests(context.Context, int64, string, string, string) (service.WorkerRequestPageResponse, error)
	CanAccessTicket(context.Context, string, int64) error
	Analytics(context.Context, int64, string, string) (service.AnalyticsResponse, error)
	Report(context.Context, int64, string, string, string) (service.GeneratedReport, error)
}

type ExpertHandler struct {
	service     expertService
	auth        authService
	attachments ticketAttachmentService
}

type expertTextRequest struct {
	Text string `json:"text"`
}

type createWorkerRequest struct {
	RequestType string `json:"request_type"`
	Reason      string `json:"reason"`
}

func NewExpertHandler(expertService expertService, auth authService, attachments ticketAttachmentService) (*ExpertHandler, error) {
	if expertService == nil || auth == nil || attachments == nil {
		return nil, errors.New("expert, auth and attachment services are required")
	}
	return &ExpertHandler{service: expertService, auth: auth, attachments: attachments}, nil
}

func (handler *ExpertHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/expert/me", handler.me)
	mux.HandleFunc("GET /api/expert/dashboard", handler.dashboard)
	mux.HandleFunc("GET /api/expert/tickets", handler.listTickets)
	mux.HandleFunc("POST /api/expert/tickets/{track_id}/open", handler.openTicket)
	mux.HandleFunc("GET /api/expert/tickets/{track_id}", handler.ticket)
	mux.HandleFunc("GET /api/expert/tickets/{track_id}/attachments/{attachment_id}", handler.attachment)
	mux.HandleFunc("POST /api/expert/tickets/{track_id}/notes", handler.addNote)
	mux.HandleFunc("POST /api/expert/tickets/{track_id}/messages", handler.addAnswer)
	mux.HandleFunc("POST /api/expert/tickets/{track_id}/requests", handler.createRequest)
	mux.HandleFunc("GET /api/expert/requests", handler.listRequests)
	mux.HandleFunc("GET /api/expert/analytics", handler.analytics)
	mux.HandleFunc("GET /api/expert/reports", handler.report)
}

func (handler *ExpertHandler) me(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	response, err := handler.service.Profile(request.Context(), user.ID)
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *ExpertHandler) dashboard(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	response, err := handler.service.Dashboard(request.Context(), user.ID)
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *ExpertHandler) listTickets(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	query := request.URL.Query()
	response, err := handler.service.ListTickets(request.Context(), user.ID, service.ExpertTicketListRequest{
		Queue: query.Get("queue"), Search: query.Get("search"), Priority: query.Get("priority"),
		CategoryID: query.Get("category_id"), ApplicantType: query.Get("applicant_type"),
		Page: query.Get("page"), Limit: query.Get("limit"),
	})
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *ExpertHandler) openTicket(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	response, err := handler.service.OpenTicket(request.Context(), request.PathValue("track_id"), user.ID)
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *ExpertHandler) ticket(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	response, err := handler.service.Ticket(request.Context(), request.PathValue("track_id"), user.ID)
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *ExpertHandler) attachment(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	if err := handler.service.CanAccessTicket(request.Context(), request.PathValue("track_id"), user.ID); err != nil {
		writeExpertError(w, err)
		return
	}
	attachmentID, err := strconv.ParseInt(request.PathValue("attachment_id"), 10, 64)
	if err != nil || attachmentID <= 0 {
		writeAPIError(w, http.StatusNotFound, "attachment_not_found", "attachment not found")
		return
	}
	file, err := handler.attachments.Open(request.Context(), request.PathValue("track_id"), attachmentID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	defer file.Reader.Close()
	w.Header().Set("Content-Type", file.MIMEType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": file.Name}))
	w.Header().Set("Content-Length", strconv.FormatInt(file.SizeBytes, 10))
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, file.Reader)
}

func (handler *ExpertHandler) addNote(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	var payload expertTextRequest
	if err := decodeJSONBody(w, request, maxExpertRequestBytes, &payload); err != nil {
		writeJSONRequestError(w, err)
		return
	}
	response, err := handler.service.AddNote(request.Context(), request.PathValue("track_id"), user, payload.Text)
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, response)
}

func (handler *ExpertHandler) addAnswer(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	var payload expertTextRequest
	if err := decodeJSONBody(w, request, maxExpertRequestBytes, &payload); err != nil {
		writeJSONRequestError(w, err)
		return
	}
	response, err := handler.service.AddAnswer(request.Context(), request.PathValue("track_id"), user.ID, payload.Text)
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, response)
}

func (handler *ExpertHandler) createRequest(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	var payload createWorkerRequest
	if err := decodeJSONBody(w, request, maxExpertRequestBytes, &payload); err != nil {
		writeJSONRequestError(w, err)
		return
	}
	response, err := handler.service.CreateWorkerRequest(request.Context(), request.PathValue("track_id"), user.ID, payload.RequestType, payload.Reason)
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, response)
}

func (handler *ExpertHandler) listRequests(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	query := request.URL.Query()
	response, err := handler.service.ListWorkerRequests(request.Context(), user.ID, query.Get("status"), query.Get("page"), query.Get("limit"))
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *ExpertHandler) analytics(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	query := request.URL.Query()
	response, err := handler.service.Analytics(request.Context(), user.ID, query.Get("date_from"), query.Get("date_to"))
	if err != nil {
		writeExpertError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *ExpertHandler) report(w http.ResponseWriter, request *http.Request) {
	user, ok := handler.expertUser(w, request)
	if !ok {
		return
	}
	query := request.URL.Query()
	report, err := handler.service.Report(request.Context(), user.ID, query.Get("date_from"), query.Get("date_to"), query.Get("format"))
	if err != nil {
		writeExpertError(w, err)
		return
	}
	w.Header().Set("Content-Type", report.ContentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": report.Name}))
	w.Header().Set("Content-Length", strconv.Itoa(len(report.Data)))
	w.Header().Set("Cache-Control", "private, no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(report.Data)
}

func (handler *ExpertHandler) expertUser(w http.ResponseWriter, request *http.Request) (service.AuthUser, bool) {
	user, err := handler.auth.Authenticate(request.Context(), request.Header.Get("Authorization"))
	if err != nil {
		writeAuthError(w, err)
		return service.AuthUser{}, false
	}
	if err := service.RequireExpert(user); err != nil {
		writeAuthError(w, err)
		return service.AuthUser{}, false
	}
	return user, true
}

func writeExpertError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTrackID), errors.Is(err, service.ErrTicketNotFound), errors.Is(err, service.ErrExpertAccessDenied):
		writeAPIError(w, http.StatusNotFound, "ticket_not_found", "ticket not found")
	case errors.Is(err, service.ErrWorkerNotFound):
		writeAPIError(w, http.StatusNotFound, "expert_not_found", "expert not found")
	case errors.Is(err, service.ErrInvalidQueue), errors.Is(err, service.ErrInvalidFilter),
		errors.Is(err, service.ErrInvalidDateRange), errors.Is(err, service.ErrInvalidReportFormat),
		errors.Is(err, service.ErrInvalidRequestType), errors.Is(err, service.ErrRequestReasonRequired),
		errors.Is(err, service.ErrNoteTextRequired), errors.Is(err, service.ErrMessageTextRequired):
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, service.ErrInvalidTransition):
		writeAPIError(w, http.StatusConflict, "invalid_status_transition", "operation is not allowed for the current ticket status")
	default:
		writeAPIError(w, http.StatusInternalServerError, "internal_error", http.StatusText(http.StatusInternalServerError))
	}
}
