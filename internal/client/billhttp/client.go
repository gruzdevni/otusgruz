package billhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"otusgruz/internal/apperr"
	internalclient "otusgruz/internal/client"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

const (
	IncomeType  = "income"
	OutcomeType = "outcome"

	userCreateEndpoint      = "api/user"
	userBalanceEndpoint     = "api/user-balance"
	increaseBalanceEndpoint = "api/user-balance/increase"
	reduceBalanceEndpoint   = "api/user-balance/reduce"
)

type Client interface {
	CreateUserRequest(ctx context.Context, guid uuid.UUID) error
	GetUserBalanceRequest(ctx context.Context, guid uuid.UUID) (decimal.Decimal, error)
	ChangeBalanceRequest(ctx context.Context, operType string, params ChangeBalanceRequest) error
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
func (c *client) GetUserBalanceRequest(ctx context.Context, guid uuid.UUID) (decimal.Decimal, error) {
	path, err := url.JoinPath(c.baseURL, userBalanceEndpoint, guid.String())
	if err != nil {
		return decimal.Zero, fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of login request", path).Msg("prepared login path")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		path,
		nil,
	)
	if err != nil {
		return decimal.Zero, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, code, err := internalclient.Do(c.doer, req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed login request")
		return decimal.Zero, fmt.Errorf("login request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed login request")
		return decimal.Zero, fmt.Errorf("not ok status for login request: path: %s, code: %d", req.URL.Path, code)
	}

	var res UserBalanceResponse

	if err := json.Unmarshal(resp, &res); err != nil {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("unmarshalling json")
		return res.Amount, fmt.Errorf("unmarshalling json: %s", err)
	}

	return res.Amount, nil
}

func (c *client) ChangeBalanceRequest(ctx context.Context, operType string, params ChangeBalanceRequest) error {
	endpoint := reduceBalanceEndpoint

	if operType == IncomeType {
		endpoint = increaseBalanceEndpoint
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("jsonBody of change balance request", string(jsonBody)).Msg("prepared change balance body")

	path, err := url.JoinPath(c.baseURL, endpoint)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of change balance request", path).Msg("prepared change balance path")

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
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed change balance request")
		return fmt.Errorf("change balance request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		if code == http.StatusForbidden {
			return apperr.ErrNotEnoughFunds
		}
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed Signup request")
		return fmt.Errorf("not ok status for change balance request: path: %s, code: %d", req.URL.Path, code)
	}

	return nil
}

func (c *client) CreateUserRequest(ctx context.Context, guid uuid.UUID) error {
	type signup struct {
		Guid uuid.UUID `json:"guid"`
	}

	requestBody := signup{Guid: guid}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("jsonBody of Signup request", string(jsonBody)).Msg("prepared Signup body")

	path, err := url.JoinPath(c.baseURL, userCreateEndpoint)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of Signup request", path).Msg("prepared Signup path")

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
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed Signup request")
		return fmt.Errorf("Signup request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		if code == http.StatusUnauthorized {
			return apperr.ErrNotCorrectData
		}
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed Signup request")
		return fmt.Errorf("not ok status for Signup request: path: %s, code: %d", req.URL.Path, code)
	}

	return nil
}
