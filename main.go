package main

import (
	"golang_training/goutils"

	"github.com/gofiber/fiber/v2"
)

// API of CREATE READ UPDATE DELETE
func main() {

	// Connect to postgreSQL
	goutils.PostgresConnection("localhost", "postgres", "1Qaz2wsx", "postgres", "5432")
	// Initialize Fiber App
	app := fiber.New()

	//CREATE
	app.Post("/save", func(c *fiber.Ctx) error {
		requestModel := &User{}
		c.BodyParser(requestModel)

		goutils.DBConnect.Debug().Table("users").Create(requestModel)
		return c.JSON(fiber.Map{
			"Result": requestModel,
		})
	})

	//READ
	app.Get("/view", func(c *fiber.Ctx) error {
		responseModel := &[]UserResponse{}
		goutils.DBConnect.Debug().Raw("SELECT * FROM users").Find(responseModel)
		return c.JSON(responseModel)
	})

	//READ SPECIFIC USER
	app.Post("/get_user", func(c *fiber.Ctx) error {
		requestUser := &UserRequest{}
		c.BodyParser(requestUser)
		responseModel := &UserResponse{}
		goutils.DBConnect.Debug().Table("users").Find(responseModel, requestUser)
		return c.JSON(responseModel)
	})

	//UPDATE
	app.Post("/update/:id?", func(c *fiber.Ctx) error {
		var isExist int64
		var1 := c.Params("id")
		requestModel := &User{}
		c.BodyParser(requestModel)

		response := &UserResponse{}
		goutils.DBConnect.Debug().Raw("SELECT COUNT(name) FROM users WHERE id=?", var1).Count(&isExist)
		if isExist > 0 {
			if updatErr := goutils.DBConnect.Debug().Table("users").Where("id=?", var1).Updates(requestModel).Scan(response).Error; updatErr != nil {
				return c.JSON(fiber.Map{
					"Error": updatErr,
				})
			}
			return c.JSON(fiber.Map{
				"Result": response,
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
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       uint   `json:"userId"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
