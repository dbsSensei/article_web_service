package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	db "github.com/dbssensei/articlewebservice/db/sqlc"
	"os"
)

type Tag struct {
	Id   string
	Name string
}

func (s Seed) Tag(store db.Store) {
	jsonFile, err := os.ReadFile("data/tags.json")
	if err != nil {
		fmt.Println("error when parse json file", err)
	}
	var tags []Tag
	json.Unmarshal(jsonFile, &tags)

	totalTag, err := store.CountTag(context.Background())
	if err != nil {
		fmt.Println("error when get totalTag", err)
	}

	if totalTag == 0 {
		for _, tag := range tags {
			store.CreateTag(context.Background(), tag.Name)
		}
		totalTag, err = store.CountTag(context.Background())
		if err != nil {
			fmt.Println("error when get totalTag", err)
		} else if int(totalTag) < len(tags) {
			fmt.Println("Tag seeding incomplete")
		} else {
			fmt.Println("Tag seeding successful")
		}
	}
}
