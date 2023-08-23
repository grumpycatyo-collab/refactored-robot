package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"refactored-robot/internal/package/models"
)

type IUserRepository interface {
	Insert(user *models.User) error
	CheckIfNameExists(name string) bool
	Delete(userID int) error
	Get(userID int) (*models.User, error)
	GetUserByName(Name string) (*models.User, error)
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

func (svc *UserService) ComparePasswordHash(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		log.Printf("ERR: %v\n", err)
		return err
	}

	return nil
}
