package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/duartqx/gochatws/core/sessions"

	as "github.com/duartqx/gochatws/core/auth/service"
	c "github.com/duartqx/gochatws/domains/controllers"
	ac "github.com/duartqx/gochatws/domains/controllers/auth"
	r "github.com/duartqx/gochatws/domains/repositories"
	s "github.com/duartqx/gochatws/domains/services"
)

func setApp(db *sqlx.DB) *fiber.App {

	app := fiber.New(
		fiber.Config{
			Views:       html.New("./domains/views", ".html"),
			ViewsLayout: "base",
		},
	)

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
	wsController := c.NewWebSocketController(chatRoomRepository)

	// Logger middleware
	app.Use(
		logger.New(
			logger.Config{TimeFormat: "2006-01-02 15:04:05"},
		),
	)

	// Static files (js, css, images)
	app.Static("/", "./domains/static")

	// Groups with prefix /api
	apiGroup := app.Group("/api")

	// Auth endpoints
	apiGroup.
		Post(
			"/register",
			authController.AuthNotLoggedMiddleware(),
			userController.Create).
		Post(
			"/login",
			authController.AuthNotLoggedMiddleware(),
			authController.Login).
		Delete(
			"/logout",
			authController.AuthMiddleware(),
			authController.Logout,
		)

	// Users endpoints
	apiGroup.Group("/user").
		// Middleware
		Use(authController.AuthMiddleware()).
		// Endpoints
		Get("/list", userController.All).
		Get("/", userController.Get).
		Put("/", userController.Update).
		Delete("/", userController.Delete)

	// Chat endpoints
	apiGroup.Group("/chat").
		// Middleware
		Use(authController.AuthMiddleware()).
		// Endpoints
		Get("/", chatRoomController.All).
		Post("/", chatRoomController.Create).
		Use("/ws", func(c *fiber.Ctx) error {
			// IsWebSocketUpgrade returns true if the client
			// requested upgrade to the WebSocket protocol.
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("allowed", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		}).
		Get("/ws/:id<int>", wsController.Chat())

	// HTML Template endpoints
	app.
		Get(
			"/",
			authController.AuthMiddlewareWithRedirect(),
			func(c *fiber.Ctx) error {
				return c.Render("index", fiber.Map{"Title": "Index"})
			}).
		Get(
			"/login",
			authController.AuthNotLoggedMiddlewareWithRedirect(),
			func(c *fiber.Ctx) error {
				return c.Render("login", fiber.Map{"Title": "Login"})
			})

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
