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
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAdminResourceNotFound = repos.ErrAdminResourceNotFound
	ErrAdminResourceInUse    = repos.ErrAdminResourceInUse
	ErrUsernameExists        = repos.ErrUsernameExists
	ErrNameExists            = repos.ErrNameExists
	ErrEmployeeHasTickets    = repos.ErrEmployeeHasTickets
	ErrQuestionInUse         = repos.ErrQuestionInUse
	ErrInvalidAdminInput     = errors.New("invalid admin input")
	ErrCannotManageAdmin     = errors.New("admin accounts cannot be managed")
	ErrCannotDeactivateSelf  = errors.New("admin cannot deactivate own account")
	ErrEmployeeNotInactive   = errors.New("employee is not inactive")
)

type adminRepository interface {
	Dashboard(context.Context) (repos.AdminDashboardRecord, error)
	GetEmployee(context.Context, int64) (repos.AdminEmployeeRecord, error)
	ListEmployees(context.Context, string, string, int, int) (repos.AdminEmployeePageRecord, error)
	CreateEmployee(context.Context, string, string, string, string, int64, int) (repos.AdminEmployeeRecord, error)
	UpdateEmployee(context.Context, int64, string, string, string, string, int64, int) (repos.AdminEmployeeRecord, error)
	DeactivateEmployee(context.Context, int64) (repos.AdminEmployeeRecord, error)
	ListGroups(context.Context, string, int, int) (repos.AdminGroupPageRecord, error)
	CreateGroup(context.Context, string) (repos.AdminGroupRecord, error)
	UpdateGroup(context.Context, int64, string) (repos.AdminGroupRecord, error)
	DeleteGroup(context.Context, int64) error
	ListCategories(context.Context, string, int, int) (repos.AdminCategoryPageRecord, error)
	GetCategory(context.Context, int64) (repos.AdminCategoryRecord, error)
	CreateCategory(context.Context, string) (repos.AdminCategoryRecord, error)
	UpdateCategory(context.Context, int64, string) (repos.AdminCategoryRecord, error)
	DeleteCategory(context.Context, int64) error
	ReplaceCategoryGroups(context.Context, int64, []int64) (repos.AdminCategoryRecord, error)
	CreateQuestion(context.Context, int64, string, []string) (repos.AdminQuestionRecord, error)
	UpdateQuestion(context.Context, int64, string, []string) (repos.AdminQuestionRecord, error)
	DeleteQuestion(context.Context, int64) error
	ListTickets(context.Context, repos.AdminTicketFilter) (repos.AdminTicketPageRecord, error)
	GetTicket(context.Context, string) (repos.AdminTicketRecord, error)
	UpdateTicketPriority(context.Context, string, models.TicketPriority, int64, time.Time) (repos.AdminTicketRecord, error)
	UpdateTicketStatus(context.Context, string, models.TicketStatus, int64, time.Time) (repos.AdminTicketRecord, error)
	SetResponsibleWorker(context.Context, string, int64, int64, time.Time) (repos.AdminTicketRecord, error)
	ListEvents(context.Context, string, int, int) (repos.AdminEventPageRecord, error)
	Analytics(context.Context, time.Time, time.Time) (repos.AnalyticsRecord, error)
	ReportTickets(context.Context, time.Time, time.Time) ([]repos.ReportTicketRecord, error)
	EmployeeLoads(context.Context, time.Time, time.Time) ([]repos.AdminEmployeeLoadRecord, error)
}

type AdminService struct {
	repository adminRepository
	location   *time.Location
	now        func() time.Time
}

func NewAdminService(repository adminRepository) (*AdminService, error) {
	if repository == nil {
		return nil, errors.New("admin repository is required")
	}
	location, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		return nil, fmt.Errorf("load admin timezone: %w", err)
	}
	return &AdminService{repository: repository, location: location, now: time.Now}, nil
}

type AdminDashboardResponse struct {
	NewCount               int `json:"new_count"`
	ActiveCount            int `json:"active_count"`
	UrgentCount            int `json:"urgent_count"`
	ReturnedCount          int `json:"returned_count"`
	RoutingIssueCount      int `json:"routing_issue_count"`
	ExpertsAtCapacityCount int `json:"experts_at_capacity_count"`
}

