package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var esClient = &http.Client{
	Timeout: 10 * time.Second,
}

func ESEnabled() bool {
	return C.ESURL != ""
}

func ESRequest(method, path string, body any) (map[string]any, int, error) {
	var r io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}

		r = bytes.NewReader(b)
	}

	req, err := http.NewRequest(
		method,
		C.ESURL+path,
		r,
	)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	if C.ESKey != "" {
		req.Header.Set(
			"Authorization",
			"ApiKey "+C.ESKey,
		)
	}

	res, err := esClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	var out map[string]any

	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, res.StatusCode, fmt.Errorf(
				"invalid Elasticsearch response: %w; body: %s",
				err,
				string(raw),
			)
		}
	}

	if res.StatusCode >= 300 {
		return out, res.StatusCode, fmt.Errorf(
			"Elasticsearch returned %d: %s",
			res.StatusCode,
			string(raw),
		)
	}

	return out, res.StatusCode, nil
}

func ConnectElasticsearch() error {
	if !ESEnabled() {
		return errors.New(
			"ELASTICSEARCH_URL not set; search is disabled",
		)
	}

	_, code, err := ESRequest("GET", "/", nil)

	if err != nil {
		return err
	}

	if code == 401 || code == 403 {
		return errors.New(
			"credentials rejected; check ELASTICSEARCH_API_KEY",
		)
	}

	if code >= 300 {
		return fmt.Errorf(
			"unexpected status %d",
			code,
		)
	}

	return nil
}