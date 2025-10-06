package deliveryhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	internalclient "otusgruz/internal/client"

	"github.com/rs/zerolog"
)

const (
	reserveDeliveryEndpoint      = "api/reserve-slot"
	unreserveDeliveryEndpoint    = "api/unreserve-slot"
	checkDeliveryReserveEndpoint = "api/check-delivery-status"
)

type Client interface {
	CheckDeliveryReserve(ctx context.Context, orderID string) (DeliveryStatus, error)
	ReserveDelivery(ctx context.Context, orderID string, slotID string) error
	UnreserveDelivery(ctx context.Context, orderID string) error
}

type client struct {
	doer    internalclient.Doer
	baseURL string
}

func NewClient(baseURL string, doer internalclient.Doer) Client {
	return &client{doer: doer, baseURL: baseURL}
}

type Error struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type withErr[T any] struct {
	Result T       `json:"result"`
	Errors []Error `json:"errors"`
	Status int     `json:"status"`
}

//nolint:exhaustruct
func (c *client) CheckDeliveryReserve(ctx context.Context, orderID string) (DeliveryStatus, error) {
	path, err := url.JoinPath(c.baseURL, checkDeliveryReserveEndpoint, orderID)
	if err != nil {
		return DeliveryStatus{}, fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of checkDeliveryReserve request", path).Msg("prepared checkDeliveryReserve path")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		path,
		nil,
	)
	if err != nil {
		return DeliveryStatus{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, code, err := internalclient.Do(c.doer, req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed checkDeliveryReserve request")
		return DeliveryStatus{}, fmt.Errorf("checkDeliveryReserve error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed checkDeliveryReserve request")
		return DeliveryStatus{}, fmt.Errorf("not ok status for checkDeliveryReserve request: path: %s, code: %d", req.URL.Path, code)
	}

	var res DeliveryStatus

	if err := json.Unmarshal(resp, &res); err != nil {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("unmarshalling json")
		return DeliveryStatus{}, fmt.Errorf("unmarshalling json: %s", err)
	}

	return res, nil
}

func (c *client) ReserveDelivery(ctx context.Context, orderID string, slotID string) error {
	params := ReserveDelivery{
		OrderID: orderID,
		SlotID:  slotID,
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("jsonBody of ReserveDelivery request", string(jsonBody)).Msg("prepared ReserveDelivery body")

	path, err := url.JoinPath(c.baseURL, reserveDeliveryEndpoint)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of ReserveDelivery request", path).Msg("prepared ReserveDelivery path")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		path,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, code, err := internalclient.Do(c.doer, req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed ReserveDeliveryrequest")
		return fmt.Errorf("ReserveDelivery request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed ReserveDelivery request")
		return fmt.Errorf("not ok status for ReserveDelivery request: path: %s, code: %d", req.URL.Path, code)
	}

	return nil
}

//nolint:exhaustruct
func (c *client) UnreserveDelivery(ctx context.Context, orderID string) error {
	path, err := url.JoinPath(c.baseURL, unreserveDeliveryEndpoint, orderID)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of UnreserveDelivery request", path).Msg("prepared UnreserveDelivery path")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		path,
		nil,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, code, err := internalclient.Do(c.doer, req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed UnreserveDelivery request")
		return fmt.Errorf("UnreserveDelivery error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed UnreserveDelivery request")
		return fmt.Errorf("not ok status for UnreserveDelivery request: path: %s, code: %d", req.URL.Path, code)
	}

	return nil
}