type AdminEmployeeResponse struct {
	ID            int64                `json:"id"`
	Username      string               `json:"username"`
	FullName      string               `json:"full_name"`
	Role          string               `json:"role"`
	Active        bool                 `json:"active"`
	ExpertGroup   *ExpertGroupResponse `json:"expert_group"`
	AverageRating *float64             `json:"avg_rating,omitempty"`
	MaxTickets    int                  `json:"max_tickets"`
	ActiveTickets int                  `json:"active_tickets"`
	LoadPercent   float64              `json:"load_percent"`
}

type AdminEmployeePageResponse struct {
	Items []AdminEmployeeResponse `json:"items"`
	Total int                     `json:"total"`
	Page  int                     `json:"page"`
	Limit int                     `json:"limit"`
}

type CreateAdminEmployeeRequest struct {
	Username, Password, FullName, Role string
	ExpertGroupID                      int64
	MaxTickets                         int
}

type UpdateAdminEmployeeRequest struct {
	Username      *string
	Password      *string
	FullName      *string
	Role          *string
	ExpertGroupID *int64
	MaxTickets    *int
}

type RestoreAdminEmployeeRequest struct {
	Role          string
	ExpertGroupID int64
	MaxTickets    int
	Password      string
}

type AdminGroupResponse struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	EmployeeCount int    `json:"employee_count"`
	CategoryCount int    `json:"category_count"`
}

