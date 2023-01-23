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
	goutils.MigrateModel(User{})
	// Initialize Fiber App
	app := fiber.New()

	//CREATE
	app.Post("/save", func(c *fiber.Ctx) error {
		requestModel := &User{}
		c.BodyParser(requestModel)

		goutils.DBConnect.Debug().Raw("INSERT INTO users(name) VALUES(?)", requestModel.Name).Scan(requestModel)
		return c.JSON(fiber.Map{
			"Result": requestModel,
		})
	})
	//READ
	app.Get("/view", func(c *fiber.Ctx) error {
		responseModel := &[]User{}
		goutils.DBConnect.Debug().Raw("SELECT * FROM users").Find(responseModel)
		return c.JSON(responseModel)
	})
	//UPDATE
	app.Post("/update/:name?", func(c *fiber.Ctx) error {
		var isExist int64
		var1 := c.Params("name")
		requestModel := &User{}
		c.BodyParser(requestModel)

		goutils.DBConnect.Debug().Raw("SELECT COUNT(name) FROM users WHERE name=?", var1).Count(&isExist)
		if isExist > 0 {
			if updatErr := goutils.DBConnect.Debug().Raw("UPDATE users SET name=? WHERE name=?", requestModel.Name, var1).Scan(requestModel).Error; updatErr != nil {
				return c.JSON(fiber.Map{
					"Error": updatErr,
				})
			}
			return c.JSON(fiber.Map{
				"Result": requestModel,
			})
		}

		return c.JSON(fiber.Map{
			"Result": "Name not found",
		})
	})
	//DELETE
	app.Delete("/delete/:name", func(c *fiber.Ctx) error {
		var1 := c.Params("name")
		responseModel := &User{}

		if deletErr := goutils.DBConnect.Debug().Raw("DELETE FROM users WHERE name=?", var1).Scan(responseModel).Error; deletErr != nil {
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
type User struct {
	Id   uint   `json:"-" gorm:"primaryKey"`
	Name string `json:"name"`
}
