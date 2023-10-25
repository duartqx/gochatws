package main

import (
	"log"

	ac "github.com/duartqx/gochatws/core/auth/controllers"
	as "github.com/duartqx/gochatws/core/auth/service"
	s "github.com/duartqx/gochatws/core/sessions"
	c "github.com/duartqx/gochatws/domains/chat"
	u "github.com/duartqx/gochatws/domains/users"

	"github.com/go-playground/validator/v10"
	// "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setApp(db *sqlx.DB) *fiber.App {

	app := fiber.New()

	// Raw dependencies
	secret := []byte("secret")
	v := validator.New()
	sessionStore := s.NewSessionStore()

	// Repositories
	userRepository := u.NewUserRepository(db, v)
	chatRoomRepository := c.NewChatRoomRepository(db, userRepository)

	// Services
	jwtAuthService := as.NewJwtAuthService(userRepository, &secret, sessionStore)

	// Controllers
	userController := u.NewUserController(userRepository)
	authController := ac.NewJwtAuthController(jwtAuthService)
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
