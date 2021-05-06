package style

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/images"
	"userstyles.world/models"
)

func StyleEditGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	s, err := models.GetStyleByID(database.DB, p)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Check if logged-in user matches style author.
	if u.ID != s.UserID {
		return c.Render("err", fiber.Map{
			"Title": "User and style author don't match",
			"User":  u,
		})
	}

	return c.Render("add", fiber.Map{
		"Title":  "Edit userstyle",
		"Method": "edit",
		"User":   u,
		"Style":  s,
	})
}

func StyleEditPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	StyleID, t := c.Params("id"), new(models.Style)

	s, err := models.GetStyleByID(database.DB, StyleID)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Check if logged-in user matches style author.
	if u.ID != s.UserID {
		return c.Render("err", fiber.Map{
			"Title": "Users don't match",
			"User":  u,
		})
	}

	q := models.Style{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Notes:       c.FormValue("notes"),
		Code:        c.FormValue("code"),
		Preview:     c.FormValue("previewUrl"),
		Homepage:    c.FormValue("homepage"),
		License:     strings.TrimSpace(c.FormValue("license", "No License")),
		Category:    strings.TrimSpace(c.FormValue("category", "unset")),
		UserID:      u.ID,
	}

	if PreviewFormValue, _ := c.FormFile("preview"); PreviewFormValue != nil {
		image, err := PreviewFormValue.Open()
		if err != nil {
			log.Println("Opening image , err:", err)
			return c.Render("err", fiber.Map{
				"Title": "Internal server error.",
				"User":  u,
			})
		}
		data, _ := io.ReadAll(image)
		err = os.WriteFile(images.CacheFolder+StyleID+".original", data, 0644)
		if err != nil {
			log.Println("Style creation failed, err:", err)
			return c.Render("err", fiber.Map{
				"Title": "Internal server error.",
				"User":  u,
			})
		}
		// Either it's removed or it didn't exist.
		// So we don't care about the error.
		_ = os.Remove(images.CacheFolder + StyleID + ".jpeg")
		_ = os.Remove(images.CacheFolder + StyleID + ".webp")

		q.Preview = "https://userstyles.world/api/preview/" + StyleID + ".jpeg"
	}

	err = database.DB.
		Model(t).
		Where("id", StyleID).
		Updates(q).
		// GORM doesn't update non-zero values in structs.
		Update("mirror_code", c.FormValue("mirrorCode") == "on").
		Update("mirror_meta", c.FormValue("mirrorMeta") == "on").
		Error

	if err != nil {
		log.Println("Updating style failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect("/style/"+c.Params("id"), fiber.StatusSeeOther)
}
