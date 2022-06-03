package main

import (
	"context"
	"os"

	"github.com/google/go-github/v45/github"
	"github.com/joho/godotenv"
	"github.com/vcokltfre/labelsync/src"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var file string

	if len(os.Args) == 1 {
		file = ".gitlabels.yml"
	} else {
		file = os.Args[1]
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	schema, err := src.LoadSchema(file)
	if err != nil {
		panic(err)
	}

	src.Sync(schema, client)
}
