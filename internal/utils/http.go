package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

var client = &http.Client{}

func SendJSONRequest(ctx context.Context, method, url string, data interface{}) (map[string]interface{}, error) {
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

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func SendRequest(ctx context.Context, method, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://kaspi.kz/")
	return client.Do(req)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
