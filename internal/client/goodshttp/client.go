package goodshttp

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
	reserveGoodsEndpoint      = "api/reserve-order-goods"
	unreserveGoodsEndpoint    = "api/unreserve-goods"
	checkGoodsReserveEndpoint = "api/check-reserve-status"
)

type Client interface {
	CheckGoodsReserve(ctx context.Context, orderID string) (ReserveStatus, error)
	ReserveGoods(ctx context.Context, orderID string, goods []Goods) error
	UnreserveGoods(ctx context.Context, orderID string) error
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
func (c *client) CheckGoodsReserve(ctx context.Context, orderID string) (ReserveStatus, error) {
	path, err := url.JoinPath(c.baseURL, checkGoodsReserveEndpoint, orderID)
	if err != nil {
		return ReserveStatus{}, fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of CheckGoodsReserve request", path).Msg("prepared CheckGoodsReserve path")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		path,
		nil,
	)
	if err != nil {
		return ReserveStatus{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, code, err := internalclient.Do(c.doer, req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed CheckGoodsReserve request")
		return ReserveStatus{}, fmt.Errorf("CheckGoodsReserve error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed CheckGoodsReserve request")
		return ReserveStatus{}, fmt.Errorf("not ok status for CheckGoodsReserve request: path: %s, code: %d", req.URL.Path, code)
	}

	var res ReserveStatus

	if err := json.Unmarshal(resp, &res); err != nil {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("unmarshalling json")
		return ReserveStatus{}, fmt.Errorf("unmarshalling json: %s", err)
	}

	return res, nil
}

func (c *client) ReserveGoods(ctx context.Context, orderID string, goods []Goods) error {
	params := ReserveGoods{
		OrderID: orderID,
		Goods:   goods,
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("jsonBody of ReserveGoods request", string(jsonBody)).Msg("prepared ReserveGoods body")

	path, err := url.JoinPath(c.baseURL, reserveGoodsEndpoint)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of ReserveGoods request", path).Msg("prepared ReserveGoods path")

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
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed ReserveGoods")
		return fmt.Errorf("ReserveGoods request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed ReserveGoods request")
		return fmt.Errorf("not ok status for ReserveGoods request: path: %s, code: %d", req.URL.Path, code)
	}

	return nil
}

//nolint:exhaustruct
func (c *client) UnreserveGoods(ctx context.Context, orderID string) error {
	path, err := url.JoinPath(c.baseURL, unreserveGoodsEndpoint, orderID)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of UnreserveGoods request", path).Msg("prepared UnreserveGoods path")

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
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed UnreserveGoods request")
		return fmt.Errorf("UnreserveGoods error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed UnreserveGoods request")
		return fmt.Errorf("not ok status for UnreserveGoods request: path: %s, code: %d", req.URL.Path, code)
	}

	return nil
}
