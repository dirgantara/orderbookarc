package service

import (
	"context"
	"errors"
	"time"
	"trader/internal/model"
	"trader/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

type RegisterInput struct {
	FullName string `json:"full_name" binding:"required,min=3"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:           uuid.New().String(),
		FullName:     in.FullName,
		Email:        in.Email,
		PasswordHash: string(hash),
		IsActive:     true,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (string, *model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, in.Email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return "", nil, errors.New("account inactive")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return signed, user, nil
}
