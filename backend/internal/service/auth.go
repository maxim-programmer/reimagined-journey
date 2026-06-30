package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/maxim-programmer/reimagined-journey/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type userRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}

type sessionStore interface {
	SetSession(ctx context.Context, token, userID string, ttl time.Duration) error
	GetSession(ctx context.Context, token string) (string, error)
	DeleteSession(ctx context.Context, token string) error
}

const sessionTTL = 24 * time.Hour

type AuthService struct {
	users    userRepo
	sessions sessionStore
}

func NewAuthService(users userRepo, sessions sessionStore) *AuthService {
	return &AuthService{users: users, sessions: sessions}
}

func (s *AuthService) Register(ctx context.Context, login, password string) (*model.User, error) {
	existing, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("check login: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("login already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Login:        login,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, *model.User, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		return "", nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := generateToken()
	if err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
	}

	if err := s.sessions.SetSession(ctx, token, user.ID, sessionTTL); err != nil {
		return "", nil, fmt.Errorf("save session: %w", err)
	}

	return token, user, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessions.DeleteSession(ctx, token)
}

func (s *AuthService) GetUserByToken(ctx context.Context, token string) (*model.User, error) {
	userID, err := s.sessions.GetSession(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	if userID == "" {
		return nil, nil
	}
	return s.users.GetByID(ctx, userID)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