type AdminGroupPageResponse struct {
	Items []AdminGroupResponse `json:"items"`
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

type AdminAnswerResponse struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}
type AdminQuestionResponse struct {
	ID      int64                 `json:"id"`
	Text    string                `json:"text"`
	Answers []AdminAnswerResponse `json:"answers"`
}
type AdminCategoryResponse struct {
	ID           int64                   `json:"id"`
	Name         string                  `json:"name"`
	ExpertGroups []AdminGroupResponse    `json:"expert_groups"`
	Questions    []AdminQuestionResponse `json:"questions"`
}
type AdminCategoryPageResponse struct {
	Items []AdminCategoryResponse `json:"items"`
	Total int                     `json:"total"`
	Page  int                     `json:"page"`
	Limit int                     `json:"limit"`
}
type AdminWorkerResponse struct {
	ID            int64               `json:"id"`
	FullName      string              `json:"full_name"`
	ExpertGroup   ExpertGroupResponse `json:"expert_group"`
	IsResponsible bool                `json:"is_responsible"`
	AssignedAt    time.Time           `json:"assigned_at"`
}
type AdminTicketResponse struct {
	TrackID       string                     `json:"track_id"`
	Category      CategoryResponse           `json:"category"`
	Status        string                     `json:"status"`
	ApplicantType string                     `json:"applicant_type"`
	Priority      string                     `json:"priority"`
	CreatedAt     time.Time                  `json:"created_at"`
	ClosedAt      *time.Time                 `json:"closed_at"`
	ReturnCount   int                        `json:"return_count"`
	Responsible   *ResponsibleWorkerResponse `json:"responsible"`
	RoutingIssue  string                     `json:"routing_issue,omitempty"`
	Workers       []AdminWorkerResponse      `json:"workers,omitempty"`
}
type AdminTicketPageResponse struct {
	Items []AdminTicketResponse `json:"items"`
	Total int                   `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}
type AdminTicketListRequest struct{ Search, Status, Priority, CategoryID, ApplicantType, ResponsibleWorkerID, RoutingIssue, Page, Limit string }
type AdminEventUserResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Role     string `json:"role,omitempty"`
}
type AdminEventResponse struct {
	ID         int64                   `json:"id"`
	EventType  string                  `json:"event_type"`
	Actor      *AdminEventUserResponse `json:"actor"`
	FromStatus string                  `json:"from_status,omitempty"`
	ToStatus   string                  `json:"to_status,omitempty"`
	FromWorker *AdminEventUserResponse `json:"from_worker,omitempty"`
	ToWorker   *AdminEventUserResponse `json:"to_worker,omitempty"`
	ReasonCode string                  `json:"reason_code,omitempty"`
	CreatedAt  time.Time               `json:"created_at"`
}
type AdminEventPageResponse struct {
	Items []AdminEventResponse `json:"items"`
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}
type AdminEmployeeLoadResponse struct {
	ID            int64   `json:"id"`
	FullName      string  `json:"full_name"`
	Role          string  `json:"role"`
	ActiveTickets int     `json:"active_tickets"`
	MaxTickets    int     `json:"max_tickets"`
	LoadPercent   float64 `json:"load_percent"`
	ActionCount   int     `json:"action_count"`
}
type AdminAnalyticsResponse struct {
	AnalyticsResponse
	Employees []AdminEmployeeLoadResponse `json:"employees"`
}

func (service *AdminService) Dashboard(ctx context.Context) (AdminDashboardResponse, error) {
	record, err := service.repository.Dashboard(ctx)
	if err != nil {
		return AdminDashboardResponse{}, err
	}
	return AdminDashboardResponse{
		NewCount: record.NewCount, ActiveCount: record.ActiveCount,
		UrgentCount: record.UrgentCount, ReturnedCount: record.ReturnedCount,
		RoutingIssueCount:      record.RoutingIssueCount,
		ExpertsAtCapacityCount: record.ExpertsAtCapacityCount,
	}, nil
}

func parseAdminPage(pageValue, limitValue string) (int, int, error) {
	page, err := positiveIntOrDefault(pageValue, 1)
	if err != nil {
		return 0, 0, ErrInvalidAdminInput
	}
	limit, err := positiveIntOrDefault(limitValue, 20)
	if err != nil || limit > 100 {
		return 0, 0, ErrInvalidAdminInput
	}
	return page, limit, nil
}

func validAdminRole(role string) bool             { return role == "operator" || role == "expert" }
func validAdminString(value string, max int) bool { return value != "" && len([]rune(value)) <= max }

func employeeResponse(record repos.AdminEmployeeRecord) AdminEmployeeResponse {
	result := AdminEmployeeResponse{
		ID: record.ID, Username: record.Username, FullName: record.FullName,
		Role: record.Role, Active: record.Role != "200", MaxTickets: record.MaxTickets,
		ActiveTickets: record.ActiveTickets,
	}
	if record.Role == "expert" {
		result.ExpertGroup = &ExpertGroupResponse{ID: record.GroupID, Title: record.GroupTitle}
		rating := record.AverageRating
		result.AverageRating = &rating
		if record.MaxTickets > 0 {
			result.LoadPercent = 100 * float64(record.ActiveTickets) / float64(record.MaxTickets)
		}
	}
	if record.Role == "200" {
		result.ExpertGroup = &ExpertGroupResponse{ID: record.GroupID, Title: record.GroupTitle}
	}
	return result
}

func validateAdminEmployee(username, password, fullName, role string, groupID int64, maxTickets int, requirePassword bool) error {
	if !validAdminString(username, 50) || !validAdminString(fullName, 100) || !validAdminRole(role) || groupID <= 0 || maxTickets < 0 {
		return ErrInvalidAdminInput
	}
	if requirePassword && password == "" {
		return ErrInvalidAdminInput
	}
	if len([]byte(password)) > 72 {
		return ErrInvalidAdminInput
	}
	return nil
}

func (service *AdminService) ListEmployees(ctx context.Context, role, search, pageValue, limitValue string) (AdminEmployeePageResponse, error) {
	role = strings.TrimSpace(role)
	if role != "" && !validAdminRole(role) && role != "200" {
		return AdminEmployeePageResponse{}, ErrInvalidAdminInput
	}
	page, limit, err := parseAdminPage(pageValue, limitValue)
	if err != nil {
		return AdminEmployeePageResponse{}, err
	}
	record, err := service.repository.ListEmployees(ctx, role, strings.TrimSpace(search), limit, (page-1)*limit)
	if err != nil {
		return AdminEmployeePageResponse{}, err
	}
	result := AdminEmployeePageResponse{Items: []AdminEmployeeResponse{}, Total: record.Total, Page: page, Limit: limit}
	for _, item := range record.Items {
		result.Items = append(result.Items, employeeResponse(item))
	}
	return result, nil
}

func (service *AdminService) CreateEmployee(ctx context.Context, request CreateAdminEmployeeRequest) (AdminEmployeeResponse, error) {
	request.Username = strings.TrimSpace(request.Username)
	request.FullName = strings.TrimSpace(request.FullName)
	request.Role = strings.TrimSpace(request.Role)
	if err := validateAdminEmployee(request.Username, request.Password, request.FullName, request.Role, request.ExpertGroupID, request.MaxTickets, true); err != nil {
		return AdminEmployeeResponse{}, err
	}
	if request.Role == "operator" {
		request.MaxTickets = 0
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	record, err := service.repository.CreateEmployee(ctx, request.Username, string(hash), request.FullName, request.Role, request.ExpertGroupID, request.MaxTickets)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	return employeeResponse(record), nil
}

func (service *AdminService) UpdateEmployee(ctx context.Context, id int64, request UpdateAdminEmployeeRequest) (AdminEmployeeResponse, error) {
	if id <= 0 {
		return AdminEmployeeResponse{}, ErrAdminResourceNotFound
	}
	if request.Username == nil && request.Password == nil && request.FullName == nil && request.Role == nil && request.ExpertGroupID == nil && request.MaxTickets == nil {
		return AdminEmployeeResponse{}, ErrInvalidAdminInput
	}
	current, err := service.repository.GetEmployee(ctx, id)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	if current.Role == "admin" || current.Role == "200" {
		return AdminEmployeeResponse{}, ErrCannotManageAdmin
	}
	username, fullName, role := current.Username, current.FullName, current.Role
	groupID, maxTickets := current.GroupID, current.MaxTickets
	if request.Username != nil {
		username = strings.TrimSpace(*request.Username)
	}
	if request.FullName != nil {
		fullName = strings.TrimSpace(*request.FullName)
	}
	if request.Role != nil {
		role = strings.TrimSpace(*request.Role)
	}
	if request.ExpertGroupID != nil {
		groupID = *request.ExpertGroupID
	}
	if request.MaxTickets != nil {
		maxTickets = *request.MaxTickets
	}
	password := ""
	if request.Password != nil {
		password = *request.Password
	}
	if role == "operator" {
		maxTickets = 0
	}
	if err = validateAdminEmployee(username, password, fullName, role, groupID, maxTickets, false); err != nil {
		return AdminEmployeeResponse{}, err
	}
	if current.Role == "expert" && role == "operator" && current.ActiveTickets > 0 {
		return AdminEmployeeResponse{}, ErrEmployeeHasTickets
	}
	hash := ""
	if password != "" {
		value, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return AdminEmployeeResponse{}, hashErr
		}
		hash = string(value)
	}
	record, err := service.repository.UpdateEmployee(ctx, id, username, hash, fullName, role, groupID, maxTickets)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	return employeeResponse(record), nil
}

func (service *AdminService) DeactivateEmployee(ctx context.Context, id, actorID int64) (AdminEmployeeResponse, error) {
	if id <= 0 {
		return AdminEmployeeResponse{}, ErrAdminResourceNotFound
	}
	if id == actorID {
		return AdminEmployeeResponse{}, ErrCannotDeactivateSelf
	}
	current, err := service.repository.GetEmployee(ctx, id)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	if current.Role == "admin" {
		return AdminEmployeeResponse{}, ErrCannotManageAdmin
	}
	if current.Role == "200" {
		return employeeResponse(current), nil
	}
	record, err := service.repository.DeactivateEmployee(ctx, id)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	return employeeResponse(record), nil
}

func (service *AdminService) RestoreEmployee(ctx context.Context, id int64, request RestoreAdminEmployeeRequest) (AdminEmployeeResponse, error) {
	current, err := service.repository.GetEmployee(ctx, id)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	if current.Role != "200" {
		return AdminEmployeeResponse{}, ErrEmployeeNotInactive
	}
	request.Role = strings.TrimSpace(request.Role)
	if request.Role == "operator" {
		request.MaxTickets = 0
	}
	if err = validateAdminEmployee(current.Username, request.Password, current.FullName, request.Role, request.ExpertGroupID, request.MaxTickets, false); err != nil {
		return AdminEmployeeResponse{}, err
	}
	hash := ""
	if request.Password != "" {
		value, hashErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return AdminEmployeeResponse{}, hashErr
		}
		hash = string(value)
	}
	record, err := service.repository.UpdateEmployee(ctx, id, current.Username, hash, current.FullName, request.Role, request.ExpertGroupID, request.MaxTickets)
	if err != nil {
		return AdminEmployeeResponse{}, err
	}
	return employeeResponse(record), nil
}

func groupResponse(record repos.AdminGroupRecord) AdminGroupResponse {
	return AdminGroupResponse{ID: record.ID, Title: record.Title, EmployeeCount: record.EmployeeCount, CategoryCount: record.CategoryCount}
}

func validateAdminName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if !validAdminString(name, 50) {
		return "", ErrInvalidAdminInput
	}
	return name, nil
}

func (service *AdminService) ListGroups(ctx context.Context, search, pageValue, limitValue string) (AdminGroupPageResponse, error) {
	page, limit, err := parseAdminPage(pageValue, limitValue)
	if err != nil {
		return AdminGroupPageResponse{}, err
	}
	record, err := service.repository.ListGroups(ctx, strings.TrimSpace(search), limit, (page-1)*limit)
	if err != nil {
		return AdminGroupPageResponse{}, err
	}
	result := AdminGroupPageResponse{Items: []AdminGroupResponse{}, Total: record.Total, Page: page, Limit: limit}
	for _, item := range record.Items {
		result.Items = append(result.Items, groupResponse(item))
	}
	return result, nil
}

func (service *AdminService) CreateGroup(ctx context.Context, title string) (AdminGroupResponse, error) {
	title, err := validateAdminName(title)
	if err != nil {
		return AdminGroupResponse{}, err
	}
	record, err := service.repository.CreateGroup(ctx, title)
	if err != nil {
		return AdminGroupResponse{}, err
	}
	return groupResponse(record), nil
}

func (service *AdminService) UpdateGroup(ctx context.Context, id int64, title string) (AdminGroupResponse, error) {
	title, err := validateAdminName(title)
	if err != nil {
		return AdminGroupResponse{}, err
	}
	record, err := service.repository.UpdateGroup(ctx, id, title)
	if err != nil {
		return AdminGroupResponse{}, err
	}
	return groupResponse(record), nil
}

func (service *AdminService) DeleteGroup(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrAdminResourceNotFound
	}
	return service.repository.DeleteGroup(ctx, id)
}

func questionResponse(record repos.AdminQuestionRecord) AdminQuestionResponse {
	result := AdminQuestionResponse{ID: record.ID, Text: record.Text, Answers: []AdminAnswerResponse{}}
	for _, answer := range record.Answers {
		result.Answers = append(result.Answers, AdminAnswerResponse{ID: answer.ID, Text: answer.Text})
	}
	return result
}

func categoryResponse(record repos.AdminCategoryRecord) AdminCategoryResponse {
	result := AdminCategoryResponse{ID: record.ID, Name: record.Name, ExpertGroups: []AdminGroupResponse{}, Questions: []AdminQuestionResponse{}}
	for _, group := range record.Groups {
		result.ExpertGroups = append(result.ExpertGroups, groupResponse(group))
	}
	for _, question := range record.Questions {
		result.Questions = append(result.Questions, questionResponse(question))
	}
	return result
}

func (service *AdminService) ListCategories(ctx context.Context, search, pageValue, limitValue string) (AdminCategoryPageResponse, error) {
	page, limit, err := parseAdminPage(pageValue, limitValue)
	if err != nil {
		return AdminCategoryPageResponse{}, err
	}
	record, err := service.repository.ListCategories(ctx, strings.TrimSpace(search), limit, (page-1)*limit)
	if err != nil {
		return AdminCategoryPageResponse{}, err
	}
	result := AdminCategoryPageResponse{Items: []AdminCategoryResponse{}, Total: record.Total, Page: page, Limit: limit}
	for _, item := range record.Items {
		detail, detailErr := service.repository.GetCategory(ctx, item.ID)
		if detailErr != nil {
			return AdminCategoryPageResponse{}, detailErr
		}
		result.Items = append(result.Items, categoryResponse(detail))
	}
	return result, nil
}

func (service *AdminService) Category(ctx context.Context, id int64) (AdminCategoryResponse, error) {
	record, err := service.repository.GetCategory(ctx, id)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	return categoryResponse(record), nil
}

func (service *AdminService) CreateCategory(ctx context.Context, name string) (AdminCategoryResponse, error) {
	name, err := validateAdminName(name)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	record, err := service.repository.CreateCategory(ctx, name)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	return categoryResponse(record), nil
}

func (service *AdminService) UpdateCategory(ctx context.Context, id int64, name string) (AdminCategoryResponse, error) {
	name, err := validateAdminName(name)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	current, err := service.repository.GetCategory(ctx, id)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	if strings.EqualFold(strings.TrimSpace(current.Name), "Не знаю, как это назвать") && !strings.EqualFold(name, current.Name) {
		return AdminCategoryResponse{}, ErrAdminResourceInUse
	}
	record, err := service.repository.UpdateCategory(ctx, id, name)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	return categoryResponse(record), nil
}

func (service *AdminService) DeleteCategory(ctx context.Context, id int64) error {
	return service.repository.DeleteCategory(ctx, id)
}

func uniqueAdminIDs(ids []int64) ([]int64, error) {
	seen := map[int64]bool{}
	result := []int64{}
	for _, id := range ids {
		if id <= 0 {
			return nil, ErrInvalidAdminInput
		}
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result, nil
}

func (service *AdminService) ReplaceCategoryGroups(ctx context.Context, id int64, ids []int64) (AdminCategoryResponse, error) {
	ids, err := uniqueAdminIDs(ids)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	record, err := service.repository.ReplaceCategoryGroups(ctx, id, ids)
	if err != nil {
		return AdminCategoryResponse{}, err
	}
	return categoryResponse(record), nil
}

func normalizeAdminQuestion(text string, answers []string) (string, []string, error) {
	text = strings.TrimSpace(text)
	if text == "" || len(answers) == 0 {
		return "", nil, ErrInvalidAdminInput
	}
	seen := map[string]bool{}
	result := []string{}
	for _, answer := range answers {
		answer = strings.TrimSpace(answer)
		key := strings.ToLower(answer)
		if answer == "" || seen[key] {
			return "", nil, ErrInvalidAdminInput
		}
		seen[key] = true
		result = append(result, answer)
	}
	return text, result, nil
}

func (service *AdminService) CreateQuestion(ctx context.Context, categoryID int64, text string, answers []string) (AdminQuestionResponse, error) {
	text, answers, err := normalizeAdminQuestion(text, answers)
	if err != nil {
		return AdminQuestionResponse{}, err
	}
	record, err := service.repository.CreateQuestion(ctx, categoryID, text, answers)
	if err != nil {
		return AdminQuestionResponse{}, err
	}
	return questionResponse(record), nil
}

func (service *AdminService) UpdateQuestion(ctx context.Context, id int64, text string, answers []string) (AdminQuestionResponse, error) {
	text, answers, err := normalizeAdminQuestion(text, answers)
	if err != nil {
		return AdminQuestionResponse{}, err
	}
	record, err := service.repository.UpdateQuestion(ctx, id, text, answers)
	if err != nil {
		return AdminQuestionResponse{}, err
	}
	return questionResponse(record), nil
}

func (service *AdminService) DeleteQuestion(ctx context.Context, id int64) error {
	return service.repository.DeleteQuestion(ctx, id)
}

func ticketResponse(record repos.AdminTicketRecord) AdminTicketResponse {
	result := AdminTicketResponse{
		TrackID: record.TrackID, Category: CategoryResponse{ID: record.CategoryID, Name: record.Category},
		Status: record.Status.String(), ApplicantType: record.ApplicantType.String(), Priority: record.Priority.String(),
		CreatedAt: record.CreatedAt.UTC(), ReturnCount: record.ReturnCount, RoutingIssue: record.RoutingIssue.String,
	}
	if record.ClosedAt.Valid {
		closedAt := record.ClosedAt.Time.UTC()
		result.ClosedAt = &closedAt
	}
	if record.ResponsibleID.Valid {
		result.Responsible = &ResponsibleWorkerResponse{ID: record.ResponsibleID.Int64, FullName: record.ResponsibleName.String}
	}
	if record.Workers != nil {
		result.Workers = []AdminWorkerResponse{}
		for _, worker := range record.Workers {
			result.Workers = append(result.Workers, AdminWorkerResponse{
				ID: worker.ID, FullName: worker.FullName,
				ExpertGroup:   ExpertGroupResponse{ID: worker.GroupID, Title: worker.GroupTitle},
				IsResponsible: worker.IsResponsible, AssignedAt: worker.CreatedAt.UTC(),
			})
		}
	}
	return result
}

func (service *AdminService) ListTickets(ctx context.Context, request AdminTicketListRequest) (AdminTicketPageResponse, error) {
	page, limit, err := parseAdminPage(request.Page, request.Limit)
	if err != nil {
		return AdminTicketPageResponse{}, err
	}
	filter := repos.AdminTicketFilter{Search: strings.TrimSpace(request.Search), RoutingIssue: strings.TrimSpace(request.RoutingIssue), Limit: limit, Offset: (page - 1) * limit}
	if filter.RoutingIssue != "" && filter.RoutingIssue != "no_group" && filter.RoutingIssue != "no_available_expert" {
		return AdminTicketPageResponse{}, ErrInvalidAdminInput
	}
	if request.Status != "" {
		filter.Status = statusFromString(request.Status)
		if filter.Status == 0 {
			return AdminTicketPageResponse{}, ErrInvalidAdminInput
		}
	}
	if request.Priority != "" {
		switch request.Priority {
		case "standard":
			filter.Priority = models.TicketPriorityStandard
		case "urgent":
			filter.Priority = models.TicketPriorityUrgent
		default:
			return AdminTicketPageResponse{}, ErrInvalidAdminInput
		}
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
			return AdminTicketPageResponse{}, ErrInvalidAdminInput
		}
	}
	filter.CategoryID, err = optionalPositiveInt64(request.CategoryID)
	if err != nil {
		return AdminTicketPageResponse{}, ErrInvalidAdminInput
	}
	filter.ResponsibleID, err = optionalPositiveInt64(request.ResponsibleWorkerID)
	if err != nil {
		return AdminTicketPageResponse{}, ErrInvalidAdminInput
	}
	record, err := service.repository.ListTickets(ctx, filter)
	if err != nil {
		return AdminTicketPageResponse{}, err
	}
	result := AdminTicketPageResponse{Items: []AdminTicketResponse{}, Total: record.Total, Page: page, Limit: limit}
	for _, item := range record.Items {
		result.Items = append(result.Items, ticketResponse(item))
	}
	return result, nil
}

func (service *AdminService) Ticket(ctx context.Context, trackID string) (AdminTicketResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return AdminTicketResponse{}, err
	}
	record, err := service.repository.GetTicket(ctx, trackID)
	if err != nil {
		return AdminTicketResponse{}, err
	}
	return ticketResponse(record), nil
}

func (service *AdminService) ChangePriority(ctx context.Context, trackID, priority string, actorID int64) (AdminTicketResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return AdminTicketResponse{}, err
	}
	var value models.TicketPriority
	switch priority {
	case "standard":
		value = models.TicketPriorityStandard
	case "urgent":
		value = models.TicketPriorityUrgent
	default:
		return AdminTicketResponse{}, ErrInvalidAdminInput
	}
	record, err := service.repository.UpdateTicketPriority(ctx, trackID, value, actorID, service.now().UTC())
	if err != nil {
		return AdminTicketResponse{}, err
	}
	return ticketResponse(record), nil
}

func (service *AdminService) ChangeStatus(ctx context.Context, trackID, status string, actorID int64) (AdminTicketResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return AdminTicketResponse{}, err
	}
	value := statusFromString(status)
	if value == 0 {
		return AdminTicketResponse{}, ErrInvalidAdminInput
	}
	record, err := service.repository.UpdateTicketStatus(ctx, trackID, value, actorID, service.now().UTC())
	if err != nil {
		return AdminTicketResponse{}, err
	}
	return ticketResponse(record), nil
}

func (service *AdminService) SetResponsible(ctx context.Context, trackID string, workerID, actorID int64) (AdminTicketResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return AdminTicketResponse{}, err
	}
	if workerID <= 0 {
		return AdminTicketResponse{}, ErrInvalidAdminInput
	}
	record, err := service.repository.SetResponsibleWorker(ctx, trackID, workerID, actorID, service.now().UTC())
	if err != nil {
		return AdminTicketResponse{}, err
	}
	return ticketResponse(record), nil
}

func (service *AdminService) Events(ctx context.Context, trackID, pageValue, limitValue string) (AdminEventPageResponse, error) {
	trackID, err := normalizeTrackID(trackID)
	if err != nil {
		return AdminEventPageResponse{}, err
	}
	page, limit, err := parseAdminPage(pageValue, limitValue)
	if err != nil {
		return AdminEventPageResponse{}, err
	}
	record, err := service.repository.ListEvents(ctx, trackID, limit, (page-1)*limit)
	if err != nil {
		return AdminEventPageResponse{}, err
	}
	result := AdminEventPageResponse{Items: []AdminEventResponse{}, Total: record.Total, Page: page, Limit: limit}
	for _, item := range record.Items {
		event := AdminEventResponse{ID: item.ID, EventType: item.EventType, ReasonCode: item.ReasonCode.String, CreatedAt: item.CreatedAt.UTC()}
		if item.ActorID.Valid {
			event.Actor = &AdminEventUserResponse{ID: item.ActorID.Int64, FullName: item.ActorName.String, Role: item.ActorRole.String}
		}
		if item.FromStatus.Valid {
			event.FromStatus = models.TicketStatus(item.FromStatus.Int64).String()
		}
		if item.ToStatus.Valid {
			event.ToStatus = models.TicketStatus(item.ToStatus.Int64).String()
		}
		if item.FromWorkerID.Valid {
			event.FromWorker = &AdminEventUserResponse{ID: item.FromWorkerID.Int64, FullName: item.FromWorkerName.String}
		}
		if item.ToWorkerID.Valid {
			event.ToWorker = &AdminEventUserResponse{ID: item.ToWorkerID.Int64, FullName: item.ToWorkerName.String}
		}
		result.Items = append(result.Items, event)
	}
	return result, nil
}

func (service *AdminService) dateRange(from, to string) (time.Time, time.Time, error) {
	start, err := time.ParseInLocation(time.DateOnly, from, service.location)
	if err != nil {
		return time.Time{}, time.Time{}, ErrInvalidDateRange
	}
	last, err := time.ParseInLocation(time.DateOnly, to, service.location)
	if err != nil || last.Before(start) {
		return time.Time{}, time.Time{}, ErrInvalidDateRange
	}
	return start.UTC(), last.AddDate(0, 0, 1).UTC(), nil
}

func (service *AdminService) Analytics(ctx context.Context, from, to string) (AdminAnalyticsResponse, error) {
	start, end, err := service.dateRange(from, to)
	if err != nil {
		return AdminAnalyticsResponse{}, err
	}
	record, err := service.repository.Analytics(ctx, start, end)
	if err != nil {
		return AdminAnalyticsResponse{}, err
	}
	base := AnalyticsResponse{
		DateFrom: from, DateTo: to, TotalTickets: record.Total,
		Categories: []DistributionResponse{}, ApplicantTypes: []DistributionResponse{}, Statuses: []DistributionResponse{},
		AverageAcceptanceSeconds: record.DurationSecondsToAccept, AverageFirstResponseSeconds: record.DurationSecondsToResponse,
		AverageCloseSeconds: record.DurationSecondsToClose, ExpertLoadPercent: record.ExpertLoadPercent,
		UrgentSharePercent: percent(record.UrgentCount, record.Total), ReturnSharePercent: percent(record.ReturnedCount, record.Total),
	}
	for _, item := range record.Categories {
		base.Categories = append(base.Categories, DistributionResponse{ID: item.ID, Value: item.Name, Count: item.Count, Percent: percent(item.Count, record.Total)})
	}
	for _, item := range record.ApplicantTypes {
		base.ApplicantTypes = append(base.ApplicantTypes, DistributionResponse{Value: models.ApplicantType(item.Value).String(), Count: item.Count, Percent: percent(item.Count, record.Total)})
	}
	for _, item := range record.Statuses {
		base.Statuses = append(base.Statuses, DistributionResponse{Value: models.TicketStatus(item.Value).String(), Count: item.Count, Percent: percent(item.Count, record.Total)})
	}
	loads, err := service.repository.EmployeeLoads(ctx, start, end)
	if err != nil {
		return AdminAnalyticsResponse{}, err
	}
	result := AdminAnalyticsResponse{AnalyticsResponse: base, Employees: []AdminEmployeeLoadResponse{}}
	for _, load := range loads {
		item := AdminEmployeeLoadResponse{ID: load.ID, FullName: load.FullName, Role: load.Role, ActiveTickets: load.ActiveTickets, MaxTickets: load.MaxTickets, ActionCount: load.ActionCount}
		if load.Role == "expert" && load.MaxTickets > 0 {
			item.LoadPercent = 100 * float64(load.ActiveTickets) / float64(load.MaxTickets)
		}
		result.Employees = append(result.Employees, item)
	}
	return result, nil
}

func (service *AdminService) Report(ctx context.Context, from, to, format string) (GeneratedReport, error) {
	start, end, err := service.dateRange(from, to)
	if err != nil {
		return GeneratedReport{}, err
	}
	records, err := service.repository.ReportTickets(ctx, start, end)
	if err != nil {
		return GeneratedReport{}, err
	}
	format = strings.ToLower(strings.TrimSpace(format))
	name := "admin-tickets-" + from + "-" + to
	switch format {
	case "csv":
		data, buildErr := buildCSVReport(records)
		return GeneratedReport{Name: name + ".csv", ContentType: "text/csv; charset=utf-8", Data: data}, buildErr
	case "xlsx":
		data, buildErr := buildXLSXReport(records)
		return GeneratedReport{Name: name + ".xlsx", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Data: data}, buildErr
	default:
		return GeneratedReport{}, ErrInvalidReportFormat
	}
}
