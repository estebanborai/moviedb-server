package controllers

import (
	"strings"
	"time"

	data "github.com/estebanborai/songs-share-server/server/src/data"
	eh "github.com/estebanborai/songs-share-server/server/src/lib/error_handlers"
	models "github.com/estebanborai/songs-share-server/server/src/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

// Payload holds the spected properties to be sent from the Front-End
// at the time the user signs up
type Payload struct {
	models.User
	Password string
	Avatar   string
}

// SignUp creates a new user, hashes/salt the user password and stores user avatar
func SignUp(c *gin.Context) {
	var id = uuid.New().String()
	decodedPayload, _ := c.MultipartForm()

	birthDay, dateError := time.Parse(time.RFC3339, decodedPayload.Value["Birthday"][0])

	if dateError != nil {
		eh.BadRequest(c, "Invalid Date")
		return
	}

	user := models.User{
		ID:         id,
		UserName:   decodedPayload.Value["UserName"][0],
		FirstName:  decodedPayload.Value["FirstName"][0],
		LastName:   decodedPayload.Value["LastName"][0],
		Email:      decodedPayload.Value["Email"][0],
		Birthday:   birthDay,
		DateJoined: time.Now(),
	}

	db := data.Connection(c)

	// TODO: Automigrate on Init
	db.AutoMigrate(&models.User{})

	defer db.Close()

	db.NewRecord(user)

	if dbc := db.Create(&user); dbc.Error != nil {
		errorString := dbc.Error.Error()
		var isUserNameTaken = strings.Contains(errorString, "Error 1062")

		if isUserNameTaken == true {
			eh.BadRequest(c, "Username "+user.UserName+" is already taken.")
			return
		}

		eh.BadRequest(c, errorString)
		return
	}

	_, passwordError := CreatePassword(c, decodedPayload.Value["Password"][0], user.ID)
	if passwordError != nil {
		eh.ResponseWithError(c, 400, "Invalid Password")
		return
	}

	_, avatarUpdateError := UpdateAvatar(c, decodedPayload.File["Avatar"][0], user.ID)
	if avatarUpdateError != nil {
		eh.BadRequest(c, "Unable to gather image")
	}

	c.JSON(200, user)
}
