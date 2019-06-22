package controllers

import (
	"strings"
	"time"

	eh "github.com/estebanborai/songs-share-server/server/src/lib/error_handlers"
	models "github.com/estebanborai/songs-share-server/server/src/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Payload struct {
	models.User
	Password string
	Avatar   string
}

func SignUp(c *gin.Context) {
	var id string = uuid.New().String()
	decodedPayload, _ := c.MultipartForm()

	birthDay, dateError := time.Parse(time.RFC3339, decodedPayload.Value["Birthday"][0])

	if dateError != nil {
		eh.BadRequest(c, "Invalid Date")
		return
	}

	user := models.User{
		Id:         id,
		UserName:   decodedPayload.Value["UserName"][0],
		FirstName:  decodedPayload.Value["FirstName"][0],
		LastName:   decodedPayload.Value["LastName"][0],
		Email:      decodedPayload.Value["Email"][0],
		Birthday:   birthDay,
		DateJoined: time.Now(),
	}

	db, err := gorm.Open("mysql", "root:root@tcp(songs-share-db)/songs-share?charset=utf8mb4&parseTime=True&loc=Local")

	if err != nil {
		eh.InternalServerError(c, "Unable to connect to the database")
		return
	} else {
		// TODO: Automigrate on Init
		db.AutoMigrate(&models.User{})
	}

	defer db.Close()

	db.NewRecord(user)

	if dbc := db.Create(&user); dbc.Error != nil {
		errorString := dbc.Error.Error()
		var isUserNameTaken bool = strings.Contains(errorString, "Error 1062")

		if isUserNameTaken == true {
			eh.BadRequest(c, "Username "+user.UserName+" is already taken.")
			return
		} else {
			eh.BadRequest(c, errorString)
			return
		}
	}

	_, passwordError := CreatePassword(decodedPayload.Value["Password"][0], user.Id)
	if passwordError != nil {
		eh.ResponseWithError(c, 400, "Invalid Password")
		return
	}

	_, avatarUpdateError := UpdateAvatar(c, decodedPayload.File["Avatar"][0], user.Id)
	if avatarUpdateError != nil {
		eh.BadRequest(c, "Unable to gather image")
	}

	c.JSON(200, user)
}
