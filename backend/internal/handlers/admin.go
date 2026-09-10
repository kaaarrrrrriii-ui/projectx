package handlers

import (
	"context"
	"errors"
	"mime"
	"net/http"
	"strconv"

	"example.com/german/backend/internal/service"
)

const maxAdminRequestBytes = 64 * 1024

type adminService interface {
	Dashboard(context.Context) (service.AdminDashboardResponse, error)
	ListEmployees(context.Context, string, string, string, string) (service.AdminEmployeePageResponse, error)
	CreateEmployee(context.Context, service.CreateAdminEmployeeRequest) (service.AdminEmployeeResponse, error)
	UpdateEmployee(context.Context, int64, service.UpdateAdminEmployeeRequest) (service.AdminEmployeeResponse, error)
	DeactivateEmployee(context.Context, int64, int64) (service.AdminEmployeeResponse, error)
	RestoreEmployee(context.Context, int64, service.RestoreAdminEmployeeRequest) (service.AdminEmployeeResponse, error)
	ListGroups(context.Context, string, string, string) (service.AdminGroupPageResponse, error)
	CreateGroup(context.Context, string) (service.AdminGroupResponse, error)
	UpdateGroup(context.Context, int64, string) (service.AdminGroupResponse, error)
	DeleteGroup(context.Context, int64) error
	ListCategories(context.Context, string, string, string) (service.AdminCategoryPageResponse, error)
	Category(context.Context, int64) (service.AdminCategoryResponse, error)
	CreateCategory(context.Context, string) (service.AdminCategoryResponse, error)
	UpdateCategory(context.Context, int64, string) (service.AdminCategoryResponse, error)
	DeleteCategory(context.Context, int64) error
	ReplaceCategoryGroups(context.Context, int64, []int64) (service.AdminCategoryResponse, error)
	CreateQuestion(context.Context, int64, string, []string) (service.AdminQuestionResponse, error)
	UpdateQuestion(context.Context, int64, string, []string) (service.AdminQuestionResponse, error)
	DeleteQuestion(context.Context, int64) error
	ListTickets(context.Context, service.AdminTicketListRequest) (service.AdminTicketPageResponse, error)
	Ticket(context.Context, string) (service.AdminTicketResponse, error)
	ChangePriority(context.Context, string, string, int64) (service.AdminTicketResponse, error)
	ChangeStatus(context.Context, string, string, int64) (service.AdminTicketResponse, error)
	SetResponsible(context.Context, string, int64, int64) (service.AdminTicketResponse, error)
	Events(context.Context, string, string, string) (service.AdminEventPageResponse, error)
	Analytics(context.Context, string, string) (service.AdminAnalyticsResponse, error)
	Report(context.Context, string, string, string) (service.GeneratedReport, error)
}

type AdminHandler struct {
	service adminService
	auth    authService
}

func NewAdminHandler(adminService adminService, auth authService) (*AdminHandler, error) {
	if adminService == nil || auth == nil {
		return nil, errors.New("admin and auth services are required")
	}
	return &AdminHandler{service: adminService, auth: auth}, nil
}

