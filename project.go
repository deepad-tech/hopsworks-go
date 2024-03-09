package hopsworks

import (
	"context"
	"net/http"
)

type Project struct{}

type GetProjectResponse struct {
	httpHeader
}

// Get project returns project for the specified name.
func (c *Client) GetProject(ctx context.Context, name string) (*Project, error) {
	url := c.url("project", "projectInfo", name)

	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var v GetProjectResponse

	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	p := &Project{}

	return p, nil
}

type GetProjectsResponse struct {
	httpHeader
}

// GetProjects returns all projects accessible by the user.
func (c *Client) GetProjects(ctx context.Context) ([]*Project, error) {
	url := c.url("project")

	req, err := c.newRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	var v GetProjectsResponse

	if err := c.sendRequest(req, &v); err != nil {
		return nil, err
	}

	projects := make([]*Project, 0)

	return projects, nil
}
