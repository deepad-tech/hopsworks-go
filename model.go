package hopsworks

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Model struct {
	ID                        int
	Name                      string
	Version                   string
	Description               string
	Created                   time.Time
	Environment               string
	ExperimentID              string
	ProjectName               string
	ExperimentProjectName     string
	TrainingMetrics           interface{} // Change to specific data structure
	Program                   string
	UserFullName              string
	InputExample              string
	Framework                 string
	ModelSchema               string
	TrainingDataset           string
	SharedRegistryProjectName string
	ModelRegistryID           string

	client *Client
}

type GetModelResponse struct {
	// TODO
}

func (c *Client) GetModel(ctx context.Context, name string, version int, registryID uint64) (*Model, error) {
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

	m := &Model{
		client: c,
	}

	return m, nil
}

func (m *Model) Download(ctx context.Context) error {
	// TODO
	return nil
}
