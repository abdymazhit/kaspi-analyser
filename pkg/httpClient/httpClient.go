package httpClient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Service struct {
	client *http.Client
}

func NewService() *Service {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100
	return &Service{
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		},
	}
}

func (s Service) SendJSONRequest(ctx context.Context, method, url string, data interface{}) (map[string]interface{}, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://kaspi.kz/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}
