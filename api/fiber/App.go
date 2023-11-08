package api

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"

	"github.com/duartqx/gochatws/infrastructure/sessions"

	c "github.com/duartqx/gochatws/api/fiber/controllers"
	ac "github.com/duartqx/gochatws/api/fiber/controllers/auth"

	s "github.com/duartqx/gochatws/application/services"

	"github.com/duartqx/gochatws/api/fiber/utils"
	w "github.com/duartqx/gochatws/api/fiber/ws"
	r "github.com/duartqx/gochatws/infrastructure/repositories/sqlite"
)

type App struct {
	app          *fiber.App
	db           *sqlx.DB
	secret       *[]byte
	port         string
	viewsBase    string
	viewsPath    string
	staticPath   string
	sessionStore *sessions.SessionStore
	v            *validator.Validate
}

func GetNewAppBuilder() *App {
	return &App{
		sessionStore: sessions.NewSessionStore(),
		v:            &validator.Validate{},
	}
}

func (a *App) SetDb(db *sqlx.DB) *App {
	a.db = db
	return a
}

func (a *App) SetPort(port string) *App {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	a.port = port
	return a
}

func (a *App) SetViewsPath(viewsPath string) *App {
	a.viewsPath = viewsPath
	return a
}

func (a *App) SetViewsBase(viewsBase string) *App {
	a.viewsBase = viewsBase
	return a
}

func (a *App) SetStaticPath(staticPath string) *App {
	a.staticPath = staticPath
	return a
}

func (a *App) SetSecret(secret string) *App {
	s := []byte(secret)
	a.secret = &s
	return a
}

func (a *App) SetValidator(v *validator.Validate) *App {
	a.v = v
	return a
}

func (a *App) SetSessionStore(sessionStore *sessions.SessionStore) *App {
	a.sessionStore = sessionStore
	return a
}

func (a App) getTemplateEngine() *html.Engine {
	templEngine := html.New(a.viewsPath, ".html")
	templEngine.AddFuncMap(map[string]interface{}{
		"GetChatCategories": utils.GetChatCategories,
	})
	return templEngine
}

func (a *App) Build() *App {

	a.app = fiber.New(
		fiber.Config{Views: a.getTemplateEngine(), ViewsLayout: a.viewsBase},
	)

	// Repositories
	userRepository := r.NewUserRepository(a.db)
	chatRoomRepository := r.NewChatRoomRepository(a.db, userRepository)
	messageRepository := r.NewMessageRepository(
		a.db, userRepository, chatRoomRepository,
	)

	// Services
	jwtAuthService := s.NewJwtAuthService(userRepository, a.secret, a.sessionStore)
	userService := s.NewUserService(userRepository, a.v)
	chatRoomService := s.NewChatRoomService(chatRoomRepository)
	messageService := s.NewMessageService(messageRepository, chatRoomRepository)
	webSocketService := w.GetWebSocketService(messageRepository)

	// Controllers
	userController := c.NewUserController(userService)
	authController := ac.NewJwtAuthController(jwtAuthService)
	chatRoomController := c.NewChatRoomController(chatRoomService)
	msgController := c.NewMessageController(messageService, webSocketService)

	// Logger middleware
	a.app.Use(
		logger.New(
			logger.Config{TimeFormat: "2006-01-02 15:04:05"},
		),
	)

	// Static files (js, css, images)
	a.app.Static("/", a.staticPath)

	// Groups with prefix /api
	apiGroup := a.app.Group("/api")

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
		Get("/ws/connect", msgController.WebSocketChatController())

	// HTML Template endpoints
	a.app.
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

	return a
}

func (a *App) Listen() error {
	if err := a.app.Listen(a.port); err != nil {
		return err
	}
	return nil
}

func (a *App) Shutdown() error {
	return a.app.Shutdown()
}
