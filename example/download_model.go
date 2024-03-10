package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deepad-tech/hopsworks-go"
)

func main() {
	ctx := context.Background()

	client := hopsworks.NewClient(os.Getenv("HOPSWORKS_API_KEY"), os.Getenv("HOPSWORKS_PROJECT"))
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

	path, err := model.Download(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(path)
}
