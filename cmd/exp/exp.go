package main

import (
	"fmt"

	"github.com/welschma/gowebdev/models"
)

type Order struct {
	ID          int
	UserID      int
	Amount      int
	Description string
}

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	us := models.UserService{DB: db}

	user, err := us.Create("bob@bob.com", "bob123")

	if err != nil {
		panic(err)
	}

	fmt.Println(user)

}
