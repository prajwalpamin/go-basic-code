package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.niveussolutions.com/prajwal.amin/gop1/database"
	"gitlab.niveussolutions.com/prajwal.amin/gop1/model"
)

// @Summary add a new user
// @ID create-user
// @Produce json
// @Param data body model.User true "user data"
// @Success 200 {object} model.User
// @Router /register [post]
func RegisterUser(context *gin.Context) {
	user := model.User{}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := database.DB.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "Name": user.Name, "Email": user.Email})
}

// @Summary login as a user
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @ID user-login
// @Produce json
// @Param data body model.LoginUser true "user data"
// @Success 200 string login_successfull
// @Router secure/login [post]
func LoginUser(context *gin.Context) {
	reqUser := model.LoginUser{}
	dbUser := model.User{}
	if err := context.ShouldBindJSON(&reqUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := database.DB.Where("email=?", reqUser.Email).First(&dbUser).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := dbUser.ComparePassword(reqUser.Password); err != nil {
		context.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": dbUser.Name + " logged in successfully"})

}

// @Summary get all users in the database
// @ID get-all-users
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]
func Users(context *gin.Context) {
	users := []model.User{}
	database.DB.Find(&users)
	context.JSON(http.StatusOK, gin.H{"data": users})

}

// @Summary get a user by ID
// @ID get-user-id
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} model.User
// @Router /users/{id} [get]
func User(context *gin.Context) {
	user := model.User{}
	userId := context.Params.ByName("id")

	if err := database.DB.Where("id=?", userId).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": user})
}

// @Summary update the user information
// @ID update-user
// @Produce json
// @Param id path string true "user ID"
// @Param user body model.User true "updating user"
// @Success 200 {object} model.User
// @Router /users/{id} [put]
func Update(context *gin.Context) {
	user := model.User{}

	userId := context.Params.ByName("id")

	if err := database.DB.Where("id=?", userId).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	data := database.DB.Save(&user)

	context.JSON(http.StatusOK, gin.H{"Updated data": data})
}

// @Summary delete a user by ID
// @ID delete-user-by-id
// @Produce json
// @Param id path string true "todo ID"
// @Success 200 {object} model.User
// @Router /users/{id} [delete]
func DeleteUser(context *gin.Context) {
	id := context.Params.ByName("id")
	user := model.User{}

	data := database.DB.Where("id=?", id).Find(&user)

	if err := database.DB.Where("id=?", id).Delete(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"Deleted user ": data})

}

// @Summary get secured page
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @ID get page
// @Produce json
// @Success 200 string secure page loaded
// @Router /secure/page [get]
func SecuredPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "This is a JWT secured web page "})
}
