package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/search"
)

func Search(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	q := c.Query("q")
	s, _ := search.FindStylesByText(q)

	return c.Render("search", fiber.Map{
		"Title":  "Search",
		"User":   u,
		"Styles": s,
		"Value":  q,
		"Root":   c.OriginalURL() == "/search",
	})
}
