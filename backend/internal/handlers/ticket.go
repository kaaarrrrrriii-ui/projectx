package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"strconv"
	"time"

	"example.com/german/backend/internal/service"
)

const (
	maxMessageRequestBytes = service.DefaultMaxAttachmentFiles*service.DefaultMaxAttachmentBytes + 1024*1024
	maxReturnRequestBytes  = 64 * 1024
)

type ticketService interface {
	Status(context.Context, string) (service.TicketStatusResponse, error)
	Chat(context.Context, string) (service.TicketChatResponse, error)
	Complete(context.Context, string) (service.CompleteTicketResponse, error)
	Return(context.Context, string, string) (service.ReturnTicketResponse, error)
}

type ticketMessageService interface {
	AddApplicantMessage(context.Context, string, string, []service.UploadedAttachment) (service.CreateMessageResponse, error)
}

type ticketAttachmentService interface {
	Open(context.Context, string, int64) (service.AttachmentFile, error)
}

type TicketHandler struct {
	tickets     ticketService
	messages    ticketMessageService
	attachments ticketAttachmentService
	statusLimit *slidingWindowLimiter
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type returnTicketRequest struct {
	Reason string `json:"reason"`
}

func NewTicketHandler(
	tickets ticketService,
	messages ticketMessageService,
	attachments ticketAttachmentService,
) (*TicketHandler, error) {
	if tickets == nil || messages == nil || attachments == nil {
		return nil, errors.New("all ticket handler services are required")
	}
	return &TicketHandler{
		tickets:     tickets,
		messages:    messages,
		attachments: attachments,
		statusLimit: newSlidingWindowLimiter(5, time.Minute),
	}, nil
}

func (handler *TicketHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tickets/{track_id}/status", handler.getStatus)
	mux.HandleFunc("GET /api/tickets/{track_id}/chat", handler.getChat)
	mux.HandleFunc("POST /api/tickets/{track_id}/messages", handler.postMessage)
	mux.HandleFunc("GET /api/tickets/{track_id}/attachments/{attachment_id}", handler.getAttachment)
	mux.HandleFunc("POST /api/tickets/{track_id}/complete", handler.completeTicket)
	mux.HandleFunc("POST /api/tickets/{track_id}/return", handler.returnTicket)
}

func (handler *TicketHandler) getStatus(w http.ResponseWriter, request *http.Request) {
	if !handler.statusLimit.Allow(clientAddress(request)) {
		w.Header().Set("Retry-After", "60")
		writeAPIError(w, http.StatusTooManyRequests, "rate_limit_reached", "too many status checks")
		return
	}
	response, err := handler.tickets.Status(request.Context(), request.PathValue("track_id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func clientAddress(request *http.Request) string {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		return host
	}
	return request.RemoteAddr
}

func (handler *TicketHandler) getChat(w http.ResponseWriter, request *http.Request) {
	response, err := handler.tickets.Chat(request.Context(), request.PathValue("track_id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *TicketHandler) postMessage(w http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(w, request.Body, maxMessageRequestBytes)
	if err := request.ParseMultipartForm(1024 * 1024); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeAPIError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large")
			return
		}
		writeAPIError(w, http.StatusBadRequest, "invalid_multipart_form", "invalid multipart form")
		return
	}
	if request.MultipartForm != nil {
		defer request.MultipartForm.RemoveAll()
	}

	fileHeaders := request.MultipartForm.File["attachments[]"]
	fileHeaders = append(fileHeaders, request.MultipartForm.File["attachments"]...)
	uploads := make([]service.UploadedAttachment, 0, len(fileHeaders))
	opened := make([]io.ReadCloser, 0, len(fileHeaders))
	defer func() {
		for _, file := range opened {
			_ = file.Close()
		}
	}()
	for _, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			writeAPIError(w, http.StatusBadRequest, "invalid_attachment", "cannot read attachment")
			return
		}
		opened = append(opened, file)
		uploads = append(uploads, service.UploadedAttachment{Source: file, Size: header.Size})
	}

	response, err := handler.messages.AddApplicantMessage(
		request.Context(),
		request.PathValue("track_id"),
		request.FormValue("text"),
		uploads,
	)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, response)
}

func (handler *TicketHandler) getAttachment(w http.ResponseWriter, request *http.Request) {
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

	disposition := mime.FormatMediaType("inline", map[string]string{"filename": file.Name})
	w.Header().Set("Content-Type", file.MIMEType)
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Length", strconv.FormatInt(file.SizeBytes, 10))
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, file.Reader)
}

func (handler *TicketHandler) completeTicket(w http.ResponseWriter, request *http.Request) {
	response, err := handler.tickets.Complete(request.Context(), request.PathValue("track_id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func (handler *TicketHandler) returnTicket(w http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(w, request.Body, maxReturnRequestBytes)
	payload := returnTicketRequest{}
	if request.Body != http.NoBody {
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&payload); err != nil && !errors.Is(err, io.EOF) {
			writeJSONRequestError(w, err)
			return
		}
		if err := ensureJSONEnd(decoder); err != nil {
			writeJSONRequestError(w, err)
			return
		}
	}

	response, err := handler.tickets.Return(request.Context(), request.PathValue("track_id"), payload.Reason)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSONRequestError(w http.ResponseWriter, err error) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		writeAPIError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large")
		return
	}
	writeAPIError(w, http.StatusBadRequest, "invalid_json", "invalid JSON request body")
}

func ensureJSONEnd(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTrackID), errors.Is(err, service.ErrTicketNotFound):
		writeAPIError(w, http.StatusNotFound, "ticket_not_found", "ticket not found")
	case errors.Is(err, service.ErrAttachmentNotFound):
		writeAPIError(w, http.StatusNotFound, "attachment_not_found", "attachment not found")
	case errors.Is(err, service.ErrMessageTextRequired):
		writeAPIError(w, http.StatusBadRequest, "message_text_required", "message text is required")
	case errors.Is(err, service.ErrTooManyAttachments), errors.Is(err, service.ErrAttachmentLimitReached):
		writeAPIError(w, http.StatusConflict, "attachment_limit_reached", "attachment limit reached")
	case errors.Is(err, service.ErrAttachmentTooLarge), errors.Is(err, service.ErrImageInputTooLarge):
		writeAPIError(w, http.StatusRequestEntityTooLarge, "attachment_too_large", "attachment is too large")
	case errors.Is(err, service.ErrUnsupportedImage):
		writeAPIError(w, http.StatusUnsupportedMediaType, "unsupported_attachment", "only JPEG and PNG images are supported")
	case errors.Is(err, service.ErrInvalidImage), errors.Is(err, service.ErrImageTooManyPixels), errors.Is(err, service.ErrImageOutputTooLarge):
		writeAPIError(w, http.StatusUnprocessableEntity, "invalid_attachment", "attachment is not a valid supported image")
	case errors.Is(err, service.ErrReturnLimitReached):
		writeAPIError(w, http.StatusConflict, "return_limit_reached", "ticket return limit reached")
	case errors.Is(err, service.ErrInvalidTransition):
		writeAPIError(w, http.StatusConflict, "invalid_status_transition", "operation is not allowed for the current ticket status")
	default:
		writeAPIError(w, http.StatusInternalServerError, "internal_error", http.StatusText(http.StatusInternalServerError))
	}
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: apiError{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
