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
	Version                   int
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

	engine *ModelEngine
}

// VersionPath returns the path to the model version in the Hopsworks file system.
func (m *Model) VersionPath() string {
	return fmt.Sprintf("%s/%d", m.ModelPath(), m.Version)
}

// ModelPath returns the path to the model in the Hopsworks file system.
func (m *Model) ModelPath() string {
	return fmt.Sprintf("/Projects/%s/Models/%s", m.ProjectName, m.Name)
}

// Download downloads the model files and return absolute path to local folder containing them.
func (m *Model) Download(ctx context.Context) (string, error) {
	return m.engine.Download(ctx, m)
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
		engine: NewModelEngine(&LocalEngine{client: c}, c),
	}

	return m, nil
}
