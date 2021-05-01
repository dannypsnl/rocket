package main

import (
	"fmt"
	"log"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/response"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db, err = gorm.Open("sqlite3", "test.db")
)

func init() {
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	gorm.Model
	Name string `route:"name"`
	Age  uint64 `route:"age"`
}

func (u *User) GetAge() (age uint64) {
	db.First(u, "name = ?", u.Name)
	return u.Age
}

func setUserAge(u *User) *response.Response {
	db.AutoMigrate(&User{})
	db.Create(u)
	return response.Redirect(fmt.Sprintf("/user/%s", u.Name))
}

func getUser(u *User) string {
	return fmt.Sprintf("User %s age is %d", u.Name, u.GetAge())
}

func main() {
	defer db.Close()

	rocket.Ignite(8080).
		Mount(
			rocket.Get("/user/:name/:age", setUserAge),
			rocket.Get("/user/:name", getUser),
		).
		Launch()
}
