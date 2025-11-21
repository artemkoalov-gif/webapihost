package http_server

import (
	"context"
	"net/http"
	"time"
)

func httpRequest(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+url, nil)
	if err != nil {
		return StatusNotOK, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return StatusNotOK, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return StatusNotOK, nil
	}

	return StatusOK, nil
}
