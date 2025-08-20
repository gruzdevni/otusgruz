package authhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	internalclient "otusgruz/internal/client"
)

const loginEndpoint = "api/login"

type Client interface {
	LoginRequest(ctx context.Context, email string, password string) (LoginResponse, error)
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
func (c *client) LoginRequest(ctx context.Context, email string, password string) (LoginResponse, error) {
	type login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	requestBody := login{Email: email, Password: password}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("marshal request body: %w", err)
	}

	path, err := url.JoinPath(c.baseURL, loginEndpoint)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("joining path: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		path,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, code, err := internalclient.DoJSON[withErr[LoginResponse]](c.doer, req)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("login request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		return LoginResponse{
			Status: resp.Status,
			Errors: resp.Errors,
		}, fmt.Errorf("not ok status for login request: %d, errors %v, path: %s, code: %d", resp.Status, resp.Errors, req.URL.Path, code)
	}

	return LoginResponse{
		Status:   resp.Status,
		UserGUID: resp.Result.UserGUID,
		Errors:   resp.Errors,
	}, nil
}
