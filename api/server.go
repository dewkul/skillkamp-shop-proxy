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

func NewServer(listenAddr, upstream, version, origins string) *Server {
	if version == "" {
		version = "dev"
	}
	return &Server{
		listenAddr:   listenAddr,
		serverUrl:    upstream,
		version:      version,
		allowOrigins: origins,
	}
}

func (s *Server) Start() error {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: s.allowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Get("/v1/api/filters", s.handleGetFilters)
	app.Get("/v1/api/products", s.handleGetProducts)
	app.Get("/v1/api/products/new_arrivals", s.handleGetNewArrival)
	app.Post("/v1/api/auth/login", s.handlePostLogin)
	app.Post("/v1/api/auth/signup", s.handlePostSignup)
	app.Get("/v1/api/cart", s.handleGetItemsInCart)
	app.Post("/v1/api/cart", s.handleAddItemsInCart)
	app.Put("/v1/api/cart", s.handleUpdateItemsInCart)
	app.Delete("/v1/api/cart", s.handleRemoveItemsInCart)
	app.Get("/v1/api/products/details/:sku", s.handleGetProductInfo)
	app.Get("/ver", s.handleGetVersion)

	return app.Listen(s.listenAddr)
}

func (s *Server) handleGetFilters(c *fiber.Ctx) error {
	return s.handleProxy(c, "GET", "/v1/api/filters")
}

func (s *Server) handleGetProducts(c *fiber.Ctx) error {
	queryStrBytes := c.Request().URI().QueryString()
	path := []string{"/v1/api/products"}
	if len(queryStrBytes) > 0 {
		path = append(path, string(queryStrBytes))
	}
	return s.handleProxy(c, "GET", strings.Join(path, "?"))
}

func (s *Server) handleGetNewArrival(c *fiber.Ctx) error {
	return s.handleProxy(c, "GET", "/v1/api/products/new_arrivals")
}

func (s *Server) handlePostLogin(c *fiber.Ctx) error {
	return s.handleProxy(c, "POST", "/v1/api/auth/login")
}

func (s *Server) handlePostSignup(c *fiber.Ctx) error {
	return s.handleProxy(c, "POST", "/v1/api/auth/signup")
}

func (s *Server) handleGetItemsInCart(c *fiber.Ctx) error {
	return s.handleProxy(c, "GET", "/v1/api/cart")
}

func (s *Server) handleAddItemsInCart(c *fiber.Ctx) error {
	return s.handleProxy(c, "POST", "/v1/api/cart")
}

func (s *Server) handleUpdateItemsInCart(c *fiber.Ctx) error {
	return s.handleProxy(c, "PUT", "/v1/api/cart")
}

func (s *Server) handleRemoveItemsInCart(c *fiber.Ctx) error {
	return s.handleProxy(c, "DELETE", "/v1/api/cart")
}

func (s *Server) handleGetProductInfo(c *fiber.Ctx) error {
	return s.handleProxy(c, "GET", "/v1/api/products/details/"+c.Params("sku"))
}

func (s *Server) handleProxy(c *fiber.Ctx, method, path string) error {
	log.Debug().Str("method", method).Str("path", path).Msg("handle proxy request")
	client := http.Client{}

	headers := c.GetReqHeaders()
	auth := headers["Authorization"]

	reqBody := bytes.NewBuffer(c.Body())

	req, err := http.NewRequest(method, s.serverUrl+path, reqBody)
	if err != nil {
		log.Error().Err(err).Str("method", method).Str("path", path).Str("upstream", s.serverUrl).Msg("Build new request error")
		return fiber.ErrBadRequest
	}

	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warn().Err(err).Str("method", method).Str("path", path).Str("upstream", s.serverUrl).Msg("Upstream connection error")
		return fiber.ErrBadGateway
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Str("method", method).Str("path", path).Str("upstream", s.serverUrl).Msg("Read resp body error")
		return fiber.ErrUnprocessableEntity
	}

	log.Debug().Str("method", method).Str("path", path).Int("status", resp.StatusCode).Msg("response")
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
	listenAddr   string
	serverUrl    string
	version      string
	allowOrigins string
}
