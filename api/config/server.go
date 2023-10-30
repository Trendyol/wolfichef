package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	identifier "trendyol.com/security/appsec/devsecops/wolfichef/api/middleware"
)

type ServerConfig struct {
	*fiber.App
	Host string `yaml:"Host" env:"HOST" env-default:"127.0.0.1"`
	Port string `yaml:"Port" env:"PORT" env-default:"8000"`
}

func (s *ServerConfig) Setup() {
	s.App = fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	s.App.Use(cors.New())
	s.App.Use(identifier.New())
}

func (s *ServerConfig) Serve() {
	err := s.Listen(fmt.Sprintf("%s:%s", s.Host, s.Port))
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
}
