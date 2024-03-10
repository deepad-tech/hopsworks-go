# Hopsworks Go Client

hopsworks-go is the SDK for interacting with a Hopsworks cluster.

_WARNING: Package is still under development. Do not use in production._

## Getting Started

Instantiate a connection and get the project object

```go
client := hopsworks.NewClient(os.Getenv("HOPSWORKS_API_KEY"))
project, _ := client.Login(ctx)
```

Get a model registry and model metadata

```go
mr, _ := project.GetModelRegistry(ctx)
model, _ := mr.GetModel(ctx, "my-model", 1)
```

Download the model for further use

```go
model.Download(ctx)
```

## Contributing

We don't have any guidelines yet. But contributors are welcome!
