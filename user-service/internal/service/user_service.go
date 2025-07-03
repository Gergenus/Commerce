package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/Gergenus/commerce/user-service/internal/repository"
	"github.com/Gergenus/commerce/user-service/pkg/hash"
	"github.com/Gergenus/commerce/user-service/pkg/jwtpkg"
	"github.com/Gergenus/commerce/user-service/pkg/utils"
	"github.com/Gergenus/commerce/user-service/pkg/validation"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrIncorrectEmail    = errors.New("incorrect email")
	ErrPasswordMismatch  = errors.New("passwords mismatch")
)

type UserService struct {
	log      *slog.Logger
	repo     repository.RepositoryInterface
	jwtToken jwtpkg.UserJWTInterface
}

func NewUserService(log *slog.Logger, repo repository.RepositoryInterface, jwtToken jwtpkg.UserJWTInterface) UserService {
	return UserService{log: log, repo: repo, jwtToken: jwtToken}
}

func (u *UserService) AddUser(ctx context.Context, user models.User) (*uuid.UUID, error) {
	const op = "service.AddUser"
	u.log.With(slog.String("op", op))
	u.log.Info("Creating user", slog.String("email", user.Email))

	if !validation.IsEmailValid(user.Email) {
		u.log.Error("failed to validate email")
		return nil, fmt.Errorf("%s: %w", op, ErrIncorrectEmail)
	}

	hashPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		slog.Error("failed to generate hashpassword", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Password = hashPassword
	uid, err := u.repo.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			u.log.Error("user already exists error", slog.String("email", user.Email))
			return nil, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		u.log.Error("creating user error", slog.String("email", user.Email), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	u.log.Info("created user", slog.String("email", user.Email))
	return uid, nil
}

func (u *UserService) Login(ctx context.Context, email, password, userAgent, ip string) (string, string, error) {
	const op = "service.Login"
	u.log.With(slog.String("op", op))
	u.log.Info("logigng the user", slog.String("email", email))

	user, err := u.repo.GetUser(ctx, email)
	if err != nil {
		u.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if !hash.CheckPassword(user.Password, password) {
		u.log.Info("passwords mismatch")
		return "", "", fmt.Errorf("%s: %w", op, ErrPasswordMismatch)
	}
	AccessToken, err := u.jwtToken.GenerateAccessToken(*user)
	if err != nil {
		u.log.Error("failed to create access token", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	RefreshToken := uuid.New().String()

	fingerprint := utils.CreateFingerprint(ip, userAgent)

	err = u.repo.CreateJWTSession(ctx, *user, RefreshToken, fingerprint, ip, time.Now().Add(7*24*time.Hour).Unix())
	if err != nil {
		u.log.Error("failed to create jwt session", slog.String("error", err.Error()))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	return AccessToken, RefreshToken, nil
}
