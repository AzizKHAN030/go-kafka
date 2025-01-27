package handlers

import (
	"context"

	"github.com/azizkhan030/go-kafka/registration/db"
	"github.com/azizkhan030/go-kafka/registration/middleware"
	"github.com/azizkhan030/go-kafka/registration/models"
	"github.com/azizkhan030/go-kafka/registration/utils"
	"github.com/gofiber/fiber/v2"
)

// Define /register handler
// Parse user registration request
// Validate the request payload
// Serialize the payload to JSON or Protobuf
// Produce message to Kafka topic
// Return success response to client

func SetupRoutes(app *fiber.App) {
	app.Post("/register", RegisterUser)
	app.Post("/login", LoginUser)

	protected := app.Group("/api", middleware.JWTMiddleware)
	protected.Get("/profile", GetProfile)
}

func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user.Password = hashedPassword
	query := "INSERT INTO users (name, email, password) VALUES($1, $2, $3) RETURNING id"
	row := db.DB.QueryRow(context.Background(), query, user.Name, user.Email, user.Password)
	if err := row.Scan(&user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	publicUser := models.PublicUser{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return c.Status(fiber.StatusCreated).JSON(publicUser)
}
