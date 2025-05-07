package main

import (
	"encoding/json"
	"time"

	"github.com/JaredMijailRE/IEEE802Sniffer/base"
	"github.com/JaredMijailRE/IEEE802Sniffer/router/conf"
	"github.com/JaredMijailRE/IEEE802Sniffer/router/view"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:                 true,
		AppName:                 "IEEE-API-SNIFFER",
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		ReadTimeout:             5 * time.Second,
		WriteTimeout:            5 * time.Second,
		IdleTimeout:             120 * time.Second,
		ReadBufferSize:          4096,
		WriteBufferSize:         4096,
		EnableTrustedProxyCheck: false,
		DisableKeepalive:        false,
		ServerHeader:            "IEEE-API-SNIFFER",
		Immutable:               true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	base.Get_divices()

	conf.SetupConf(app)
	view.SetupView(app)

	app.Listen(":3000")
}
