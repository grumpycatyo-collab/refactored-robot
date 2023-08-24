package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"refactored-robot/internal/models"
)

type UserRepository struct {
	dbClient *gorm.DB
}

func NewUserRepository(dbClient *gorm.DB) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}

func (repo *UserRepository) Insert(user *models.User) error {
	err := repo.dbClient.Debug().
		Model(models.User{}).
		Create(user).Error
	if err != nil {
		log.Printf("failed to insest user in database: %v\n", err)
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(userID int) error {
	err := repo.dbClient.Debug().
		Model(models.User{}).
		Where("id = ?", userID).
		Delete(&models.User{}).Error
	if err != nil {
		log.Printf("failed to delete user from database: %v\n", err)
		return err
	}

	return nil
}

func (repo *UserRepository) Get(userID int) (*models.User, error) {
	user := &models.User{}
	err := repo.dbClient.Debug().
		Model(user).
		First(user, userID).Error
	if err != nil {
		log.Printf("failed to get user from database: %v\n", err)
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) GetUserByName(Name string) (*models.User, error) {
	user := &models.User{}
	err := repo.dbClient.Debug().
		Model(user).
		Where("name = ?", Name).
		First(user).Error
	if err != nil {
		log.Printf("failed to get user from database: %v\n", err)
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) CheckIfNameExists(name string) bool {
	var user models.User
	err := repo.dbClient.Debug().Model(models.User{}).Find(&user).Where("name = ?", name).Error
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (repo *UserRepository) UploadImage(userID int, image []byte) error {
	var user models.User
	result := repo.dbClient.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}

	user.Image = image
	result = repo.dbClient.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
