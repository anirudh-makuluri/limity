package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Store struct {
	baseURL string
	token   string
	client  *http.Client
}

type response struct {
	Result interface{} `json:"result"`
	Error  *string     `json:"error"`
}

func NewStore(baseURL, token string) *Store {
	return &Store{baseURL: baseURL, token: token, client: http.DefaultClient}
}

func (s *Store) Incr(ctx context.Context, key string) (int, error) {
	if s.baseURL == "" || s.token == "" {
		return 0, fmt.Errorf("missing redis configuration")
	}

	url := fmt.Sprintf("%s/incr/%s", s.baseURL, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var redisResp response
	if err := json.Unmarshal(body, &redisResp); err != nil {
		return 0, err
	}
	if redisResp.Error != nil {
		return 0, fmt.Errorf("redis error: %s", *redisResp.Error)
	}

	count, ok := redisResp.Result.(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected redis response type")
	}

	return int(count), nil
}

func (s *Store) Expire(ctx context.Context, key string, seconds int) error {
	if s.baseURL == "" || s.token == "" {
		return fmt.Errorf("missing redis configuration")
	}

	url := fmt.Sprintf("%s/expire/%s/%d", s.baseURL, key, seconds)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var redisResp response
	if err := json.Unmarshal(body, &redisResp); err != nil {
		return err
	}
	if redisResp.Error != nil {
		return fmt.Errorf("redis error: %s", *redisResp.Error)
	}

	return nil
}
