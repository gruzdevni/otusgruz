package order

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"otusgruz/internal/apperr"
	"otusgruz/internal/client/billhttp"
	"otusgruz/internal/models"
	query "otusgruz/internal/repo"
	"otusgruz/pkg/http"
)

type repo interface {
	CreateOrder(ctx context.Context, arg query.CreateOrderParams) error
}

type billClient interface {
	CreateUserRequest(ctx context.Context, guid uuid.UUID) error
	GetUserBalanceRequest(ctx context.Context, guid uuid.UUID) (decimal.Decimal, error)
	ChangeBalanceRequest(ctx context.Context, operType string, params billhttp.ChangeBalanceRequest) error
}

type service struct {
	repo       repo
	billClient billClient
}

type Service interface {
	ProcessOrder(ctx context.Context, params models.NewOrder) (*models.CreatedOrderData, error)
	// IncreaseBalance(ctx context.Context) (*models.DefaultStatusResponse, error)
	// GetBalance(ctx context.Context) (*models.DefaultStatusResponse, error)
}

func NewService(repo repo, billClient billClient) Service {
	return &service{
		repo:       repo,
		billClient: billClient,
	}
}

func (s *service) ProcessOrder(ctx context.Context, params models.NewOrder) (*models.CreatedOrderData, error) {
	ctxUserGUID := http.UserGUIDFromContext(ctx)

	userGUID, err := uuid.Parse(params.UserGUID.String())
	if err != nil {
		return nil, fmt.Errorf("parsing user guid: %w", err)
	}

	if ctxUserGUID == uuid.Nil || ctxUserGUID != userGUID {
		return nil, apperr.ErrNoPermission
	}

	orderstatus := query.OrderStatusDraft

	balance, err := s.billClient.GetUserBalanceRequest(ctx, userGUID)
	if err != nil {
		return nil, fmt.Errorf("getting user balance: %w", err)
	}

	orderAmount := decimal.NewFromFloat(params.Amount)

	orderNumber, err := strconv.Atoi(time.Now().Format("060102150405999"))
	if err != nil {
		return nil, fmt.Errorf("making order number: %w", err)
	}

	if balance.GreaterThanOrEqual(orderAmount) {
		err := s.billClient.ChangeBalanceRequest(ctx, billhttp.OutcomeType, billhttp.ChangeBalanceRequest{
			UserGUID:     userGUID,
			OperationRef: fmt.Sprintf("Заказ №%d", orderNumber),
			Amount:       orderAmount,
		})

		if err == nil {
			orderstatus = query.OrderStatusCompleted
		}
	}

	err = s.repo.CreateOrder(ctx, query.CreateOrderParams{
		Guid:     uuid.New(),
		UserGuid: userGUID,
		Number:   int32(orderNumber),
		Amount:   orderAmount,
		Status:   orderstatus,
	})

	return &models.CreatedOrderData{
			OrderNumber: "01",
			Status:      string(orderstatus),
		},
		nil
}
