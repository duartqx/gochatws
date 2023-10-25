package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	// "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	as "github.com/duartqx/gochatws/core/auth/service"
	ac "github.com/duartqx/gochatws/domains/controllers/auth"

	"github.com/duartqx/gochatws/core/sessions"

	c "github.com/duartqx/gochatws/domains/controllers"
	r "github.com/duartqx/gochatws/domains/repositories"
	s "github.com/duartqx/gochatws/domains/services"
)

func setApp(db *sqlx.DB) *fiber.App {

	app := fiber.New()

	// Raw dependencies
	secret := []byte("secret")
	v := validator.New()
	sessionStore := sessions.NewSessionStore()

	// Repositories
	userRepository := r.NewUserRepository(db)
	chatRoomRepository := r.NewChatRoomRepository(db, userRepository)

	// Services
	jwtAuthService := as.NewJwtAuthService(userRepository, &secret, sessionStore)
	userService := s.NewUserService(userRepository, v)
	chatRoomService := s.NewChatRoomService(chatRoomRepository)

	// Controllers
	userController := c.NewUserController(userService)
	authController := ac.NewJwtAuthController(jwtAuthService)
	chatRoomController := c.NewChatRoomController(chatRoomService)

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
		// Middleware
		Use(authController.AuthMiddleware).
		// Endpoints
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
