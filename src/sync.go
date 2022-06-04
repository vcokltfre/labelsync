package src

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/v45/github"
)

func syncRepository(repo Repository, client *github.Client) {
	labels, _, err := client.Issues.ListLabels(context.Background(), repo.Owner, repo.Name, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		panic(err)
	}

	required := map[string]Label{}
	existing := map[string]github.Label{}

	toCreate := []Label{}
	toDelete := []string{}
	toUpdate := map[string]Label{}

	for _, label := range repo.Labels {
		required[label.Name] = label
	}

	for _, label := range labels {
		existing[*label.Name] = *label
	}

	for name, label := range required {
		lab, ok := existing[name]; if !ok {
			toCreate = append(toCreate, label)
		} else if *lab.Color != label.Color || *lab.Description != label.Description {
			toUpdate[name] = label
		}
	}

	for name := range existing {
		_, ok := required[name]; if !ok {
			toDelete = append(toDelete, name)
		}
	}

	todo := len(toDelete) + len(toCreate) + len(toUpdate)

	if todo == 0 {
		fmt.Printf("No changes required for %s/%s\n", repo.Owner, repo.Name)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(todo)

	for _, label := range toCreate {
		go func(lbl Label) {
			fmt.Println(fmt.Sprintf("[%s/%s] Creating label %s", repo.Owner, repo.Name, lbl.Name))

			_, _, err := client.Issues.CreateLabel(context.Background(), repo.Owner, repo.Name, &github.Label{
				Name:        &lbl.Name,
				Description: &lbl.Description,
				Color:       &lbl.Color,
			})
			if err != nil {
				panic(err)
			}

			wg.Done()
		}(label)
	}

	for _, name := range toDelete {
		go func(labelName string) {
			fmt.Println(fmt.Sprintf("[%s/%s] Deleting label %s", repo.Owner, repo.Name, labelName))

		_, err := client.Issues.DeleteLabel(context.Background(), repo.Owner, repo.Name, labelName)
		if err != nil {
			panic(err)
		}

		wg.Done()
		}(name)
	}

	for name, label := range toUpdate {
		go func(labelName string, lbl Label) {
			fmt.Println(fmt.Sprintf("[%s/%s] Updating label %s", repo.Owner, repo.Name, lbl.Name))

		_, _, err := client.Issues.EditLabel(context.Background(), repo.Owner, repo.Name, labelName, &github.Label{
			Description: &lbl.Description,
			Color:       &lbl.Color,
		})
		if err != nil {
			panic(err)
		}

		wg.Done()
		}(name, label)
	}

	wg.Wait()
}

func Sync(schema *Schema, client *github.Client) {
	wg := sync.WaitGroup{}
	wg.Add(len(schema.Repositories))

	for _, repo := range schema.Repositories {
		go func(repo Repository) {
			syncRepository(repo, client)

			wg.Done()
		}(repo)
	}

	wg.Wait()
}
