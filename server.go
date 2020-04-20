package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type User struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

var (
	userstable = "users"
	seq        = 1
	conn, err  = dbr.Open("mysql", "root:@tcp(localhost:3306)/test", nil)
	sess       = conn.NewSession(nil)
)

func InsertUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	sess.InsertInto(userstable).Columns("id", "name").Values(user.Id, user.Name).Exec()
	return c.NoContent(http.StatusOK)
}

func SelectUser(c echo.Context) error {
	user := new(User)
	id := c.Param("id")
	sess.Select("*").From(userstable).Where("id=?", id).Load(&user)
	return c.JSON(http.StatusOK, user)
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/authors", InsertUser)
	e.GET("/author/:id", SelectUser)

	e.Logger.Fatal(e.Start(":8080"))
}
