package controller

import (
	"context"
	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"net/http"
	"refactored-robot/internal/models"
	"strconv"
	"time"
)

type IUserService interface {
	Register(user *models.User) error
	Delete(userID int) error
	Get(userID int) (*models.User, error)
	GetUserByName(Name string) (*models.User, error)
	ComparePasswordHash(hash, pass string) error
	SetImage(userID int, image []byte) error
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

		// Status 500 sau nu ?
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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

func (ctrl *UserController) GetImage(c *gin.Context) {
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

	c.Status(http.StatusOK)
	c.Writer.Write([]byte(user.ImagePath))
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 8).Unix(),
	})
	const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	err = storeJWTTokenInRedis(user.Id, tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to send to Redis"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*8, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func (ctrl *UserController) SetImage(c *gin.Context) {
	// Parse user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	imageData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image data"})
		return
	}

	err = ctrl.userService.SetImage(userID, imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func storeJWTTokenInRedis(id int, token string) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       1,                // Default DB
	})
	ctx := context.Background()
	err := redisClient.Set(ctx, strconv.Itoa(id), token, 8*time.Hour).Err() // Token expires in 24 hours
	return err
}
