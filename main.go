package main

import (
	a "gochatws/domains/auth/auth"
	u "gochatws/domains/auth/users"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setApp(db *sqlx.DB) *fiber.App {

	app := fiber.New()

	st := session.New()

	v := validator.New()

	userRepository := u.NewUserRepo(db, v)

	userController := u.NewUserController(userRepository)
	authController := a.NewBasicAuthController(userRepository, st)

	// Unauthenticated endpoints
	app.
		Post("/login", authController.Login).
		Post("/register", userController.Create)

	// Authenticated endpoints
	app.Group("/users").
		// Middleware
		Use(authController.AuthenticationMiddleware).
		// Endpoints
		Get("/", userController.All).
		Get("/:id<int>", userController.Get).
		Put("/:id<int>", userController.Update).
		Delete("/:id<int>", userController.Delete)

	return app
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