func (handler *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/admin/me", handler.me)
	mux.HandleFunc("GET /api/admin/dashboard", handler.dashboard)
	mux.HandleFunc("GET /api/admin/employees", handler.listEmployees)
	mux.HandleFunc("POST /api/admin/employees", handler.createEmployee)
	mux.HandleFunc("PATCH /api/admin/employees/{employee_id}", handler.updateEmployee)
	mux.HandleFunc("POST /api/admin/employees/{employee_id}/deactivate", handler.deactivateEmployee)
	mux.HandleFunc("POST /api/admin/employees/{employee_id}/restore", handler.restoreEmployee)
	mux.HandleFunc("GET /api/admin/expert-groups", handler.listGroups)
	mux.HandleFunc("POST /api/admin/expert-groups", handler.createGroup)
	mux.HandleFunc("PATCH /api/admin/expert-groups/{group_id}", handler.updateGroup)
	mux.HandleFunc("DELETE /api/admin/expert-groups/{group_id}", handler.deleteGroup)
	mux.HandleFunc("GET /api/admin/categories", handler.listCategories)
	mux.HandleFunc("POST /api/admin/categories", handler.createCategory)
	mux.HandleFunc("GET /api/admin/categories/{category_id}", handler.category)
	mux.HandleFunc("PATCH /api/admin/categories/{category_id}", handler.updateCategory)
	mux.HandleFunc("DELETE /api/admin/categories/{category_id}", handler.deleteCategory)
	mux.HandleFunc("PUT /api/admin/categories/{category_id}/expert-groups", handler.replaceCategoryGroups)
	mux.HandleFunc("POST /api/admin/categories/{category_id}/questions", handler.createQuestion)
	mux.HandleFunc("PATCH /api/admin/questions/{question_id}", handler.updateQuestion)
	mux.HandleFunc("DELETE /api/admin/questions/{question_id}", handler.deleteQuestion)
	mux.HandleFunc("GET /api/admin/tickets", handler.listTickets)
	mux.HandleFunc("GET /api/admin/tickets/{track_id}", handler.ticket)
	mux.HandleFunc("PATCH /api/admin/tickets/{track_id}/priority", handler.changePriority)
	mux.HandleFunc("PATCH /api/admin/tickets/{track_id}/status", handler.changeStatus)
	mux.HandleFunc("PUT /api/admin/tickets/{track_id}/responsible-worker", handler.setResponsible)
	mux.HandleFunc("GET /api/admin/tickets/{track_id}/events", handler.events)
	mux.HandleFunc("GET /api/admin/analytics", handler.analytics)
	mux.HandleFunc("GET /api/admin/reports", handler.report)
}

func (handler *AdminHandler) adminUser(w http.ResponseWriter, request *http.Request) (service.AuthUser, bool) {
	user, err := handler.auth.Authenticate(request.Context(), request.Header.Get("Authorization"))
	if err != nil {
		writeAuthError(w, err)
		return service.AuthUser{}, false
	}
	if err = service.RequireAdmin(user); err != nil {
		writeAuthError(w, err)
		return service.AuthUser{}, false
	}
	return user, true
}
func adminID(request *http.Request, name string) (int64, error) {
	id, err := strconv.ParseInt(request.PathValue(name), 10, 64)
	if err != nil || id <= 0 {
		return 0, service.ErrAdminResourceNotFound
	}
	return id, nil
}
func (handler *AdminHandler) decode(w http.ResponseWriter, request *http.Request, destination any) bool {
	if err := decodeJSONBody(w, request, maxAdminRequestBytes, destination); err != nil {
		writeJSONRequestError(w, err)
		return false
	}
	return true
}

func (handler *AdminHandler) me(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.adminUser(w, r)
	if ok {
		writeJSON(w, http.StatusOK, map[string]service.AuthUser{"user": user})
	}
}
func (handler *AdminHandler) dashboard(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	x, err := handler.service.Dashboard(r.Context())
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) listEmployees(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	q := r.URL.Query()
	x, err := handler.service.ListEmployees(r.Context(), q.Get("role"), q.Get("search"), q.Get("page"), q.Get("limit"))
	handler.write(w, http.StatusOK, x, err)
}

type adminEmployeePayload struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	FullName      string `json:"full_name"`
	Role          string `json:"role"`
	ExpertGroupID int64  `json:"expert_group_id"`
	MaxTickets    int    `json:"max_tickets"`
}

func (handler *AdminHandler) createEmployee(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	var p adminEmployeePayload
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.CreateEmployee(r.Context(), service.CreateAdminEmployeeRequest{Username: p.Username, Password: p.Password, FullName: p.FullName, Role: p.Role, ExpertGroupID: p.ExpertGroupID, MaxTickets: p.MaxTickets})
	handler.write(w, http.StatusCreated, x, err)
}

type adminEmployeePatch struct {
	Username      *string `json:"username"`
	Password      *string `json:"password"`
	FullName      *string `json:"full_name"`
	Role          *string `json:"role"`
	ExpertGroupID *int64  `json:"expert_group_id"`
	MaxTickets    *int    `json:"max_tickets"`
}

