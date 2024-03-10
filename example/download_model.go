package main

import (
	"context"
	"log"
	"os"

	"github.com/deepad-tech/hopsworks-go"
)

func main() {
	ctx := context.Background()

	client := hopsworks.NewClient(os.Getenv("HOPSWORKS_API_KEY"))
	project, err := client.Login(ctx)
	if err != nil {
		log.Fatal(err)
	}

	mr, err := project.GetModelRegistry(ctx)
	if err != nil {
		log.Fatal(err)
	}

	model, err := mr.GetModel(ctx, "my-model", 1)
	if err != nil {
		log.Fatal(err)
	}

	if err := model.Download(ctx); err != nil {
		log.Fatal(err)
	}
}
