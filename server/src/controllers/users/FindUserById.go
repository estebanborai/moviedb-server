package controllers

import (
	"fmt"

	data "github.com/estebanborai/songs-share-server/server/src/data"
	eh "github.com/estebanborai/songs-share-server/server/src/lib/error_handlers"
	models "github.com/estebanborai/songs-share-server/server/src/models"
	"github.com/gin-gonic/gin"
)

// FindUserByID returs the user with the specified ID
func FindUserByID(c *gin.Context, ID string) {
	var user models.User

	db := data.Connection(c)

	if userResult := db.Where(&models.User{ID: ID}).First(&user); userResult.Error == nil {
		c.JSON(200, user)
	} else {
		eh.NotFound(c, fmt.Sprintf("%s doesn't exists", ID))
	}
}
