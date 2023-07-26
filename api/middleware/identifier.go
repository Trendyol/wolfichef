package identifier

import "github.com/gofiber/fiber/v2"

type Config struct {
	Filter         func(*fiber.Ctx) bool
	ErrorHandler   fiber.ErrorHandler
	SuccessHandler fiber.Handler
	ContextKey     string
}

func New(config ...Config) fiber.Handler {
	var cfg Config
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).SendString("")
		}
	}
	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "identifier"
	}

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}

		identifier := c.GetReqHeaders()["X-User-Identifier"]

		if identifier == "" {
			return cfg.ErrorHandler(c, nil)
		}

		c.Locals(cfg.ContextKey, identifier)
		return cfg.SuccessHandler(c)
	}
}
