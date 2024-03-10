package hopsworks

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type ModelRegistry struct {
	ProjectID         uint64
	ProjectName       string
	ID                uint64
	SharedProjectName string

	client *Client
}

type GetModelRegistryResponse struct {
	Items []struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}

func (c *Client) GetModelRegistry(ctx context.Context, project string) (*ModelRegistry, error) {
	url := c.url("project", strconv.Itoa(int(c.projectID)), "modelregistries")

	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var v GetModelRegistryResponse
	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	var (
		registryID = c.projectID
		sharedName string
	)
	for _, item := range v.Items {
		if item.Name == project {
			registryID = item.ID
			sharedName = project
		}
	}

	if sharedName == "" {
		return nil, fmt.Errorf(
			"no model registry shared with current project %s, from project %s",
			c.config.Project, project)
	}

	r := &ModelRegistry{
		ProjectID:         c.projectID,
		ProjectName:       c.config.Project,
		ID:                registryID,
		SharedProjectName: sharedName,
		client:            c,
	}

	return r, nil
}

type GetModelResponse struct {
	// TODO
}

func (r *ModelRegistry) GetModel(ctx context.Context, name string, version int, registryID uint64) (*Model, error) {
	url := r.client.url(
		"project",
		fmt.Sprintf("%d", r.client.projectID),
		"modelregistries",
		fmt.Sprintf("%d", registryID),
		"models",
		fmt.Sprintf("%s_%d", name, version),
	)
	queryArgs := map[string]string{
		"expand": "trainingdatasets",
	}

	req, err := r.client.newRequest(ctx, http.MethodGet, url, withQueryArgs(queryArgs))
	if err != nil {
		return nil, err
	}

	var v GetModelResponse
	if err := r.client.sendRequest(req, &v); err != nil {
		return nil, err
	}

	m := &Model{}

	return m, nil
}
