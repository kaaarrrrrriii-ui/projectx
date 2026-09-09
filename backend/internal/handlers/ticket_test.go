package handlers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/german/backend/internal/service"
)

type stubTicketService struct {
	statusResponse   service.TicketStatusResponse
	statusErr        error
	chatResponse     service.TicketChatResponse
	chatErr          error
	completeResponse service.CompleteTicketResponse
	completeErr      error
	returnResponse   service.ReturnTicketResponse
	returnErr        error
	returnReason     string
}

func (stub *stubTicketService) Status(context.Context, string) (service.TicketStatusResponse, error) {
	return stub.statusResponse, stub.statusErr
}

func (stub *stubTicketService) Chat(context.Context, string) (service.TicketChatResponse, error) {
	return stub.chatResponse, stub.chatErr
}

func (stub *stubTicketService) Complete(context.Context, string) (service.CompleteTicketResponse, error) {
	return stub.completeResponse, stub.completeErr
}

func (stub *stubTicketService) Return(_ context.Context, _, reason string) (service.ReturnTicketResponse, error) {
	stub.returnReason = reason
	return stub.returnResponse, stub.returnErr
}

type stubMessageService struct {
	response service.CreateMessageResponse
	err      error
	text     string
	uploads  []service.UploadedAttachment
}

func (stub *stubMessageService) AddApplicantMessage(
	_ context.Context,
	_ string,
	text string,
	uploads []service.UploadedAttachment,
) (service.CreateMessageResponse, error) {
	stub.text = text
	stub.uploads = uploads
	return stub.response, stub.err
}

type stubAttachmentService struct {
	file service.AttachmentFile
	err  error
}

func (stub *stubAttachmentService) Open(context.Context, string, int64) (service.AttachmentFile, error) {
	return stub.file, stub.err
}

