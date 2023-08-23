package controller

import (
	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"refactored-robot/internal/package/models"
	"strconv"
)

type IUserService interface {
	Register(user *models.User) error
	Delete(userID int) error
	Get(userID int) (*models.User, error)
	GetUserByName(Name string) (*models.User, error)
	ComparePasswordHash(hash, pass string) error
}

type UserController struct {
	userService IUserService
}

func NewUserController(userService IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	err := ctrl.userService.Register(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully create an user",
	})
}

func (ctrl *UserController) Delete(c *gin.Context) {
	userIDStr := c.Param("id")
	fmt.Printf(userIDStr)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = ctrl.userService.Delete(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctrl *UserController) Get(c *gin.Context) {
	userIDStr := c.Param("id")
	fmt.Printf(userIDStr)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	user, err := ctrl.userService.Get(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Login(c *gin.Context) {
	var loginInfo struct {
		Name     string `json:"name"`
		Password string `json:"pass"`
	}

	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	name := loginInfo.Name
	pass := loginInfo.Password

	user, err := ctrl.userService.GetUserByName(name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err1 := ctrl.userService.ComparePasswordHash(user.Password, pass)
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}
