package user

import (
	"context"
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	"otusgruz/internal/models"
	query "otusgruz/internal/repo"
)

type repo interface {
	GetUser(ctx context.Context, guid uuid.UUID) (query.User, error)
	DeleteUser(ctx context.Context, guid uuid.UUID) error
	InsertUser(ctx context.Context, arg query.InsertUserParams) error
	UpdateUser(ctx context.Context, arg query.UpdateUserParams) error
}

type service struct {
	repo repo
}

type Service interface {
	GetUser(ctx context.Context, guid uuid.UUID) (*models.UserData, error)
}

func NewService(repo repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUser(ctx context.Context, guid uuid.UUID) (*models.UserData, error) {
	res, err := s.repo.GetUser(ctx, guid)
	if err != nil {
		return nil, fmt.Errorf("getting user: %w", err)
	}

	return &models.UserData{
		GUID:       strfmt.UUID(res.Guid.String()),
		IsDeleted:  res.IsDeleted,
		Name:       res.Name,
		Occupation: res.Occupation,
	}, nil
}
