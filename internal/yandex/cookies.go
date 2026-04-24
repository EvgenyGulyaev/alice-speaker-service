package yandex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	tokenBySessionIDURL          = "https://mobileproxy.passport.yandex.net/1/bundle/oauth/token_by_sessionid"
	tokenBySessionIDClientID     = "c0ebe342af7d48fbbbfcf2d2eedb8f9e"
	tokenBySessionIDClientSecret = "ad0a908f0aa341a182a37ecd75bc319e"
)

type exportedCookie struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

type tokenBySessionResponse struct {
	AccessToken string   `json:"access_token"`
	Errors      []string `json:"errors"`
	Message     string   `json:"message"`
}

func (c *Client) ExtractXTokenFromCookies(raw string) (string, error) {
	cookieHeader, host, err := normalizeExportedCookies(raw)
	if err != nil {
		return "", err
	}

	form := map[string]string{
		"client_id":     tokenBySessionIDClientID,
		"client_secret": tokenBySessionIDClientSecret,
	}
	values := url.Values{}
	for key, value := range form {
		values.Set(key, value)
	}

	request, err := http.NewRequest(http.MethodPost, tokenBySessionIDURL, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Ya-Client-Host", host)
	request.Header.Set("Ya-Client-Cookie", cookieHeader)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var payload tokenBySessionResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}
	if response.StatusCode >= http.StatusBadRequest {
		if payload.Message != "" {
			return "", errors.New(payload.Message)
		}
		if len(payload.Errors) > 0 {
			return "", errors.New(strings.Join(payload.Errors, ", "))
		}
		return "", fmt.Errorf("yandex token_by_sessionid failed with status %d", response.StatusCode)
	}
	if strings.TrimSpace(payload.AccessToken) == "" {
		return "", errors.New("yandex did not return x_token from cookies")
	}
	return strings.TrimSpace(payload.AccessToken), nil
}

func normalizeExportedCookies(raw string) (cookieHeader string, host string, err error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", "", errors.New("cookies payload is empty")
	}

	if strings.HasPrefix(trimmed, "[") {
		var cookies []exportedCookie
		if err := json.Unmarshal([]byte(trimmed), &cookies); err != nil {
			return "", "", err
		}
		if len(cookies) == 0 {
			return "", "", errors.New("cookies list is empty")
		}

		pairs := make([]string, 0, len(cookies))
		host = "passport.yandex.ru"
		for _, cookie := range cookies {
			name := strings.TrimSpace(cookie.Name)
			value := strings.TrimSpace(cookie.Value)
			if name == "" || value == "" {
				continue
			}
			pairs = append(pairs, name+"="+value)
			domain := strings.TrimSpace(cookie.Domain)
			if strings.HasPrefix(domain, ".yandex.") {
				host = domain
			}
		}
		if len(pairs) == 0 {
			return "", "", errors.New("cookies list does not contain name/value pairs")
		}
		return strings.Join(pairs, "; "), host, nil
	}

	return trimmed, "passport.yandex.ru", nil
}