func (handler *AdminHandler) updateEmployee(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "employee_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	var p adminEmployeePatch
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.UpdateEmployee(r.Context(), id, service.UpdateAdminEmployeeRequest{Username: p.Username, Password: p.Password, FullName: p.FullName, Role: p.Role, ExpertGroupID: p.ExpertGroupID, MaxTickets: p.MaxTickets})
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) deactivateEmployee(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.adminUser(w, r)
	if !ok {
		return
	}
	id, err := adminID(r, "employee_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	x, err := handler.service.DeactivateEmployee(r.Context(), id, user.ID)
	handler.write(w, http.StatusOK, x, err)
}

type adminRestorePayload struct {
	Role          string `json:"role"`
	ExpertGroupID int64  `json:"expert_group_id"`
	MaxTickets    int    `json:"max_tickets"`
	Password      string `json:"password,omitempty"`
}

func (handler *AdminHandler) restoreEmployee(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "employee_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	var p adminRestorePayload
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.RestoreEmployee(r.Context(), id, service.RestoreAdminEmployeeRequest{Role: p.Role, ExpertGroupID: p.ExpertGroupID, MaxTickets: p.MaxTickets, Password: p.Password})
	handler.write(w, http.StatusOK, x, err)
}

type adminNamePayload struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (handler *AdminHandler) listGroups(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	q := r.URL.Query()
	x, err := handler.service.ListGroups(r.Context(), q.Get("search"), q.Get("page"), q.Get("limit"))
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) createGroup(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	var p struct {
		Title string `json:"title"`
	}
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.CreateGroup(r.Context(), p.Title)
	handler.write(w, http.StatusCreated, x, err)
}
func (handler *AdminHandler) updateGroup(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "group_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	var p struct {
		Title string `json:"title"`
	}
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.UpdateGroup(r.Context(), id, p.Title)
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) deleteGroup(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "group_id")
	if err == nil {
		err = handler.service.DeleteGroup(r.Context(), id)
	}
	if err != nil {
		writeAdminError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *AdminHandler) listCategories(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	q := r.URL.Query()
	x, err := handler.service.ListCategories(r.Context(), q.Get("search"), q.Get("page"), q.Get("limit"))
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) category(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "category_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	x, err := handler.service.Category(r.Context(), id)
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	var p struct {
		Name string `json:"name"`
	}
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.CreateCategory(r.Context(), p.Name)
	handler.write(w, http.StatusCreated, x, err)
}
func (handler *AdminHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "category_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	var p struct {
		Name string `json:"name"`
	}
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.UpdateCategory(r.Context(), id, p.Name)
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "category_id")
	if err == nil {
		err = handler.service.DeleteCategory(r.Context(), id)
	}
	if err != nil {
		writeAdminError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (handler *AdminHandler) replaceCategoryGroups(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "category_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	var p struct {
		ExpertGroupIDs []int64 `json:"expert_group_ids"`
	}
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.ReplaceCategoryGroups(r.Context(), id, p.ExpertGroupIDs)
	handler.write(w, http.StatusOK, x, err)
}

type adminQuestionPayload struct {
	Text    string   `json:"text"`
	Answers []string `json:"answers"`
}

