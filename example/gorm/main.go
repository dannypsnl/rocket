package main

import (
	"fmt"

	"github.com/dannypsnl/rocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	gorm.Model
	Name string `route:"name"`
	Age  uint64
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})
	db.Create(&User{Name: "Danny", Age: 21})

	rocket.Ignite(":8080").
		Mount("/user", rocket.Get("/:name", func(u *User) string {
			var user User
			db.First(&user, "name = ?", u.Name)
			return fmt.Sprintf("User %s age is %d", user.Name, user.Age)
		})).
		Launch()
}