func newTicketTestMux(t *testing.T, tickets *stubTicketService, messages *stubMessageService, attachments *stubAttachmentService) *http.ServeMux {
	t.Helper()
	handler, err := NewTicketHandler(tickets, messages, attachments)
	if err != nil {
		t.Fatalf("NewTicketHandler() error = %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux
}

func TestTicketStatusHandler(t *testing.T) {
	t.Parallel()

	tickets := &stubTicketService{statusResponse: service.TicketStatusResponse{
		TrackID:   "ОТК-ABCD-2345",
		Status:    "answer_ready",
		CreatedAt: time.Date(2026, 9, 9, 9, 48, 0, 0, time.UTC),
	}}
	mux := newTicketTestMux(t, tickets, &stubMessageService{}, &stubAttachmentService{})
	request := httptest.NewRequest(http.MethodGet, "/api/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/status", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"status":"answer_ready"`) {
		t.Fatalf("body = %s", response.Body.String())
	}
}

func TestTicketStatusHandlerLimitsChecksPerAddress(t *testing.T) {
	t.Parallel()

	tickets := &stubTicketService{statusResponse: service.TicketStatusResponse{Status: "new"}}
	mux := newTicketTestMux(t, tickets, &stubMessageService{}, &stubAttachmentService{})
	for attempt := 1; attempt <= 6; attempt++ {
		request := httptest.NewRequest(http.MethodGet, "/api/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/status", nil)
		request.RemoteAddr = "203.0.113.5:12345"
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)

		want := http.StatusOK
		if attempt == 6 {
			want = http.StatusTooManyRequests
		}
		if response.Code != want {
			t.Fatalf("attempt %d status = %d, want %d", attempt, response.Code, want)
		}
		if attempt == 6 && response.Header().Get("Retry-After") != "60" {
			t.Fatalf("Retry-After = %q", response.Header().Get("Retry-After"))
		}
	}
}

func TestTicketReturnHandlerAcceptsMissingReason(t *testing.T) {
	t.Parallel()

	tickets := &stubTicketService{returnResponse: service.ReturnTicketResponse{Status: "returned", ReturnCount: 1}}
	mux := newTicketTestMux(t, tickets, &stubMessageService{}, &stubAttachmentService{})
	request := httptest.NewRequest(http.MethodPost, "/api/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/return", strings.NewReader(`{}`))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK || tickets.returnReason != "" {
		t.Fatalf("status = %d, reason = %q, body = %s", response.Code, tickets.returnReason, response.Body.String())
	}
}

func TestTicketReturnHandlerMapsLimitToConflict(t *testing.T) {
	t.Parallel()

	tickets := &stubTicketService{returnErr: service.ErrReturnLimitReached}
	mux := newTicketTestMux(t, tickets, &stubMessageService{}, &stubAttachmentService{})
	request := httptest.NewRequest(http.MethodPost, "/api/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/return", strings.NewReader(`{"reason":"Нет"}`))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusConflict || !strings.Contains(response.Body.String(), `"code":"return_limit_reached"`) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestPostMessageHandlerReadsMultipartFiles(t *testing.T) {
	t.Parallel()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("text", "Дополнительная информация"); err != nil {
		t.Fatalf("WriteField() error = %v", err)
	}
	part, err := writer.CreateFormFile("attachments[]", "proof.png")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := part.Write([]byte("image bytes")); err != nil {
		t.Fatalf("write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	messages := &stubMessageService{response: service.CreateMessageResponse{
		Message:      service.MessageResponse{ID: 15, Type: "applicant", Attachments: []service.AttachmentResponse{}},
		TicketStatus: "in_progress",
	}}
	mux := newTicketTestMux(t, &stubTicketService{}, messages, &stubAttachmentService{})
	request := httptest.NewRequest(http.MethodPost, "/api/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/messages", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if messages.text != "Дополнительная информация" || len(messages.uploads) != 1 || messages.uploads[0].Size != int64(len("image bytes")) {
		t.Fatalf("message input text = %q, uploads = %+v", messages.text, messages.uploads)
	}
}

func TestAttachmentHandlerUsesInlinePrivateResponse(t *testing.T) {
	t.Parallel()

	attachments := &stubAttachmentService{file: service.AttachmentFile{
		Reader:    io.NopCloser(strings.NewReader("stored image")),
		Name:      "image-abcdef123456.png",
		MIMEType:  "image/png",
		SizeBytes: int64(len("stored image")),
	}}
	mux := newTicketTestMux(t, &stubTicketService{}, &stubMessageService{}, attachments)
	request := httptest.NewRequest(http.MethodGet, "/api/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/attachments/42", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK || response.Body.String() != "stored image" {
		t.Fatalf("status = %d, body = %q", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Disposition") != `inline; filename=image-abcdef123456.png` {
		t.Fatalf("Content-Disposition = %q", response.Header().Get("Content-Disposition"))
	}
	if response.Header().Get("Cache-Control") != "private, no-store" || response.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("security headers = %#v", response.Header())
	}
}

func TestTicketHandlerHidesInvalidTrackAsNotFound(t *testing.T) {
	t.Parallel()

	tickets := &stubTicketService{statusErr: service.ErrInvalidTrackID}
	mux := newTicketTestMux(t, tickets, &stubMessageService{}, &stubAttachmentService{})
	request := httptest.NewRequest(http.MethodGet, "/api/tickets/bad/status", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound || !strings.Contains(response.Body.String(), `"code":"ticket_not_found"`) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestTicketHandlerMapsUnexpectedErrorsToInternalServerError(t *testing.T) {
	t.Parallel()

	tickets := &stubTicketService{completeErr: errors.New("database unavailable")}
	mux := newTicketTestMux(t, tickets, &stubMessageService{}, &stubAttachmentService{})
	request := httptest.NewRequest(http.MethodPost, "/api/tickets/%D0%9E%D0%A2%D0%9A-ABCD-2345/complete", nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), `"code":"internal_error"`) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}
