package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"otusgruz/internal/apperr"
	authhttp "otusgruz/internal/client/authhttp"
	"otusgruz/internal/models"
	query "otusgruz/internal/repo"
)

type repo interface {
	GetUserByEmail(ctx context.Context, email string) (query.User, error)
	IsAuth(ctx context.Context, guid uuid.UUID) (query.LoggedIn, error)
	InsertSession(ctx context.Context, userGuid uuid.UUID) error
	InsertUser(ctx context.Context, arg query.InsertUserParams) error
}

type authClient interface {
	LoginRequest(ctx context.Context, email string, password string) (authhttp.LoginResponse, error)
	SignupRequest(ctx context.Context, email string, password string) (authhttp.SignupResponse, error)
}

type service struct {
	repo       repo
	authClient authClient
}

type Service interface {
	Auth(ctx context.Context, guid uuid.UUID) (uuid.UUID, error)
	Login(ctx context.Context, email string, pwd string) (uuid.UUID, error)
	Singup(ctx context.Context, params models.UserSignup) (*models.DefaultStatusResponse, error)
}

func NewService(repo repo, authClient authClient) Service {
	return &service{
		repo:       repo,
		authClient: authClient,
	}
}

func (s *service) Auth(ctx context.Context, guid uuid.UUID) (uuid.UUID, error) {
	isAuth, err := s.repo.IsAuth(ctx, guid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil
		}

		return uuid.Nil, fmt.Errorf("checking is user auth: %w", err)
	}

	return isAuth.UserGuid, nil
}

func (s *service) Login(ctx context.Context, email string, pwd string) (uuid.UUID, error) {
	zerolog.Ctx(ctx).Info().Msg("entered into login method")
	resp, err := s.authClient.LoginRequest(ctx, email, pwd)
	if err != nil {
		if errors.Is(err, apperr.ErrNotCorrectData) {
			return uuid.Nil, apperr.ErrNotCorrectData
		}

		return uuid.Nil, fmt.Errorf("login request: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("response", resp).Msg("finished auth client request")

	return resp.UserGUID, nil
}

func (s *service) Singup(ctx context.Context, params models.UserSignup) (*models.DefaultStatusResponse, error) {
	zerolog.Ctx(ctx).Info().Msg("entered into Singup method")

	email := string(params.Email)
	occupation := params.Occupation
	name := params.Name
	pwd := params.Password

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("getting user by email: %w", err)
	}

	if !lo.IsEmpty(user) {
		return nil, apperr.ErrEmailAlreadyUsed
	}

	resp, err := s.authClient.SignupRequest(ctx, email, pwd)
	if err != nil {
		if errors.Is(err, apperr.ErrNotCorrectData) {
			return nil, apperr.ErrNotCorrectData
		}

		return nil, fmt.Errorf("login request failed: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("response", resp).Msg("finished auth client request")

	err = s.repo.InsertUser(ctx, query.InsertUserParams{
		Guid:       resp.UserGUID,
		Occupation: occupation,
		Name:       name,
		Email:      email,
	})
	if err != nil {
		return nil, fmt.Errorf("inserting user: %w", err)
	}

	return &models.DefaultStatusResponse{
			Code:    "01",
			Message: "Successfully created",
		},
		nil
}
