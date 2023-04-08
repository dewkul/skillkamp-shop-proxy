package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(listenAddr string, upstream string) *Server {
	return &Server{
		listenAddr: listenAddr,
		serverUrl:  upstream,
	}
}

func (s *Server) Start() error {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Get("/v1/api/filters", s.handleGetFilters)
	app.Post("/v1/api/auth/login", s.handlePostLogin)
	app.Post("/v1/api/auth/signup", s.handlePostSignup)
	return app.Listen(s.listenAddr)
}

func (s *Server) handleGetFilters(c *fiber.Ctx) error {
	return s.handleGetProxy(c, "/v1/api/filters")
}

func (s *Server) handlePostLogin(c *fiber.Ctx) error {
	return s.handlePostProxy(c, "/v1/api/auth/login")
}

func (s *Server) handlePostSignup(c *fiber.Ctx) error {
	return s.handlePostProxy(c, "/v1/api/auth/signup")
}

func (s *Server) handleGetProxy(c *fiber.Ctx, path string) error {
	resp, err := http.Get(s.serverUrl + path)
	if err != nil {
		return fiber.ErrBadGateway
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	return c.Send(body)
}

func (s *Server) handlePostProxy(c *fiber.Ctx, path string) error {
	reqBody := bytes.NewBuffer(c.Body())
	resp, err := http.Post(s.serverUrl+path, "application/json", reqBody)
	if err != nil {
		return fiber.ErrBadGateway
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	return c.Status(resp.StatusCode).Send(body)
}

type Server struct {
	listenAddr string
	serverUrl  string
}
