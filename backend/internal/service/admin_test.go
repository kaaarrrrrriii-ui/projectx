package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
	"golang.org/x/crypto/bcrypt"
)

type stubAdminRepository struct {
	adminRepository
	employee            repos.AdminEmployeeRecord
	createdPasswordHash string
	createdGroupID      int64
	createdMaxTickets   int
	deactivateCalled    bool
	category            repos.AdminCategoryRecord
	status              models.TicketStatus
}

func (stub *stubAdminRepository) GetEmployee(context.Context, int64) (repos.AdminEmployeeRecord, error) {
	return stub.employee, nil
}

func (stub *stubAdminRepository) CreateEmployee(_ context.Context, username, passwordHash, fullName, role string, groupID int64, maxTickets int) (repos.AdminEmployeeRecord, error) {
	stub.createdPasswordHash = passwordHash
	stub.createdGroupID = groupID
	stub.createdMaxTickets = maxTickets
	return repos.AdminEmployeeRecord{ID: 9, Username: username, FullName: fullName, Role: role, GroupID: groupID, MaxTickets: maxTickets}, nil
}

func (stub *stubAdminRepository) DeactivateEmployee(context.Context, int64) (repos.AdminEmployeeRecord, error) {
	stub.deactivateCalled = true
	stub.employee.Role = "200"
	return stub.employee, nil
}

func (stub *stubAdminRepository) GetCategory(context.Context, int64) (repos.AdminCategoryRecord, error) {
	return stub.category, nil
}

func (stub *stubAdminRepository) UpdateTicketStatus(_ context.Context, trackID string, status models.TicketStatus, actorID int64, changedAt time.Time) (repos.AdminTicketRecord, error) {
	stub.status = status
	return repos.AdminTicketRecord{TrackID: trackID, Status: status, Category: "Категория", CategoryID: 1, Priority: models.TicketPriorityStandard, ApplicantType: models.ApplicantTypeSchoolchild, CreatedAt: changedAt}, nil
}

func newStubAdminService(t *testing.T, repository adminRepository) *AdminService {
	t.Helper()
	result, err := NewAdminService(repository)
	if err != nil {
		t.Fatalf("NewAdminService() error = %v", err)
	}
	return result
}

func TestAdminCreateOperatorHashesPasswordAndUsesSelectedGroup(t *testing.T) {
	t.Parallel()
	repository := &stubAdminRepository{}
	admin := newStubAdminService(t, repository)
	response, err := admin.CreateEmployee(context.Background(), CreateAdminEmployeeRequest{
		Username: " operator ", Password: "secret", FullName: " Оператор ", Role: "operator", ExpertGroupID: 7, MaxTickets: 99,
	})
	if err != nil {
		t.Fatalf("CreateEmployee() error = %v", err)
	}
	if repository.createdGroupID != 7 || repository.createdMaxTickets != 0 || response.ExpertGroup != nil {
		t.Fatalf("created employee = %+v, group = %d, max = %d", response, repository.createdGroupID, repository.createdMaxTickets)
	}
	if bcrypt.CompareHashAndPassword([]byte(repository.createdPasswordHash), []byte("secret")) != nil {
		t.Fatal("password was not bcrypt hashed")
	}
}

func TestAdminCannotDeactivateSelfOrAdmin(t *testing.T) {
	t.Parallel()
	repository := &stubAdminRepository{employee: repos.AdminEmployeeRecord{ID: 2, Role: "operator"}}
	admin := newStubAdminService(t, repository)
	if _, err := admin.DeactivateEmployee(context.Background(), 2, 2); !errors.Is(err, ErrCannotDeactivateSelf) {
		t.Fatalf("self deactivation error = %v", err)
	}
	repository.employee.Role = "admin"
	if _, err := admin.DeactivateEmployee(context.Background(), 2, 1); !errors.Is(err, ErrCannotManageAdmin) {
		t.Fatalf("admin deactivation error = %v", err)
	}
	if repository.deactivateCalled {
		t.Fatal("repository deactivation must not be called")
	}
}

func TestAdminProtectsRequiredCategory(t *testing.T) {
	t.Parallel()
	repository := &stubAdminRepository{category: repos.AdminCategoryRecord{ID: 1, Name: "Не знаю, как это назвать"}}
	admin := newStubAdminService(t, repository)
	if _, err := admin.UpdateCategory(context.Background(), 1, "Другое"); !errors.Is(err, ErrAdminResourceInUse) {
		t.Fatalf("UpdateCategory() error = %v", err)
	}
}

func TestAdminRejectsDuplicateQuestionAnswers(t *testing.T) {
	t.Parallel()
	admin := newStubAdminService(t, &stubAdminRepository{})
	if _, err := admin.CreateQuestion(context.Background(), 1, "Вопрос", []string{"Да", " да "}); !errors.Is(err, ErrInvalidAdminInput) {
		t.Fatalf("CreateQuestion() error = %v", err)
	}
}

func TestAdminMaySetTerminalStatus(t *testing.T) {
	t.Parallel()
	repository := &stubAdminRepository{}
	admin := newStubAdminService(t, repository)
	admin.now = func() time.Time { return time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC) }
	if _, err := admin.ChangeStatus(context.Background(), "ОТК-ABCD-2345", "closed_without_answer", 3); err != nil {
		t.Fatalf("ChangeStatus() error = %v", err)
	}
	if repository.status != models.TicketStatusClosedWithoutAnswer {
		t.Fatalf("status = %v", repository.status)
	}
}
