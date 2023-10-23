package main

import (
	"log"

	a "github.com/duartqx/gochatws/domains/auth/auth"
	s "github.com/duartqx/gochatws/domains/auth/sessions"
	u "github.com/duartqx/gochatws/domains/auth/users"
	c "github.com/duartqx/gochatws/domains/chat"

	"github.com/go-playground/validator/v10"
	// "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setApp(db *sqlx.DB) *fiber.App {

	app := fiber.New()

	secret := []byte("secret")
	v := validator.New()
	sessionStore := s.NewSessionStore()

	userRepository := u.NewUserRepository(db, v)
	chatRoomRepository := c.NewChatRoomRepository(db, userRepository)

	userController := u.NewUserController(userRepository)
	authController := a.NewJwtAuthController(
		userRepository, &secret, sessionStore,
	)
	chatRoomController := c.NewChatRoomController(chatRoomRepository)

	// Logger middleware
	app.Use(
		logger.New(
			logger.Config{TimeFormat: "2006-01-02 15:04:05"},
		),
	)

	// Auth endpoints
	app.
		Post(
			"/register",
			authController.AuthNotLoggedMiddleware,
			userController.Create).
		Post(
			"/login",
			authController.AuthNotLoggedMiddleware,
			authController.Login).
		Delete("/logout", authController.AuthMiddleware, authController.Logout)

	// Users endpoints
	app.Group("/user").
		// Middleware
		Use(authController.AuthMiddleware).
		// Endpoints
		Get("/list", userController.All).
		Get("/", userController.Get).
		Put("/", userController.Update).
		Delete("/", userController.Delete)

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
