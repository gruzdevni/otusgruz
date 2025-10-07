package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	"otusgruz/internal/apperr"
	"otusgruz/internal/client/billhttp"
	"otusgruz/internal/client/goodshttp"
	"otusgruz/internal/client/notifyhttp"
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

type notifyClient interface {
	CreateNotificationRequest(ctx context.Context, params notifyhttp.CreateNotificationRequest) error
}

type deliveryClient interface {
	ReserveDelivery(ctx context.Context, orderID string, slotID string) error
	UnreserveDelivery(ctx context.Context, orderID string) error
}

type goodsClient interface {
	ReserveGoods(ctx context.Context, orderID string, goods []goodshttp.Goods) error
	UnreserveGoods(ctx context.Context, orderID string) error
}

type service struct {
	repo           repo
	billClient     billClient
	notifyClient   notifyClient
	deliveryClient deliveryClient
	goodsClient    goodsClient
}

type Service interface {
	ProcessOrder(ctx context.Context, params models.NewOrder) (*models.CreatedOrderData, error)
	// IncreaseBalance(ctx context.Context) (*models.DefaultStatusResponse, error)
	// GetBalance(ctx context.Context) (*models.DefaultStatusResponse, error)
}

func NewService(repo repo, billClient billClient, notifyClient notifyClient, deliveryClient deliveryClient, goodsClient goodsClient) Service {
	return &service{
		repo:           repo,
		billClient:     billClient,
		notifyClient:   notifyClient,
		deliveryClient: deliveryClient,
		goodsClient:    goodsClient,
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

	balance, err := s.billClient.GetUserBalanceRequest(ctx, userGUID)
	if err != nil {
		return nil, fmt.Errorf("getting user balance: %w", err)
	}

	orderstatus := query.OrderStatusDraft
	notifyType := notifyhttp.FailureType
	orderAmount := decimal.NewFromFloat(params.Amount)
	orderNumber := time.Now().Format("060102150405")
	operationRef := fmt.Sprintf("Заказ №%s", orderNumber)
	failReason := ""
	isNeedToStopProccess := false

	if err = s.deliveryClient.ReserveDelivery(ctx, orderNumber, params.DeliverySlot); err != nil {
		failReason = "delivery reserve problem"
		isNeedToStopProccess = true
	}

	if !isNeedToStopProccess {
		if err = s.goodsClient.ReserveGoods(ctx, orderNumber, lo.Map(params.Goods, func(item *models.NewOrderGoodsItems0, _ int) goodshttp.Goods {
			return goodshttp.Goods{
				Nomenclature: item.Nomenclature,
				Quantity:     int(item.Quantity),
			}
		})); err != nil {
			err = s.deliveryClient.UnreserveDelivery(ctx, orderNumber)
			if err != nil {
				return nil, fmt.Errorf("unreserving delivery slot: %w", err)
			}

			failReason = "goods reserve problem"
			isNeedToStopProccess = true
		}
	}

	if !isNeedToStopProccess {
		if balance.GreaterThanOrEqual(orderAmount) {
			err := s.billClient.ChangeBalanceRequest(ctx, billhttp.OutcomeType, billhttp.ChangeBalanceRequest{
				UserGUID:     userGUID,
				OperationRef: operationRef,
				Amount:       orderAmount.InexactFloat64(),
			})

			if err == nil {
				orderstatus = query.OrderStatusCompleted
				notifyType = notifyhttp.SuccessType
			} else {
				err = s.deliveryClient.UnreserveDelivery(ctx, orderNumber)
				if err != nil {
					return nil, fmt.Errorf("unreserving delivery slot: %w", err)
				}

				err = s.goodsClient.UnreserveGoods(ctx, orderNumber)
				if err != nil {
					return nil, fmt.Errorf("unreserving goods: %w", err)
				}

				failReason = "billing problem"
			}
		} else {
			err = s.deliveryClient.UnreserveDelivery(ctx, orderNumber)
			if err != nil {
				return nil, fmt.Errorf("unreserving delivery slot: %w", err)
			}

			err = s.goodsClient.UnreserveGoods(ctx, orderNumber)
			if err != nil {
				return nil, fmt.Errorf("unreserving goods: %w", err)
			}

			failReason = "not enough balance"
		}
	}

	err = s.repo.CreateOrder(ctx, query.CreateOrderParams{
		Guid:     uuid.New(),
		UserGuid: userGUID,
		Number:   orderNumber,
		Amount:   orderAmount,
		Status:   orderstatus,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new order: %w", err)
	}

	err = s.notifyClient.CreateNotificationRequest(ctx, notifyhttp.CreateNotificationRequest{
		UserGUID:     userGUID,
		OperationRef: operationRef,
		NotifyType:   notifyType,
	})
	if err != nil {
		return nil, fmt.Errorf("creating new notification: %w", err)
	}

	return &models.CreatedOrderData{
			OrderNumber: orderNumber,
			Status:      string(orderstatus),
			FailReason:  failReason,
		},
		nil
}
