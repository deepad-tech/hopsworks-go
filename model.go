package hopsworks

import "time"

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
}
