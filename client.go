package hopsworks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	config    *ClientConfig
	projectID uint64
}

// NewClient creates new Hopsworks API client.
func NewClient(apiKey, project string) *Client {
	config := DefaultConfig(apiKey)
	config.Project = project

	return NewClientWithConfig(config)
}

// NewClientWithConfig creates new Hopsworks API client for specified config.
func NewClientWithConfig(config *ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

// Login connects to the Hopsworks server and returns project.
func (c *Client) Login(ctx context.Context) (*ProjectClient, error) {
	p, err := c.GetProject(ctx, c.config.Project)
	if err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}
	p.projectID = p.ID

	return p, nil
}

// DownloadDatasetFile downloads a file from the dataset to the specified local path.
func (c *Client) DownloadDatasetFile(ctx context.Context, remotePath, localPath string) error {
	url := c.url(
		"project",
		fmt.Sprintf("%d", c.projectID),
		"dataset",
		"download",
		"with_auth",
		remotePath,
	)
	queryArgs := map[string]string{
		"type": "DATASET",
	}

	req, err := c.newRequest(ctx, http.MethodGet, url, withQueryArgs(queryArgs))
	if err != nil {
		return err
	}

	f, err := os.OpenFile(localPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	v := new(bytes.Buffer)
	if err := c.sendRequest(req, &v); err != nil {
		return err
	}

	_, err = io.Copy(f, v)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

type requestOptions struct {
	body      any
	header    http.Header
	queryArgs map[string]string
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(opts *requestOptions) {
		opts.body = body
	}
}

func withContentType(contentType string) requestOption {
	return func(opts *requestOptions) {
		opts.header.Set("Content-Type", contentType)
	}
}

func withQueryArgs(args map[string]string) requestOption {
	return func(opts *requestOptions) {
		opts.queryArgs = args
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string, setters ...requestOption) (*http.Request, error) {
	// Default Options
	opts := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(opts)
	}

	var bodyReader io.Reader
	if opts.body != nil {
		if v, ok := opts.body.(io.Reader); ok {
			bodyReader = v
		} else {
			var reqBytes []byte
			reqBytes, err := json.Marshal(opts.body)
			if err != nil {
				return nil, fmt.Errorf("marshal: %w", err)
			}
			bodyReader = bytes.NewBuffer(reqBytes)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	c.setCommonHeaders(req)

	return req, nil
}

func (c *Client) sendRequest(req *http.Request, v any) error {
	req.Header.Set("Accept", "application/json")

	// Check whether Content-Type is already set, Upload Files API requires
	// Content-Type == multipart/form-data
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if isFailureStatusCode(res) {
		return c.handleErrorResp(res)
	}

	return decodeResponse(res.Body, v)
}

func (c *Client) setCommonHeaders(req *http.Request) {
	if c.config.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("ApiKey %s", c.config.apiKey))
	}
}

func (c *Client) url(path ...string) string {
	base := fmt.Sprintf("https://%s:%d/hopsworks-api/api", c.config.Host, c.config.Port)
	// Ignore error as base is validated.
	full, _ := url.JoinPath(base, path...)
	return full
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes *APIError
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil {
		reqErr := &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
		}
		return reqErr
	}

	errRes.HTTPStatusCode = resp.StatusCode
	return errRes
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		return json.NewDecoder(body).Decode(v)
	}
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}
