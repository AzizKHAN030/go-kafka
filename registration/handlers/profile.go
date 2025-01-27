package handlers

import (
	"context"

	"github.com/azizkhan030/go-kafka/registration/db"
	"github.com/azizkhan030/go-kafka/registration/models"
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	user := new(models.PublicUser)

	userId := c.Locals("userID")

	query := "SELECT * FROM users WHERE id=$1"
	row := db.DB.QueryRow(context.Background(), query, userId)
	if err := row.Scan(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
	})
}
