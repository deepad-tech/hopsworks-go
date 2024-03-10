package hopsworks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deepad-tech/hopsworks-go/hsml"
)

type GetModelResponse struct {
	// TODO
}

func (c *Client) GetModel(ctx context.Context, name string, version int, registryID uint64) (*hsml.Model, error) {
	url := c.url(
		"project",
		fmt.Sprintf("%d", c.projectID),
		"modelregistries",
		fmt.Sprintf("%d", registryID),
		"models",
		fmt.Sprintf("%s_%d", name, version),
	)
	queryArgs := map[string]string{
		"expand": "trainingdatasets",
	}

	req, err := c.newRequest(ctx, http.MethodGet, url, withQueryArgs(queryArgs))
	if err != nil {
		return nil, err
	}

	var v GetModelResponse
	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	m := &hsml.Model{}

	return m, nil
}
