package main

import (
	"fmt"

	"github.com/dannypsnl/rocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db, err = gorm.Open("sqlite3", "test.db")
)

func init() {
	if err != nil {
		panic(err)
	}
}

type User struct {
	gorm.Model
	Name string `route:"name"`
	Age  uint64
}

func (u *User) GetAge() (age uint64) {
	db.First(u, "name = ?", u.Name)
	return u.Age
}

func main() {
	defer db.Close()

	db.AutoMigrate(&User{})
	db.Create(&User{Name: "Danny", Age: 21})

	rocket.Ignite(":8080").
		Mount("/user", rocket.Get("/:name", func(u *User) string {
			return fmt.Sprintf("User %s age is %d", u.Name, u.GetAge())
		})).
		Launch()
}
