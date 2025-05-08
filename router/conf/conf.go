package conf

import (
	"strconv"

	"github.com/JaredMijailRE/IEEE802Sniffer/base"
	"github.com/gofiber/fiber/v2"
)

func SetupConf(app *fiber.App) {

	app.Get("/status", statusBackend)
	app.Get("/devices", getDevices)

	app.Post("/monitor/:index", setMonitor)

}

func statusBackend(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"devices": base.Devices != nil,
		"monitor": base.Monitor,
	})
}

func getDevices(c *fiber.Ctx) error {
	deviceNames := make([]string, len(base.Devices))
	for i, device := range base.Devices {
		if device.Description == "" {
			deviceNames[i] = device.Name
		} else {
			deviceNames[i] = device.Description
		}
	}
	return c.JSON(deviceNames)
}

func setMonitor(c *fiber.Ctx) error {
	index := c.Params("index")
	indexInt, _ := strconv.Atoi(index)
	base.Set_monitor(uint8(indexInt))
	return c.JSON(fiber.Map{
		"message": "Monitor set",
	})
}
