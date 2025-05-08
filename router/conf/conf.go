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

// returns the status of the backend
func statusBackend(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"devices": base.Devices != nil,
		"monitor": base.Monitor,
	})
}

// returns the diveces available by description, or if there is none by name
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

// sets the monitor to the device with the index
func setMonitor(c *fiber.Ctx) error {
	index := c.Params("index")
	indexInt, _ := strconv.Atoi(index)
	base.Set_monitor(uint8(indexInt))
	return c.JSON(fiber.Map{
		"message": "Monitor set",
	})
}
