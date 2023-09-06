package http_transport

import (
	"context"
	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"net/http"
	"refactored-robot/internal/controller/rabbitmq_transport"
	"refactored-robot/internal/models"
	"refactored-robot/utils"
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
	amqURL := "amqp://guest:guest@localhost:5672/"
	controller, err := rabbitmq_transport.NewAMQController(amqURL)
	controller.CreateQueueAndPublishMessage("RefreshTokenHandlerQueue", "token succsesfully refreshed")

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

	const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
	tokenString, err := utils.CreateToken(user.Id, hmacSampleSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}
	refreshToken, err := utils.CreateToken(user.Id, hmacSampleSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	err = storeJWTTokenInRedis(user.Id, refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to send to Redis"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("AccessToken", tokenString, 60*15, "/", "localhost", false, true)
	c.SetCookie("RefreshToken", refreshToken, 3600*2, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"Authorization": tokenString})
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

func (ctrl *UserController) RefreshAccessToken(c *gin.Context) {
	const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
	cookie, err := c.Cookie("RefreshToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail"})
		return
	}

	err = utils.ValidateToken(cookie, hmacSampleSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	UserTokenId, err := getIDFromRedisByToken(cookie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User ID not found in Redis"})
		return
	}
	user, err := ctrl.userService.Get(UserTokenId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 8).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "access_token": tokenString})
}

func storeJWTTokenInRedis(id int, token string) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       1,                // Default DB
	})
	ctx := context.Background()
	err := redisClient.Set(ctx, token, strconv.Itoa(id), 2*time.Hour).Err()
	return err
}

func getIDFromRedisByToken(token string) (int, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	ctx := context.Background()

	idStr, err := redisClient.Get(ctx, token).Result()
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}
