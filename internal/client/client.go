package internalclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func DoJSON[T any](d Doer, req *http.Request) (T, int, error) {
	var val T

	body, code, err := Do(d, req)
	if err != nil {
		return val, code, err
	}

	if err := json.Unmarshal(body, &val); err != nil {
		return val, code, fmt.Errorf("unmarshalling json: %w", err)
	}

	return val, code, nil
}

func Do(d Doer, req *http.Request) ([]byte, int, error) {
	resp, err := d.Do(req)
	if err != nil {
		code := 0

		if resp != nil && resp.Body != nil {
			code = resp.StatusCode
			defer resp.Body.Close() //nolint:errcheck
		}

		return nil, code, fmt.Errorf("making request: %w", err)
	}

	defer resp.Body.Close() //nolint:errcheck

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading body: %w", err)
	}

	return body, resp.StatusCode, nil
}
