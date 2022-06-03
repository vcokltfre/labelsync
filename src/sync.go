package src

import (
	"context"
	"fmt"

	"github.com/google/go-github/v45/github"
)

func syncRepository(repo Repository, client *github.Client) error {
	labels, _, err := client.Issues.ListLabels(context.Background(), repo.Owner, repo.Name, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return err
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

	if len(toDelete) == 0 && len(toCreate) == 0 && len(toUpdate) == 0 {
		fmt.Printf("No changes required for %s/%s\n", repo.Owner, repo.Name)
		return nil
	}

	for _, label := range toCreate {
		fmt.Println(fmt.Sprintf("[%s/%s] Creating label %s", repo.Owner, repo.Name, label.Name))

		_, _, err := client.Issues.CreateLabel(context.Background(), repo.Owner, repo.Name, &github.Label{
			Name:        &label.Name,
			Description: &label.Description,
			Color:       &label.Color,
		})
		if err != nil {
			return err
		}
	}

	for _, name := range toDelete {
		fmt.Println(fmt.Sprintf("[%s/%s] Deleting label %s", repo.Owner, repo.Name, name))

		_, err := client.Issues.DeleteLabel(context.Background(), repo.Owner, repo.Name, name)
		if err != nil {
			return err
		}
	}

	for name, label := range toUpdate {
		fmt.Println(fmt.Sprintf("[%s/%s] Updating label %s", repo.Owner, repo.Name, label.Name))

		_, _, err := client.Issues.EditLabel(context.Background(), repo.Owner, repo.Name, name, &github.Label{
			Description: &label.Description,
			Color:       &label.Color,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func Sync(schema *Schema, client *github.Client) error {
	for _, repo := range schema.Repositories {
		if err := syncRepository(repo, client); err != nil {
			return err
		}
	}

	return nil
}
