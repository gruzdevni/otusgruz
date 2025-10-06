package notifyhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"otusgruz/internal/apperr"
	internalclient "otusgruz/internal/client"

	"github.com/rs/zerolog"
)

const (
	SuccessType = "success"
	FailureType = "failure"

	userNotificationCreateEndpoint = "api/user-notification/create"
)

type Client interface {
	CreateNotificationRequest(ctx context.Context, params CreateNotificationRequest) error
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

func (c *client) CreateNotificationRequest(ctx context.Context, params CreateNotificationRequest) error {
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("jsonBody of new notification request", string(jsonBody)).Msg("prepared new notification body")

	path, err := url.JoinPath(c.baseURL, userNotificationCreateEndpoint)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of new notification request", path).Msg("prepared new notification path")

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
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed new notification request")
		return fmt.Errorf("new notification request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		if code == http.StatusForbidden {
			return apperr.ErrNotEnoughFunds
		}
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed new notification request")
		return fmt.Errorf("not ok status for new notification request: path: %s, code: %d", req.URL.Path, code)
	}

	return nil
}
