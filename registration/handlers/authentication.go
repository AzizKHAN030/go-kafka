package handlers

import (
	"log"

	"github.com/azizkhan030/go-kafka/registration/db"
	"github.com/azizkhan030/go-kafka/registration/utils"
	"github.com/gofiber/fiber/v2"
)

func LoginUser(c *fiber.Ctx) error {
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	var hashedPassword string
	var userID int

	query := "SELECT id, password FROM users WHERE email = $1"

	err := db.DB.QueryRow(c.Context(), query, data.Email).Scan(&userID, &hashedPassword)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if err := utils.VerifyPassword(hashedPassword, data.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email or password is invalid",
		})
	}

	token, err := utils.GenerateToken(userID)
	if err != nil {
		log.Fatalf("%s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
