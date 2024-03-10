package hopsworks

import (
	"context"
	"fmt"
	"net/http"
)

type Dataset struct{}

// Item represents a file or directory in a dataset
type Item struct {
	Attributes map[string]interface{} `json:"attributes"`
	Path       string                 `json:"path"`
}

func (c *Client) ListDataset(ctx context.Context, remotePath string) ([]Item, error) {
	url := c.url("project", fmt.Sprintf("%d", c.projectID), "dataset", remotePath)

	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var v []Item
	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	return v, nil
}
