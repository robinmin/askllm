package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/robinmin/askllm/pkg/utils/log"
)

var client *retryablehttp.Client

func init() {
	client = retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 30 * time.Second

	// Set timeout on the underlying http.Client
	client.HTTPClient.Timeout = 60 * time.Second

	// Custom retry policy
	client.CheckRetry = customRetryPolicy
	client.Logger = log.GetDefaultLogger()
}

func customRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// First, check the default retry policy
	shouldRetry, checkErr := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

	// If there's an error in checking, return it
	if checkErr != nil {
		return false, checkErr
	}

	// If the default policy says to retry, or there's no response, return that decision
	if shouldRetry || resp == nil {
		return shouldRetry, nil
	}

	// Add custom logic for specific status codes
	switch resp.StatusCode {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		// We're returning true to indicate that we should retry, but we're not sleeping here.
		// The retryablehttp client will handle the delay between retries based on its configuration.
		return true, nil
	}

	// For all other cases, return the default retry decision
	return shouldRetry, nil
}

func APIRequestCore(method string, url string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := retryablehttp.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	log.Infof("[API] ====> : %s %s", method, url)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Infof("[API] <==== : %s %s - %d", method, url, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	log.Infof("[API] Response: %s", string(responseBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return responseBody, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return responseBody, nil
}

func APIPost[request any, response any](url string, body request, headers map[string]string) (*response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Errorf("error making request: %v", err)
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	responseBody, err := APIRequestCore(http.MethodPost, url, jsonBody, headers)
	if err != nil {
		log.Errorf("error making request: %v", err)
		return nil, err
	}

	var result response
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		log.Errorf("error making request: %v", err)
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

func APIGet[response any](url string, headers map[string]string) (*response, error) {
	responseBody, err := APIRequestCore(http.MethodGet, url, nil, headers)
	if err != nil {
		log.Errorf("error making request: %v", err)
		return nil, err
	}

	var result response
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		log.Errorf("error making request: %v", err)
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}
