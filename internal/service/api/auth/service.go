package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"

	authhttp "otusgruz/internal/client/authhttp"
	query "otusgruz/internal/repo"
)

var (
	ErrNotCorrectData   = errors.New("Not correct password or email")
	ErrNoSuchUser       = errors.New("No such email registered")
	ErrEmailAlreadyUsed = errors.New("Email is already registered. Please login")
)

type repo interface {
	GetUserByEmail(ctx context.Context, email string) (query.User, error)
	IsAuth(ctx context.Context, guid uuid.UUID) (query.LoggedIn, error)
	InsertSession(ctx context.Context, userGuid uuid.UUID) error
	InsertUser(ctx context.Context, arg query.InsertUserParams) error
}

type authClient interface {
	LoginRequest(ctx context.Context, email string, password string) (authhttp.LoginResponse, error)
}

type service struct {
	repo       repo
	authClient authClient
}

type Service interface {
	Auth(ctx context.Context, guid uuid.UUID) (uuid.UUID, error)
	Login(ctx context.Context, email string, pwd string) (uuid.UUID, error)
	Singup(ctx context.Context, email string, pwd string) error
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
	resp, err := s.authClient.LoginRequest(ctx, email, pwd)
	if err != nil {
		return uuid.Nil, fmt.Errorf("login request: %w", err)
	}

	return uuid.UUID(resp.UserGUID), nil
}

func (s *service) Singup(ctx context.Context, email string, pwd string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("getting user by email: %w", err)
	}

	if !lo.IsEmpty(user) {
		return ErrEmailAlreadyUsed
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), 8)
	if err != nil {
		return fmt.Errorf("encrypting password: %w", err)
	}

	err = s.repo.InsertUser(ctx, query.InsertUserParams{
		Guid:       uuid.New(),
		Occupation: "",
		Name:       "",
		Email:      email,
		Pwd:        string(hashedPwd),
	})
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}
