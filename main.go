package main

import (
	"golang_training/goutils"

	"github.com/gofiber/fiber/v2"
)

// API of CREATE READ UPDATE DELETE
func main() {

	// Connect to postgreSQL
	goutils.ConnectDB("localhost", "postgres", "1Qaz2wsx", "postgres", "5432")
	// Create table
	goutils.MigrateModel(UserModel{})
	// Initialize Fiber App
	app := fiber.New()

	//CREATE
	app.Post("/save", func(c *fiber.Ctx) error {
		requestModel := &UserModel{}
		c.BodyParser(requestModel)

		goutils.DBConnect.Debug().Raw("INSERT INTO user_models(name) VALUES(?)", requestModel.Name).Scan(requestModel)
		return c.JSON(requestModel)
	})
	//READ
	app.Get("/view", func(c *fiber.Ctx) error {
		responseModel := &[]UserModel{}

		goutils.DBConnect.Debug().Raw("SELECT * FROM user_models").Find(responseModel)

		return c.JSON(responseModel)
	})
	//UPDATE
	app.Post("/update/:name?", func(c *fiber.Ctx) error {
		var1 := c.Params("name")
		requestModel := &UserModel{}
		c.BodyParser(requestModel)

		goutils.DBConnect.Debug().Raw("UPDATE user_models SET name=? WHERE name=?", requestModel.Name, var1).Scan(requestModel)

		return c.JSON(fiber.Map{
			"Result": requestModel,
		})
	})
	//DELETE
	app.Delete("/delete/:name", func(c *fiber.Ctx) error {
		var1 := c.Params("name")
		responseModel := &UserModel{}

		if deletErr := goutils.DBConnect.Debug().Raw("DELETE FROM user_models WHERE name=?", var1).Scan(responseModel).Error; deletErr != nil {
			c.JSON(fiber.Map{
				"Error": deletErr,
			})
		}
		return c.JSON(fiber.Map{
			"Result": "successfully deleted",
		})
	})

	// Application Port
	app.Listen(":3000")
}

// Model
type UserModel struct {
	Id   uint   `json:"-" gorm:"primaryKey"`
	Name string `json:"name"`
}
