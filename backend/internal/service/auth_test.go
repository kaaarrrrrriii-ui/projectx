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

type stubAuthUserRepository struct {
	user models.User
	err  error
}

func (repository *stubAuthUserRepository) GetByUsername(context.Context, string) (models.User, error) {
	return repository.user, repository.err
}

func (repository *stubAuthUserRepository) GetByID(context.Context, int64) (models.User, error) {
	return repository.user, repository.err
}

func TestAuthServiceLoginAndAuthenticate(t *testing.T) {
	t.Parallel()
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	repository := &stubAuthUserRepository{user: models.User{
		ID: 12, Username: "operator", PasswordHash: string(hash), FullName: "Олег", Role: "operator",
	}}
	auth, err := NewAuthService(repository, "01234567890123456789012345678901")
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	auth.now = func() time.Time { return time.Date(2026, 9, 10, 8, 0, 0, 0, time.UTC) }

	login, err := auth.Login(context.Background(), " operator ", "correct-password")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if login.Token == "" || login.User.ID != 12 || login.ExpiresAt.Sub(auth.now()) != DefaultAuthTokenTTL {
		t.Fatalf("Login() response = %+v", login)
	}
	user, err := auth.Authenticate(context.Background(), "Bearer "+login.Token)
	if err != nil || user.ID != 12 || user.Role != "operator" {
		t.Fatalf("Authenticate() = (%+v, %v)", user, err)
	}
}

func TestAuthServiceRejectsBadPasswordAndToken(t *testing.T) {
	t.Parallel()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repository := &stubAuthUserRepository{user: models.User{ID: 1, PasswordHash: string(hash), Role: "operator"}}
	auth, _ := NewAuthService(repository, "01234567890123456789012345678901")
	if _, err := auth.Login(context.Background(), "operator", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
	if _, err := auth.Authenticate(context.Background(), "Bearer invalid"); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Authenticate() error = %v, want %v", err, ErrInvalidToken)
	}
	repository.err = repos.ErrUserNotFound
	if _, err := auth.Login(context.Background(), "missing", "password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() missing user error = %v", err)
	}
}

func TestAuthServiceRejectsDeactivatedUser(t *testing.T) {
	t.Parallel()
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repository := &stubAuthUserRepository{user: models.User{ID: 1, PasswordHash: string(hash), Role: "200"}}
	auth, _ := NewAuthService(repository, "01234567890123456789012345678901")
	if _, err := auth.Login(context.Background(), "inactive", "correct"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestRequireOperator(t *testing.T) {
	t.Parallel()
	if err := RequireOperator(AuthUser{Role: "operator"}); err != nil {
		t.Fatalf("RequireOperator(operator) error = %v", err)
	}
	if err := RequireOperator(AuthUser{Role: "admin"}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("RequireOperator(admin) error = %v", err)
	}
	if err := RequireOperator(AuthUser{Role: "expert"}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("RequireOperator(expert) error = %v", err)
	}
}

func TestRequireAdmin(t *testing.T) {
	if err := RequireAdmin(AuthUser{Role: "admin"}); err != nil {
		t.Fatalf("RequireAdmin(admin) error = %v", err)
	}
	if err := RequireAdmin(AuthUser{Role: "operator"}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("RequireAdmin(operator) error = %v", err)
	}
}
