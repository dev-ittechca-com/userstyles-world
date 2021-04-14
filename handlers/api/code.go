package api

import (
	"crypto/md5"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyleSource(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// Override updateURL field for Stylus integration.
	// TODO: Also override it in the database on demand?
	uc := usercss.ParseFromString(style.Code)
	url := "https://userstyles.world/api/style/" + id + ".user.css"
	uc.OverrideUpdateURL(url)

	// Count installs.
	_, err = models.AddStatsToStyle(database.DB, id, c.IP(), true)
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "Internal server error",
		})
	}

	c.Set("Content-Type", "text/css")
	return c.SendString(uc.SourceCode)
}

var normalizedHeaderETag = []byte("Etag")

func GetStyleEtag(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// TODO: add a possible update stat?

	// Follows the format "source code length - MD5 Checksum of source code"
	etagValue := fmt.Sprintf("\"%d-%x\"", len(style.Code), md5.Sum([]byte(style.Code)))

	// Set the value for "Etag" header
	c.Response().Header.SetCanonical(normalizedHeaderETag, []byte(etagValue))
	return nil
}
