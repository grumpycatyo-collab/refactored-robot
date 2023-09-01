package repository

import (
	"bytes"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"refactored-robot/internal/models"
	"time"
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

func (repo *UserRepository) SetImage(userID int, image []byte) error {

	var user models.User
	result := repo.dbClient.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}

	// Salvarea imaginii se face in servicii + ai doua apeluri la baza de date intr-o singura metoda
	// am spus ca, fiecare metoda a repositoriului face doar o singura chestie
	img, err := jpeg.Decode(bytes.NewReader(image))
	if err != nil {
		return err
	}

	// Generate a unique filename using timestamp
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("user_%d_%d.jpg", userID, timestamp)

	// Define the folder path to save images
	imageFolderPath := "C:\\Users\\Max\\refactored-robot\\web\\static" // Update this with your actual path

	// Create the image file
	imagePath := filepath.Join(imageFolderPath, filename)
	imageFile, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()

	// Write the JPEG image to the file
	if err := jpeg.Encode(imageFile, img, nil); err != nil {
		return err
	}

	// Update the user model with the image path
	user.ImagePath = imagePath
	result = repo.dbClient.Save(&user)
	if result.Error != nil {
		// Clean up the saved image file if database update fails
		_ = os.Remove(imagePath)
		return result.Error
	}

	return nil
}
