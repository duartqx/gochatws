package main

import (
	"log"

	u "gochatws/domains/users"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func SetUserRoutes(r *fiber.Router, db *sqlx.DB, v *validator.Validate) {
	userController := u.NewUserController(u.NewUserRepo(db, v))

	(*r).Get("/", userController.All)
	(*r).Post("/", userController.Create)
	(*r).Get("/:id<int>", userController.Get)
	(*r).Put("/:id<int>", userController.Update)
}

func main() {

	db, err := sqlx.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	v := validator.New()

	app := fiber.New()

	UserRouter := app.Group("/users")
	SetUserRoutes(&UserRouter, db, v)

	log.Fatal(app.Listen(":8000"))
}
