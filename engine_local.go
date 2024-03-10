package hopsworks

import "context"

type LocalEngine struct {
	client *Client
}

var _ Engine = (*LocalEngine)(nil)

func (e *LocalEngine) Download(ctx context.Context, remotePath, localPath string) error {
	return e.client.DownloadDatasetFile(ctx, remotePath, localPath)
}
