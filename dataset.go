package hopsworks

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) ListDataset(ctx context.Context, remotePath string) ([]string, error) {
	url := c.url("project", fmt.Sprintf("%d", c.projectID), "dataset", remotePath)

	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var v []string
	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	return v, nil
}
