package stores

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	databaseURL string
	db          *sqlx.DB
}

func NewStore(databaseURL string) (Store, error) {
	fmt.Println(databaseURL)
	database, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return Store{}, err
	}

	return Store{
		db:          database,
		databaseURL: databaseURL,
	}, nil
}

func (s *Store) GetDB() *sqlx.DB {
	return s.db
}
