package authhttp

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
	loginEndpoint  = "api/login"
	signupEndpoint = "api/signup"
)

type Client interface {
	LoginRequest(ctx context.Context, email string, password string) (LoginResponse, error)
	SignupRequest(ctx context.Context, email string, password string) (SignupResponse, error)
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

	zerolog.Ctx(ctx).Info().Any("jsonBody of login request", string(jsonBody)).Msg("prepared login body")

	path, err := url.JoinPath(c.baseURL, loginEndpoint)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of login request", path).Msg("prepared login path")

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

	resp, code, err := internalclient.Do(c.doer, req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed login request")
		return LoginResponse{}, fmt.Errorf("login request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		if code == http.StatusUnauthorized {
			return LoginResponse{}, apperr.ErrNotCorrectData
		}
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed login request")
		return LoginResponse{}, fmt.Errorf("not ok status for login request: path: %s, code: %d", req.URL.Path, code)
	}

	var res LoginResponse

	if err := json.Unmarshal(resp, &res); err != nil {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("unmarshalling json")
		return res, fmt.Errorf("unmarshalling json: %s", err)
	}

	return res, nil
}

func (c *client) SignupRequest(ctx context.Context, email string, password string) (SignupResponse, error) {
	type signup struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	requestBody := signup{Email: email, Password: password}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return SignupResponse{}, fmt.Errorf("marshal request body: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("jsonBody of Signup request", string(jsonBody)).Msg("prepared Signup body")

	path, err := url.JoinPath(c.baseURL, signupEndpoint)
	if err != nil {
		return SignupResponse{}, fmt.Errorf("joining path: %w", err)
	}

	zerolog.Ctx(ctx).Info().Any("path of Signup request", path).Msg("prepared Signup path")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		path,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return SignupResponse{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, code, err := internalclient.Do(c.doer, req)
	if err != nil {
		zerolog.Ctx(ctx).Err(err).Any("response", resp).Msg("failed Signup request")
		return SignupResponse{}, fmt.Errorf("Signup request error: %w, path: %s, code: %d", err, req.URL.Path, code)
	}

	if code != http.StatusOK {
		if code == http.StatusUnauthorized {
			return SignupResponse{}, apperr.ErrNotCorrectData
		}
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("failed Signup request")
		return SignupResponse{}, fmt.Errorf("not ok status for Signup request: path: %s, code: %d", req.URL.Path, code)
	}

	var res SignupResponse

	if err := json.Unmarshal(resp, &res); err != nil {
		zerolog.Ctx(ctx).Info().Any("response", resp).Any("code", code).Msg("unmarshalling json")
		return res, fmt.Errorf("unmarshalling json: %s", err)
	}

	return res, nil
}
