package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	db "github.com/dbssensei/articlewebservice/db/sqlc"
	"os"
)

type Category struct {
	Id   string
	Name string
}

func (s Seed) Category(store db.Store) {
	jsonFile, err := os.ReadFile("data/categories.json")
	if err != nil {
		fmt.Println("error when parse json file", err)
	}
	var categories []Category
	json.Unmarshal(jsonFile, &categories)

	totalCategory, err := store.CountCategory(context.Background())
	if err != nil {
		fmt.Println("error when get totalCategory", err)
	}

	if totalCategory == 0 {
		for _, category := range categories {
			store.CreateCategory(context.Background(), category.Name)
		}
		totalCategory, err = store.CountCategory(context.Background())
		if err != nil {
			fmt.Println("error when get totalCategory", err)
		} else if int(totalCategory) < len(categories) {
			fmt.Println("Category seeding incomplete")
		} else {
			fmt.Println("Category seeding successful")
		}
	}
}
