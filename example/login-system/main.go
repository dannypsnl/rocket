package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dannypsnl/rocket"
	"github.com/dannypsnl/rocket/cookie"
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
	Name     string `form:"username"`
	Password string `form:"password"`
}

func login(u *User) *response.Response {
	lookupU := &User{
		Name:     u.Name,
		Password: "",
	}
	db.First(lookupU, "name = ?", lookupU.Name)
	if u.Password != lookupU.Password {
		return response.Redirect("/")
	}
	return response.Redirect("/").
		Cookies(cookie.New("logged_user", u.Name))
}

type UserStatus struct {
	Logged *http.Cookie `cookie:"logged_user"`
}

func home(s *UserStatus) response.Html {
	if s.Logged != nil {
		return response.Html(fmt.Sprintf(`
<html>
<body>
	<h1>%s logged</h1>
    <form action="/logout" method="POST" enctype="form-data">
        <input type="submit" value="logout" name="submit" class="btn btn-success">
	</form>
</body>
</html>
`, s.Logged.Value))
	}
	return `
<html>
<body>
	<form action="/login" method="POST" enctype="form-data">
        <input type="text" name="username" id="username"><br>
        <input type="text" name="password" id="password"><br>
        <input type="submit" value="login" name="submit" class="btn btn-success">
	</form>
</body>
</html>
`
}

func logout() *response.Response {
	return response.Redirect("/").
		Cookies(cookie.New("logged_user", "").MaxAge(-1))
}

func main() {
	db.AutoMigrate(&User{})
	db.Create(&User{
		Model:    gorm.Model{},
		Name:     "aaa",
		Password: "aaa",
	})
	db.Create(&User{
		Model:    gorm.Model{},
		Name:     "bbb",
		Password: "bbb",
	})

	rocket.Ignite(8080).
		Mount(
			rocket.Get("/", home),
			rocket.Post("/login", login),
			rocket.Post("/logout", logout),
		).
		OnClose(func() error {
			return db.Close()
		}).
		Launch()
}
