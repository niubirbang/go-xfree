package goxfree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type api struct {
	client *http.Client
}

func (a *api) url(path string, query url.Values) string {
	var q string
	if query != nil {
		q = query.Encode()
	}
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("http://unix/%s?%s", path, q)
}
func (a *api) do(req *http.Request) ([]byte, error) {
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}
	var data struct {
		Message string `json:"message"`
	}
	json.Unmarshal(body, &data)
	if data.Message != "" {
		return nil, errors.New(data.Message)
	}
	return nil, fmt.Errorf("status code: %d", resp.StatusCode)
}
func (a *api) get(path string, query url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, a.url(path, query), nil)
	if err != nil {
		return nil, err
	}
	return a.do(req)
}
func (a *api) post(path string, query url.Values, data interface{}) ([]byte, error) {
	dataBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, a.url(path, query), bytes.NewReader(dataBody))
	if err != nil {
		return nil, err
	}
	return a.do(req)
}
func (a *api) put(path string, query url.Values, param interface{}) ([]byte, error) {
	paramBody, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, a.url(path, query), bytes.NewReader(paramBody))
	if err != nil {
		return nil, err
	}
	return a.do(req)
}
func (a *api) patch(path string, query url.Values, data interface{}) ([]byte, error) {
	dataBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPatch, a.url(path, query), bytes.NewReader(dataBody))
	if err != nil {
		return nil, err
	}
	return a.do(req)
}
func (a *api) delete(path string, query url.Values, data interface{}) ([]byte, error) {
	dataBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodDelete, a.url(path, query), bytes.NewReader(dataBody))
	if err != nil {
		return nil, err
	}
	return a.do(req)
}
