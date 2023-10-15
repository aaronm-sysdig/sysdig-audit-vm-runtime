package sysdighttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SysdigRequestConfig struct {
	Method     string
	URL        string
	Headers    map[string]string
	Params     map[string]interface{}
	JSON       interface{}
	Data       map[string]string
	Auth       [2]string
	Verify     bool
	Stream     bool
	MaxRetries int
	BaseDelay  int
	MaxDelay   int
	Timeout    int
}

func DefaultSysdigRequestConfig() SysdigRequestConfig {
	return SysdigRequestConfig{
		Method:     "GET",
		Verify:     true,
		MaxRetries: 5,
		BaseDelay:  5,
		MaxDelay:   60,
		Timeout:    10,
	}
}

func SysdigRequest(SysdigRequest SysdigRequestConfig) (*http.Response, error) {
	retries := 0
	var resp *http.Response
	var err error

	for retries <= SysdigRequest.MaxRetries {
		var jsonData io.Reader
		if SysdigRequest.JSON != nil {
			byteData, err := json.Marshal(SysdigRequest.JSON)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
			}
			jsonData = bytes.NewBuffer(byteData)
		}

		u, err := url.Parse(SysdigRequest.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}

		params := url.Values{}
		for k, v := range SysdigRequest.Params {
			switch value := v.(type) {
			case int:
				params.Add(k, strconv.Itoa(value))
			case string:
				params.Add(k, value)
			// You can add more cases if needed
			default:
				// Handle unexpected types if necessary, or ignore them
			}
		}

		u.RawQuery = params.Encode()

		req, err := http.NewRequest(SysdigRequest.Method, u.String(), jsonData)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		for k, v := range SysdigRequest.Headers {
			req.Header.Set(k, v)
		}

		if len(SysdigRequest.Auth) == 2 {
			if SysdigRequest.Auth[0] != "" {
				req.SetBasicAuth(SysdigRequest.Auth[0], SysdigRequest.Auth[1])
			}
		}
		client := &http.Client{
			Timeout: time.Duration(SysdigRequest.Timeout) * time.Second,
		}
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if err == nil {
			// log status code if request did not fail
			log.Printf("Received HTTP status code: %d", resp.StatusCode)
			respBody, _ := io.ReadAll(resp.Body)
			log.Printf("Response body: %s", string(respBody))
			resp.Body.Close() // ensure response body is closed
		}

		log.Printf("Error: %v. Retrying in %d seconds...", err, SysdigRequest.BaseDelay)
		log.Printf("Retry %d, Sleeping for %d seconds", retries, SysdigRequest.BaseDelay)
		time.Sleep(time.Duration(SysdigRequest.BaseDelay) * time.Second)
		retries++
	}

	log.Printf("Failed to fetch data from %s after %d retries.", SysdigRequest.URL, SysdigRequest.MaxRetries)
	log.Printf("Error making request to %s: %v", SysdigRequest.URL, err)

	// Manually create an HTTP response with a 503 status code
	resp = &http.Response{
		Status:     "503 Service Unavailable",
		StatusCode: http.StatusServiceUnavailable,
		Body:       io.NopCloser(bytes.NewBufferString("Service is unavailable after retries.")),
	}
	return resp, fmt.Errorf("failed after %d retries", SysdigRequest.MaxRetries)
}

func ResponseBodyToJson(resp *http.Response, target interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}

	return nil
}
