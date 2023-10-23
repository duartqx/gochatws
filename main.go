package main

import (
	"log"

	a "github.com/duartqx/gochatws/domains/auth/auth"
	u "github.com/duartqx/gochatws/domains/auth/users"
	c "github.com/duartqx/gochatws/domains/chat"

	"github.com/go-playground/validator/v10"
	// "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setApp(db *sqlx.DB) *fiber.App {

	app := fiber.New()

	secret := []byte("secret")
	v := validator.New()

	userRepository := u.NewUserRepository(db, v)
	chatRoomRepository := c.NewChatRoomRepository(db, userRepository)

	userController := u.NewUserController(userRepository)
	authController := a.NewJwtAuthController(userRepository, &secret)
	chatRoomController := c.NewChatRoomController(chatRoomRepository)

	// Auth endpoints
	app.
		Post("/register", userController.Create).
		Post("/login", authController.Login).
		Delete("/logout", authController.AuthMiddleware, authController.Logout)

	// Users endpoints
	app.Group("/users").
		// Middleware
		Use(authController.AuthMiddleware).
		// Endpoints
		Get("/", userController.All).
		Get("/:id<int>", userController.Get).
		Put("/:id<int>", userController.Update).
		Delete("/:id<int>", userController.Delete)

	app.Group("/chat").
		Get("/", chatRoomController.All).
		Post("/", chatRoomController.Create).
		Get("/:id<int>", chatRoomController.One)

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
