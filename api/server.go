package api

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
)

func NewServer(listenAddr, upstream, version string) *Server {
	if version == "" {
		version = "dev"
	}
	return &Server{
		listenAddr: listenAddr,
		serverUrl:  upstream,
		version:    version,
	}
}

func (s *Server) Start() error {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Get("/v1/api/filters", s.handleGetFilters)
	app.Get("/v1/api/products", s.handleGetProducts)
	app.Get("/v1/api/products/new_arrivals", s.handleGetNewArrival)
	app.Post("/v1/api/auth/login", s.handlePostLogin)
	app.Post("/v1/api/auth/signup", s.handlePostSignup)
	app.Get("/v1/api/cart", s.handleGetItemsInCart)
	app.Post("/v1/api/cart", s.handleAddItemsInCart)
	app.Get("/v1/api/products/details/:sku", s.handleGetProductInfo)
	app.Get("/ver", s.handleGetVersion)

	return app.Listen(s.listenAddr)
}

func (s *Server) handleGetFilters(c *fiber.Ctx) error {
	return s.handleGetProxy(c, "/v1/api/filters")
}

func (s *Server) handleGetProducts(c *fiber.Ctx) error {
	queryStrBytes := c.Request().URI().QueryString()
	path := []string{"/v1/api/products"}
	if len(queryStrBytes) > 0 {
		path = append(path, string(queryStrBytes))
	}
	return s.handleGetProxy(c, strings.Join(path, "?"))
}

func (s *Server) handleGetNewArrival(c *fiber.Ctx) error {
	return s.handleGetProxy(c, "/v1/api/products/new_arrivals")
}

func (s *Server) handlePostLogin(c *fiber.Ctx) error {
	return s.handlePostProxy(c, "/v1/api/auth/login")
}

func (s *Server) handlePostSignup(c *fiber.Ctx) error {
	return s.handlePostProxy(c, "/v1/api/auth/signup")
}

func (s *Server) handleGetItemsInCart(c *fiber.Ctx) error {
	return s.handleGetProxy(c, "/v1/api/cart")
}

func (s *Server) handleAddItemsInCart(c *fiber.Ctx) error {
	return s.handlePostProxy(c, "/v1/api/cart")
}

func (s *Server) handleGetProductInfo(c *fiber.Ctx) error {
	return s.handleGetProxy(c, "/v1/api/products/details/"+c.Params("sku"))
}

func (s *Server) handleGetProxy(c *fiber.Ctx, path string) error {
	log.Debug().Str("path", path).Msg("handle get proxy request")
	resp, err := http.Get(s.serverUrl + path)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Str("upstream", s.serverUrl).Msg("GET: Upstream connection error")
		return fiber.ErrBadGateway
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Str("upstream", s.serverUrl).Msg("GET: Read resp body error")
		return fiber.ErrUnprocessableEntity
	}
	log.Debug().Str("path", path).Int("status", resp.StatusCode).Msg("GET: response")
	return c.Status(resp.StatusCode).Send(body)
}

func (s *Server) handlePostProxy(c *fiber.Ctx, path string) error {
	reqBody := bytes.NewBuffer(c.Body())
	resp, err := http.Post(s.serverUrl+path, "application/json", reqBody)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Str("upstream", s.serverUrl).Msg("POST: Upstream connection error")
		return fiber.ErrBadGateway
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Str("upstream", s.serverUrl).Msg("POST: Read resp body error")
		return fiber.ErrUnprocessableEntity
	}
	log.Debug().Str("path", path).Int("status", resp.StatusCode).Msg("POST: response")
	return c.Status(resp.StatusCode).Send(body)
}

func (s *Server) handleGetVersion(c *fiber.Ctx) error {
	resp := VersionResponse{
		Version: s.version,
	}
	return c.JSON(resp)
}

type VersionResponse struct {
	Version string `json:"version"`
}

type Server struct {
	listenAddr string
	serverUrl  string
	version    string
}
