package hopsworks

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Engine interface {
	Download(ctx context.Context, remotePath, localPath string) error
}

type ModelEngine struct {
	Engine

	client *Client
}

func NewModelEngine(base Engine, client *Client) *ModelEngine {
	return &ModelEngine{Engine: base, client: client}
}

func (e *ModelEngine) Download(ctx context.Context, m *Model) (string, error) {
	dir, err := os.MkdirTemp("tmp", "*")
	if err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	modelPath := path.Join(dir, m.Name, strconv.Itoa(m.Version))

	if err := os.MkdirAll(modelPath, 0755); err != nil {
		return "", fmt.Errorf("mkdir all: %w", err)
	}

	fromPath := m.VersionPath()
	if strings.HasSuffix(fromPath, "hdsf:/") {
		idx := strings.Index(fromPath, "/Projects")
		if idx == -1 {
			return "", fmt.Errorf("invalid path: %s", fromPath)
		}
		fromPath = fromPath[idx:]
	}

	if err := e.downloadFromHWFS(ctx, fromPath, modelPath); err != nil {
		return "", err
	}

	return modelPath, nil
}

func (e *ModelEngine) downloadFromHWFS(ctx context.Context, fromPath, localPath string) error {
	_, _, err := e.downloadRecursive(ctx, fromPath, localPath)
	if err != nil {
		return err
	}

	return nil
}

func (e *ModelEngine) downloadRecursive(ctx context.Context, fromPath, localPath string) (int, int, error) {
	var (
		nDirs  int
		nFiles int
		err    error
	)

	items, err := e.client.ListDataset(ctx, fromPath)
	if err != nil {
		return 0, 0, fmt.Errorf("list dataset: %w", err)
	}

	for _, item := range items {
		pathAttr := item.Attributes
		path := item.Path
		basename := filepath.Base(path)

		if pathAttr["dir"].(bool) {
			if basename == "Artifacts" {
				continue // skip Artifacts subfolder
			}
			localFolderPath := filepath.Join(localPath, basename)
			if err := os.Mkdir(localFolderPath, os.ModePerm); err != nil {
				return 0, 0, fmt.Errorf("mkdir: %w", err)
			}
			nDirs, nFiles, err = e.downloadRecursive(ctx, path, localFolderPath)
			if err != nil {
				return 0, 0, err
			}
			nDirs++
		} else {
			localFilePath := filepath.Join(localPath, basename)
			e.Engine.Download(ctx, path, localFilePath)
			nFiles++
		}
	}

	return nDirs, nFiles, nil
}
