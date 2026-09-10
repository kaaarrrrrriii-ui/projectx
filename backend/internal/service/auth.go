package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"example.com/german/backend/internal/models"
	"example.com/german/backend/internal/repos"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid authentication token")
	ErrForbidden          = errors.New("operation is forbidden")
)

const DefaultAuthTokenTTL = 8 * time.Hour

type authUserRepository interface {
	GetByUsername(context.Context, string) (models.User, error)
	GetByID(context.Context, int64) (models.User, error)
}

type AuthService struct {
	repository authUserRepository
	secret     []byte
	tokenTTL   time.Duration
	now        func() time.Time
}

type AuthUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      AuthUser  `json:"user"`
}

type tokenClaims struct {
	Subject   string `json:"sub"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func NewAuthService(repository authUserRepository, secret string) (*AuthService, error) {
	if repository == nil {
		return nil, errors.New("user repository is required")
	}
	if len(secret) < 32 {
		return nil, errors.New("JWT secret must contain at least 32 characters")
	}
	return &AuthService{
		repository: repository,
		secret:     []byte(secret),
		tokenTTL:   DefaultAuthTokenTTL,
		now:        time.Now,
	}, nil
}

func (service *AuthService) Login(ctx context.Context, username, password string) (LoginResponse, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return LoginResponse{}, ErrInvalidCredentials
	}
	user, err := service.repository.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repos.ErrUserNotFound) {
			return LoginResponse{}, ErrInvalidCredentials
		}
		return LoginResponse{}, err
	}
	if user.Role == "200" {
		return LoginResponse{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	now := service.now().UTC()
	expiresAt := now.Add(service.tokenTTL)
	token, err := service.sign(tokenClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		Role:      user.Role,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	})
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{Token: token, ExpiresAt: expiresAt, User: authUser(user)}, nil
}

func (service *AuthService) Authenticate(ctx context.Context, authorization string) (AuthUser, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return AuthUser{}, ErrInvalidToken
	}
	claims, err := service.verify(strings.TrimSpace(strings.TrimPrefix(authorization, prefix)))
	if err != nil {
		return AuthUser{}, err
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return AuthUser{}, ErrInvalidToken
	}
	user, err := service.repository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repos.ErrUserNotFound) {
			return AuthUser{}, ErrInvalidToken
		}
		return AuthUser{}, err
	}
	if user.Role != claims.Role {
		return AuthUser{}, ErrInvalidToken
	}
	return authUser(user), nil
}

func RequireOperator(user AuthUser) error {
	if user.Role != "operator" {
		return ErrForbidden
	}
	return nil
}

func RequireAdmin(user AuthUser) error {
	if user.Role != "admin" {
		return ErrForbidden
	}
	return nil
}

func RequireExpert(user AuthUser) error {
	if user.Role != "expert" {
		return ErrForbidden
	}
	return nil
}

func authUser(user models.User) AuthUser {
	return AuthUser{ID: user.ID, Username: user.Username, FullName: user.FullName, Role: user.Role}
}

func (service *AuthService) sign(claims tokenClaims) (string, error) {
	header, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("encode JWT claims: %w", err)
	}
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(service.signature(unsigned)), nil
}

func (service *AuthService) verify(token string) (tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return tokenClaims{}, ErrInvalidToken
	}
	unsigned := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(signature, service.signature(unsigned)) {
		return tokenClaims{}, ErrInvalidToken
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return tokenClaims{}, ErrInvalidToken
	}
	var header map[string]string
	if json.Unmarshal(headerBytes, &header) != nil || header["alg"] != "HS256" {
		return tokenClaims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return tokenClaims{}, ErrInvalidToken
	}
	var claims tokenClaims
	if json.Unmarshal(payload, &claims) != nil || claims.Subject == "" || claims.ExpiresAt <= service.now().UTC().Unix() {
		return tokenClaims{}, ErrInvalidToken
	}
	return claims, nil
}

func (service *AuthService) signature(value string) []byte {
	mac := hmac.New(sha256.New, service.secret)
	_, _ = mac.Write([]byte(value))
	return mac.Sum(nil)
}
