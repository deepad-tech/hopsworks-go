package hopsworks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Response interface {
	SetHeader(http.Header)
}

type httpHeader http.Header

func (h *httpHeader) SetHeader(header http.Header) {
	*h = httpHeader(header)
}

func (h *httpHeader) Header() http.Header {
	return http.Header(*h)
}

type Client struct {
	config *ClientConfig
}

// NewClient creates new Hopsworks API client.
func NewClient(apiKey string) *Client {
	return NewClientWithConfig(DefaultConfig(apiKey))
}

// NewClientWithConfig creates new Hopsworks API client for specified config.
func NewClientWithConfig(config *ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

// Login connects to the Hopsworks server and returns project.
func (c *Client) Login(ctx context.Context) (*Project, error) {
	return c.GetProject(ctx, c.config.Project)
}

type requestOptions struct {
	body   any
	header http.Header
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

func withContentType(contentType string) requestOption {
	return func(args *requestOptions) {
		args.header.Set("Content-Type", contentType)
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string, setters ...requestOption) (*http.Request, error) {
	// Default Options
	args := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(args)
	}

	var bodyReader io.Reader
	if args.body != nil {
		if v, ok := args.body.(io.Reader); ok {
			bodyReader = v
		} else {
			var reqBytes []byte
			reqBytes, err := json.Marshal(args.body)
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

func (c *Client) sendRequest(req *http.Request, v Response) error {
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

	if v != nil {
		v.SetHeader(res.Header)
	}

	return decodeResponse(res.Body, v)
}

func (c *Client) setCommonHeaders(req *http.Request) {
	if c.config.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.apiKey))
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
	var errRes ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
		}
		if errRes.Error != nil {
			reqErr.Err = errRes.Error
		}
		return reqErr
	}

	errRes.Error.HTTPStatusCode = resp.StatusCode
	return errRes.Error
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
