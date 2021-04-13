package style

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/images"
	"userstyles.world/models"
)

func StyleCreateGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("add", fiber.Map{
		"Title":  "Add userstyle",
		"User":   u,
		"Method": "add",
	})
}

func StyleCreatePost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	s := &models.Style{
		Name:        c.FormValue("name"),
		Description: c.FormValue("description"),
		Notes:       c.FormValue("notes"),
		Homepage:    c.FormValue("homepage"),
		Preview:     c.FormValue("previewUrl"),
		Code:        c.FormValue("code"),
		License:     strings.TrimSpace(c.FormValue("license", "No License")),
		Category:    strings.TrimSpace(c.FormValue("category", "unset")),
		UserID:      u.ID,
	}

	code := usercss.ParseFromString(c.FormValue("code"))
	if errs := usercss.BasicMetadataValidation(code); errs != nil {
		return c.Render("add", fiber.Map{
			"Title":  "Add userstyle",
			"User":   u,
			"Style":  s,
			"Method": "add",
			"Errors": errs,
		})
	}

	// Prevent adding multiples of the same style.
	err := models.CheckDuplicateStyleName(database.DB, s)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": err,
			"User":  u,
		})
	}

	var image multipart.File
	if s.Preview == "" {
		if ff, _ := c.FormFile("preview"); ff != nil {
			image, err = ff.Open()
			if err != nil {
				log.Println("Opening image , err:", err)
				return c.Render("err", fiber.Map{
					"Title": "Internal server error.",
					"User":  u,
				})
			}
		}
	}
	s, err = models.CreateStyle(database.DB, s)
	if err != nil {
		log.Println("Style creation failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	if image != nil {
		ID := strconv.FormatUint(uint64(s.ID), 10)
		data, _ := io.ReadAll(image)
		err = os.WriteFile(images.CacheFolder+ID+".originial", data, 0644)
		if err != nil {
			log.Println("Style creation failed, err:", err)
			return c.Render("err", fiber.Map{
				"Title": "Internal server error.",
				"User":  u,
			})
		}
		if s.Preview == "" {
			s.Preview = "https://userstyles.world/api/screenshot/" + ID + ".jpeg"
			database.DB.
				Model(new(models.Style)).
				Where("id", ID).
				Updates(s)
		}

	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}
