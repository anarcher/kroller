package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// create a Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	build(ctx, client)
	test(ctx, client)
}

func build(ctx context.Context, client *dagger.Client) {
	c := client.Container().
		From("golang:1.18").
		WithMountedDirectory("/src", client.Host().Directory(".")).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-o", "app"})

	stdout, err := c.Stdout(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("build: OK %s\n", stdout)
}

func test(ctx context.Context, client *dagger.Client) {
	c := client.Container().
		From("golang:1.18").
		WithMountedDirectory("/src", client.Host().Directory(".")).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "test", "-v", "./..."})

	stdout, err := c.Stdout(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("test: %s\n", stdout)
}
