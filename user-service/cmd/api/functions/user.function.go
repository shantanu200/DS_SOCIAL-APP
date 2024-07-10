package functions

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"user/cmd/api/models"
	"user/cmd/api/util"
	"user/cmd/database"
	rabbitmq "user/cmd/rabbitMQ"

	"github.com/goccy/go-json"
)

type LoginUserResponse struct {
	UserId      int    `json:"userId"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
	Name        string `json:"name"`
}

func CreateUser(doc map[string]interface{}) (*LoginUserResponse, error) {
	db := database.DB

	user := models.User{
		Name:         doc["name"].(string),
		Email:        doc["email"].(string),
		Password:     doc["password"].(string),
		IsHotUser:    false,
		CreationTime: time.Now(),
	}

	err := db.Create(&user).Error

	if err != nil {
		if exists := strings.Contains(err.Error(), "duplicate key"); exists {
			return nil, errors.New("email already exists")
		} else {
			return nil, err
		}
	}

	userToken, err := util.CreateAuthToken(user.Email, user.Id)

	if err != nil {
		return nil, err
	}

	return &LoginUserResponse{
		UserId:      user.Id,
		Email:       user.Email,
		Name:        user.Name,
		AccessToken: *userToken,
	}, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	db := database.DB

	var user *models.User

	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserById(Id float64) (*models.User, error) {
	db := database.DB

	var user *models.User

	err := db.Where("id = ?", Id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func LoginUser(email string, password string) (*LoginUserResponse, error) {
	db := database.DB
	userLogin, err := GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if userLogin.Password != password {
		return nil, errors.New("invalid Credential Passed")
	}

	err = db.Model(models.User{}).Where("email = ?", email).Updates(map[string]interface{}{"last_login": time.Now().UTC()}).Error
	if err != nil {
		return nil, err
	}

	userToken, err := util.CreateAuthToken(email, userLogin.Id)

	fmt.Println(userLogin)
	if err != nil {
		return nil, err
	}

	loginResponse := &LoginUserResponse{
		UserId:      userLogin.Id,
		Name:        userLogin.Name,
		Email:       email,
		AccessToken: *userToken,
	}

	return loginResponse, nil
}

func UpdateUser(userId int64, data models.User) (*models.User, error) {
	db := database.DB

	err := db.Model(models.User{}).Where("id = ?", userId).Updates(data).Error
	if err != nil {
		return nil, err
	}

	var user *models.User

	err = db.Model(models.User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}

	go func() {
		payload := rabbitmq.UpdateUserPayload{
			UserId: userId,
		}

		jsonMarshal, err := json.Marshal(payload)

		if err != nil {
			log.Fatalf("[UPDATE USER RABBITMQ ERROR] %s", err)
			return
		}

		fmt.Println(string(jsonMarshal))

		err = rabbitmq.PublishQueue(rabbitmq.EXCHANGENAME, "UPDATE_USER", string(jsonMarshal))

		if err != nil {
			log.Fatalf("[UPDATE USER RABBITMQ ERROR] %s", err)
			return
		}
	}()

	return user, nil
}
