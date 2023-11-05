package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"

	"github.com/duartqx/gochatws/infrastructure/sessions"

	c "github.com/duartqx/gochatws/api/controllers"
	ac "github.com/duartqx/gochatws/api/controllers/auth"

	s "github.com/duartqx/gochatws/application/services"

	"github.com/duartqx/gochatws/api/utils"
	w "github.com/duartqx/gochatws/api/ws"
	r "github.com/duartqx/gochatws/infrastructure/repositories/sqlite"
)

func getTemplateEngine() *html.Engine {
	templEngine := html.New("./presentation/views", ".html")
	templEngine.AddFuncMap(map[string]interface{}{
		"GetChatCategories": utils.GetChatCategories,
	})
	return templEngine
}

func Setup(db *sqlx.DB, secret []byte) *fiber.App {

	app := fiber.New(
		fiber.Config{Views: getTemplateEngine(), ViewsLayout: "base"},
	)

	// Raw dependencies
	v := validator.New()
	sessionStore := sessions.NewSessionStore()
	connStore := w.GetConnectionStore()

	// Repositories
	userRepository := r.NewUserRepository(db)
	chatRoomRepository := r.NewChatRoomRepository(db, userRepository)
	messageRepository := r.NewMessageRepository(
		db, userRepository, chatRoomRepository,
	)

	// Services
	jwtAuthService := s.NewJwtAuthService(userRepository, &secret, sessionStore)
	userService := s.NewUserService(userRepository, v)
	chatRoomService := s.NewChatRoomService(chatRoomRepository)
	messageService := s.NewMessageService(messageRepository)

	// Controllers
	userController := c.NewUserController(userService)
	authController := ac.NewJwtAuthController(jwtAuthService)
	chatRoomController := c.NewChatRoomController(chatRoomService)
	msgController := c.NewMessageController(chatRoomRepository, messageService, connStore)

	// Logger middleware
	app.Use(
		logger.New(
			logger.Config{TimeFormat: "2006-01-02 15:04:05"},
		),
	)

	// Static files (js, css, images)
	app.Static("/", "./presentation/static")

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
		// Chat id param
		Group("/:chat_id<int>").
		Get("/msg", msgController.GetChatMessages).
		Post("/msg", msgController.CreateMessage).
		Use("/ws", func(c *fiber.Ctx) error {
			// IsWebSocketUpgrade returns true if the client
			// requested upgrade to the WebSocket protocol.
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("allowed", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		}).
		Get("/ws/connect", msgController.WebSocketChat())

	// HTML Template endpoints
	app.
		// Authenticated endpoints
		Get(
			"/",
			authController.AuthMiddlewareWithRedirect(),
			func(c *fiber.Ctx) error {
				return c.Render(
					"index",
					utils.BuildTemplateContext(c, &fiber.Map{"Title": "Index"}))
			},
		).
		Get(
			"/chat/:id<int>",
			authController.AuthMiddlewareWithRedirect(),
			chatRoomController.ChatView,
		).
		// Unauthenticated endpoints
		Get(
			"/login",
			authController.AuthNotLoggedMiddlewareWithRedirect(),
			func(c *fiber.Ctx) error {
				return c.Render("login", fiber.Map{"Title": "Login"})
			}).
		Get(
			"/register",
			authController.AuthNotLoggedMiddleware(),
			func(c *fiber.Ctx) error {
				return c.Render("register", fiber.Map{"Title": "Register"})
			})

	return app
}
