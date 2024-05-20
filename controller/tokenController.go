package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.niveussolutions.com/prajwal.amin/gop1/auth"
	"gitlab.niveussolutions.com/prajwal.amin/gop1/database"
	"gitlab.niveussolutions.com/prajwal.amin/gop1/model"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary generate token
// @ID gen-Token
// @Produce json
// @Param data body TokenRequest true "user credentials"
// @Success 200 {object} TokenRequest
// @Router /token [post]
func GenrateToken(c *gin.Context) {
	var request TokenRequest
	var user model.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	record := database.DB.Where("email=?", request.Email).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	credentialErr := user.ComparePassword(request.Password)
	if credentialErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": credentialErr.Error()})
		c.Abort()
		return
	}
	tokenString, err := auth.GenrateJWT(user.Name, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
