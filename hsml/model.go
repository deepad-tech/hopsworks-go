package hsml

import (
	"context"
	"fmt"
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
}

func (m *Model) VersionPath() string {
	return fmt.Sprintf("%s/%d", m.ModelPath(), m.Version)
}

func (m *Model) ModelPath() string {
	return fmt.Sprintf("/Projects/%s/Models/%s", m.ProjectName, m.Name)
}

type GetModelResponse struct {
	// TODO
}

// Download downloads the model files and return absolute path to local folder containing them.
func (m *Model) Download(ctx context.Context) (string, error) {
	// TODO
	return "", nil
}
