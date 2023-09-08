package service

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"log"
	"refactored-robot/internal/models"
	"refactored-robot/utils"
	"strconv"
	"time"
)

type IUserRepository interface {
	Insert(user *models.User) error
	CheckIfNameExists(name string) bool
	Delete(userID int) error
	Get(userID int) (*models.User, error)
	GetUserByName(Name string) (*models.User, error)
	SetImage(userID int, image []byte) error
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) Register(user *models.User) error {
	if svc.userRepo.CheckIfNameExists(user.Name) {
		return errors.New("user already exists")
	}

	hash, err := svc.generatePasswordHash(user.Password)
	if err != nil {
		return errors.New("err can't register user")
	}

	user.Password = hash

	return svc.userRepo.Insert(user)
}

func (svc *UserService) Delete(userID int) error {
	err := svc.userRepo.Delete(userID)
	if err != nil {
		log.Printf("failed to delete user from database: %v\n", err)
		return err
	}
	return nil
}

func (svc *UserService) Get(userID int) (*models.User, error) {
	user, err := svc.userRepo.Get(userID)
	if err != nil {
		log.Printf("failed to get user from database: %v\n", err)
		return nil, err
	}
	return user, nil
}

func (svc *UserService) GetUserByName(Name string) (*models.User, error) {
	user, err := svc.userRepo.GetUserByName(Name)
	if err != nil {
		log.Printf("failed to get user by name from database: %v\n", err)
		return nil, err
	}
	return user, nil
}

func (svc *UserService) generatePasswordHash(pass string) (string, error) {
	const salt = 14
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), salt)
	if err != nil {
		log.Printf("ERR: %v\n", err)
		return "", err
	}

	return string(hash), nil
}

func (svc *UserService) LoginUser(Name string, pass string) (string, string, error) {
	user, err := svc.userRepo.GetUserByName(Name)
	err = svc.ComparePasswordHash(user.Password, pass)
	if err != nil {
		return "", "", err
	}
	const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
	tokenString, _ := utils.CreateToken(user.Id, hmacSampleSecret)
	refreshToken, _ := utils.CreateToken(user.Id, hmacSampleSecret)
	err = storeJWTTokenInRedis(user.Id, refreshToken)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshToken, nil

}

func (svc *UserService) ComparePasswordHash(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		log.Printf("ERR: %v\n", err)
		return err
	}

	return nil
}

func (svc *UserService) SetImage(userID int, image []byte) error {
	return svc.userRepo.SetImage(userID, image)
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

func GetIDFromRedisByToken(token string) (int, error) {
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
