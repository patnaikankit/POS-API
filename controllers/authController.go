package controllers

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/models"
)

func Login(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"Message": "Invalid post request",
		})
	}

	//check if passcode is empty
	if data["passcode"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Passcode is required",
			"error":   map[string]interface{}{},
		})
	}
	var cashier models.Cashier
	db.DB.Where("id = ?", cashierId).First(&cashier)

	//check if cashier exist
	if cashier.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not found",
			"error":   map[string]interface{}{},
		})
	}

	if cashier.Passcode != data["passcode"] {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Passcode Not Match",
			"error":   map[string]interface{}{},
		})
	}

	// generating jwt token associated with the cashier
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Issuer":    strconv.Itoa(int(cashier.ID)),
		"ExpiresAt": time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Token Expired or invalid",
		})
	}

	cashierData := make(map[string]interface{})
	cashierData["token"] = tokenString

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashierData,
	})
}

func Logout(c *fiber.Ctx) error {
	cashierID := c.Params("cashierID")
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// check if passcode is empty
	if data["passcode"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Passcode is Required!",
		})
	}

	var cashier models.Cashier
	db.DB.Where("id=?", cashierID).First(&cashier)

	// check if cashier exist
	if cashier.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not found!",
		})
	}

	// expiring the jwt token stored in the cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Cashier LoggedOut Successfully!",
	})
}

func Passcode(c *fiber.Ctx) error {
	cashierID := c.Params("cashierID")
	var cashier models.Cashier
	db.DB.Select("id,name,passcode").Where("id=?", cashierID).First(&cashier)

	if cashier.Name == "" || cashier.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	cashierData := make(map[string]interface{})
	cashierData["passcode"] = cashier.Passcode

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashierData,
	})
}
