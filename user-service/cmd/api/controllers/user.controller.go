package controllers

import (
	"strings"
	"user/cmd/api/functions"
	"user/cmd/api/handlers"
	"user/cmd/api/middleware"
	"user/cmd/api/models"
	"user/cmd/s3"

	"github.com/gofiber/fiber/v2"
)

func CreateUserModel(c *fiber.Ctx) error {
	var user map[string]interface{}

	if err := c.BodyParser(&user); err != nil {
		return handlers.ErrorRouter(c, "Invalid body passed")
	}

	result, err := functions.CreateUser(user)

	if err != nil {
		if exists := strings.Contains(err.Error(), "email already exists"); exists {
			return handlers.ErrorRouter(c, "Email already exists")
		}
		return handlers.ErrorRouter(c, "Unable to create user")
	}

	return handlers.SuccessRouter(c, "User created successfully", result)
}

func LoginUserModel(c *fiber.Ctx) error {
	var user map[string]interface{}

	if err := c.BodyParser(&user); err != nil {
		return handlers.ErrorRouter(c, "Invalid body passed")
	}

	result, err := functions.LoginUser(user["email"].(string), user["password"].(string))

	if err != nil {
		return handlers.ErrorRouter(c, err.Error())
	}

	return handlers.SuccessRouter(c, "User logged in", result)
}

func GetUserByToken(c *fiber.Ctx) error {
	logId, err := middleware.GetLocalUser(c)

	if err != nil {
		return handlers.InvalidUserRouter(c)
	}

	result, err := functions.GetUserById(logId)

	if err != nil {
		return handlers.ErrorRouter(c, err.Error())
	}

	return handlers.SuccessRouter(c, "User Found", result)
}

func UpdateUser(c *fiber.Ctx) error {
	logId, err := middleware.GetLocalUser(c)

	if err != nil {
		return handlers.InvalidUserRouter(c)
	}

	form, err := c.MultipartForm()

	if err != nil {
		return handlers.ErrorRouter(c, "Invalid Request Body passed")
	}

	profileImage := form.File["profileImage"]

	var profileUrl string

	if profileImage != nil {
		result, err := s3.FileUploader(profileImage)

		if err != nil {
			return handlers.ErrorRouter(c, "Unable to upload image on storage bucket")
		}

		profileUrl = result[0]
	}

	var user = models.User{
		Name:         form.Value["name"][0],
		Email:        form.Value["email"][0],
		ProfileImage: profileUrl,
	}

	result, err := functions.UpdateUser(int64(logId), user)

	if err != nil {
		return handlers.ErrorRouter(c, "Unable to update user details")
	}

	return handlers.SuccessRouter(c, "User details updated successfully", result)

}