func (handler *AdminHandler) createQuestion(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "category_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	var p adminQuestionPayload
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.CreateQuestion(r.Context(), id, p.Text, p.Answers)
	handler.write(w, http.StatusCreated, x, err)
}
func (handler *AdminHandler) updateQuestion(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "question_id")
	if err != nil {
		writeAdminError(w, err)
		return
	}
	var p adminQuestionPayload
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.UpdateQuestion(r.Context(), id, p.Text, p.Answers)
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) deleteQuestion(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	id, err := adminID(r, "question_id")
	if err == nil {
		err = handler.service.DeleteQuestion(r.Context(), id)
	}
	if err != nil {
		writeAdminError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (handler *AdminHandler) listTickets(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	q := r.URL.Query()
	x, err := handler.service.ListTickets(r.Context(), service.AdminTicketListRequest{Search: q.Get("search"), Status: q.Get("status"), Priority: q.Get("priority"), CategoryID: q.Get("category_id"), ApplicantType: q.Get("applicant_type"), ResponsibleWorkerID: q.Get("responsible_worker_id"), RoutingIssue: q.Get("routing_issue"), Page: q.Get("page"), Limit: q.Get("limit")})
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) ticket(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	x, err := handler.service.Ticket(r.Context(), r.PathValue("track_id"))
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) changePriority(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.adminUser(w, r)
	if !ok {
		return
	}
	var p struct {
		Priority string `json:"priority"`
	}
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.ChangePriority(r.Context(), r.PathValue("track_id"), p.Priority, user.ID)
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) changeStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.adminUser(w, r)
	if !ok {
		return
	}
	var p struct {
		Status string `json:"status"`
	}
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.ChangeStatus(r.Context(), r.PathValue("track_id"), p.Status, user.ID)
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) setResponsible(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.adminUser(w, r)
	if !ok {
		return
	}
	var p workerRequest
	if !handler.decode(w, r, &p) {
		return
	}
	x, err := handler.service.SetResponsible(r.Context(), r.PathValue("track_id"), p.WorkerID, user.ID)
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) events(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	q := r.URL.Query()
	x, err := handler.service.Events(r.Context(), r.PathValue("track_id"), q.Get("page"), q.Get("limit"))
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) analytics(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	q := r.URL.Query()
	x, err := handler.service.Analytics(r.Context(), q.Get("date_from"), q.Get("date_to"))
	handler.write(w, http.StatusOK, x, err)
}
func (handler *AdminHandler) report(w http.ResponseWriter, r *http.Request) {
	if _, ok := handler.adminUser(w, r); !ok {
		return
	}
	q := r.URL.Query()
	report, err := handler.service.Report(r.Context(), q.Get("date_from"), q.Get("date_to"), q.Get("format"))
	if err != nil {
		writeAdminError(w, err)
		return
	}
	w.Header().Set("Content-Type", report.ContentType)
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": report.Name}))
	w.Header().Set("Content-Length", strconv.Itoa(len(report.Data)))
	w.Header().Set("Cache-Control", "private, no-store")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(report.Data)
}

func (handler *AdminHandler) write(w http.ResponseWriter, status int, value any, err error) {
	if err != nil {
		writeAdminError(w, err)
		return
	}
	writeJSON(w, status, value)
}
func writeAdminError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTrackID), errors.Is(err, service.ErrTicketNotFound):
		writeAPIError(w, http.StatusNotFound, "ticket_not_found", "ticket not found")
	case errors.Is(err, service.ErrAdminResourceNotFound):
		writeAPIError(w, http.StatusNotFound, "resource_not_found", "resource not found")
	case errors.Is(err, service.ErrInvalidAdminInput), errors.Is(err, service.ErrInvalidDateRange), errors.Is(err, service.ErrInvalidReportFormat):
		writeAPIError(w, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, service.ErrUsernameExists):
		writeAPIError(w, http.StatusConflict, "username_already_exists", err.Error())
	case errors.Is(err, service.ErrNameExists):
		writeAPIError(w, http.StatusConflict, "name_already_exists", err.Error())
	case errors.Is(err, service.ErrEmployeeHasTickets):
		writeAPIError(w, http.StatusConflict, "employee_has_active_tickets", err.Error())
	case errors.Is(err, service.ErrQuestionInUse):
		writeAPIError(w, http.StatusConflict, "question_in_use", err.Error())
	case errors.Is(err, service.ErrAdminResourceInUse):
		writeAPIError(w, http.StatusConflict, "resource_in_use", err.Error())
	case errors.Is(err, service.ErrCannotManageAdmin), errors.Is(err, service.ErrCannotDeactivateSelf), errors.Is(err, service.ErrEmployeeNotInactive):
		writeAPIError(w, http.StatusConflict, "operation_not_allowed", err.Error())
	case errors.Is(err, service.ErrWorkerNotFound):
		writeAPIError(w, http.StatusNotFound, "worker_not_found", "worker not found")
	case errors.Is(err, service.ErrWorkerUnavailable), errors.Is(err, service.ErrResponsibleRequired):
		writeAPIError(w, http.StatusConflict, "ticket_conflict", err.Error())
	default:
		writeAPIError(w, http.StatusInternalServerError, "internal_error", http.StatusText(http.StatusInternalServerError))
	}
}
