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
	seq        = 2 //Сделал так, потому что в оригинальной таблице у меня уже есть пользователь с id = 1
	conn, err  = dbr.Open("mysql", "root:@tcp(localhost:3306)/test", nil)
	sess       = conn.NewSession(nil)
)

func InsertUser(c echo.Context) error {
	user := &User{
		Id: seq,
	}
	if err := c.Bind(user); err != nil {
		return err
	}
	sess.InsertInto(userstable).Columns("name").Values(user.Name).Exec()
	seq++
	return c.JSON(http.StatusCreated, user)
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

	e.POST("/user", InsertUser)
	e.GET("/user/:id", SelectUser)

	e.Logger.Fatal(e.Start(":8080"))
}
