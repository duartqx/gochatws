package main

import (
	u "gochatws/domains/auth/users"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setApp(db *sqlx.DB) *fiber.App {
	app := fiber.New()
	v := validator.New()

	authRouter := app.Group("/users")
	setAuthRoutes(&authRouter, db, v)

	return app
}

func setAuthRoutes(r *fiber.Router, db *sqlx.DB, v *validator.Validate) {

	userRepository := u.NewUserRepo(db, v)
	userController := u.NewUserController(userRepository)

	(*r).
		Get("/", userController.All).
		Post("/", userController.Create).
		Get("/:id<int>", userController.Get).
		Put("/:id<int>", userController.Update).
		Delete("/:id<int>", userController.Delete)

}

func main() {

	db, err := sqlx.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	app := setApp(db)

	log.Fatalln(app.Listen(":8000"))
}
