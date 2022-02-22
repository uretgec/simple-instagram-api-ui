package main

import (
	"github.com/gofiber/fiber/v2"
)

type MagicItems struct {
	Links []string `json:"links"`
}

// Homepage
func home() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("templates/home", fiber.Map{
			"Title": "Reposter - Instagram Media Downloader",
		})
	}
}

// Parser
func parser() fiber.Handler {
	return func(c *fiber.Ctx) error {

		mItem := &MagicItems{}
		if err := c.BodyParser(mItem); err != nil {
			return c.JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		// Start InstaPost Parse Action
		instaItem := NewMagicItem()
		if !instaItem.validateUrl(mItem.Links[0]) {
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "Please send correct instagram url!",
			})
		}

		body, err := instaItem.makeRequest()
		if err != nil {
			return c.JSON(fiber.Map{
				"status":  false,
				"message": "Instagram url not found!",
			})
		}

		err = instaItem.convertRequest(body)
		if err != nil {
			return c.JSON(fiber.Map{
				"status":  false,
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status":  true,
			"result":  instaItem,
			"message": "Magic generated",
		})
	}
}

/*func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}*/
