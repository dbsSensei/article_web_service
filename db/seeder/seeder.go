package seeder

import db "github.com/dbssensei/articlewebservice/db/sqlc"

type Seed struct {
}

func Execute(store db.Store) {
	var seed Seed
	seed.Category(store)
	seed.Tag(store)
}
