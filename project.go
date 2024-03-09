package hopsworks

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ProjectClient struct {
	*Client

	ID          uint64
	Name        string
	Owner       string
	Description string
	Created     time.Time
}

func (c *ProjectClient) GetModelRegistry(ctx context.Context) (*ModelRegistry, error) {
	return c.Client.GetModelRegistry(ctx, c.Name)
}

type GetProjectResponse struct {
	ID          uint64    `json:"projectId"`
	Name        string    `json:"projectName"`
	Owner       string    `json:"owner"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

// Get project returns project for the specified name.
func (c *Client) GetProject(ctx context.Context, name string) (*ProjectClient, error) {
	url := c.url("project", "getProjectInfo", name)

	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var v GetProjectResponse
	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	p := &ProjectClient{
		Client:      c,
		ID:          v.ID,
		Name:        v.Name,
		Owner:       v.Owner,
		Description: v.Description,
		Created:     v.Created,
	}

	return p, nil
}

type GetProjectsResponse []struct {
	Project struct {
		Name string `json:"name"`
	} `json:"project"`
}

// GetProjects returns all projects accessible by the user.
func (c *Client) GetProjects(ctx context.Context) ([]*ProjectClient, error) {
	url := c.url("project")

	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var v GetProjectsResponse
	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	projects := make([]*ProjectClient, 0, len(v))

	for _, vv := range v {
		p, err := c.GetProject(ctx, vv.Project.Name)
		if err != nil {
			return nil, fmt.Errorf("get project %s:%w", vv.Project.Name, err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}
