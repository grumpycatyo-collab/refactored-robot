package mongodbrepo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"refactored-robot/internal/models"
	"time"
)

type UserRepository struct {
	dbClient *mongo.Client
}

func NewUserRepository(dbClient *mongo.Client) *UserRepository {
	return &UserRepository{
		dbClient: dbClient,
	}
}
func (repo *UserRepository) Insert(user *models.User) error {
	collection := repo.dbClient.Database("mongodb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("failed to insert user into MongoDB: %v\n", err)
		return err
	}

	return nil
}

func (repo *UserRepository) Delete(userID int) error {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get a handle to the users collection in your MongoDB
	usersCollection := repo.dbClient.Database("mongodb").Collection("users")

	// Define a filter to find the user with the specified ID
	filter := bson.M{"id": userID}

	// Delete the user document from MongoDB
	result, err := usersCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Failed to delete user from MongoDB: %v\n", err)
		return err
	}

	// Check if no documents matched the filter
	if result.DeletedCount == 0 {
		return fmt.Errorf("No user with ID %d found", userID)
	}

	return nil
}

func (repo *UserRepository) Get(userID int) (*models.User, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get a handle to the users collection in your MongoDB
	usersCollection := repo.dbClient.Database("mongodb").Collection("users")

	// Define a filter to find the user with the specified ID
	filter := bson.M{"id": userID}

	// Find the user document in MongoDB
	var user models.User
	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User with ID %d not found", userID)
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetUserByName(name string) (*models.User, error) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get a handle to the users collection in your MongoDB
	usersCollection := repo.dbClient.Database("mongodb").Collection("users")

	// Define a filter to find the user by name
	filter := bson.M{"name": name}

	// Find the user document in MongoDB
	var user models.User
	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("User with name %s not found", name)
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) CheckIfNameExists(name string) bool {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get a handle to the users collection in your MongoDB
	usersCollection := repo.dbClient.Database("mongodb").Collection("users")

	// Define a filter to check if the name exists
	filter := bson.M{"name": name}

	// Find one user with the given name
	var user models.User
	err := usersCollection.FindOne(ctx, filter).Decode(&user)

	// If there's no error and a user is found, the name exists
	return err == nil
}

func (repo *UserRepository) SetImage(userID int, image []byte) error {
	// Create a MongoDB client and establish a connection to your MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use your MongoDB database and collection names
	databaseName := "mongodb"
	collectionName := "users"

	// Connect to the database and collection
	db := repo.dbClient.Database(databaseName)
	collection := db.Collection(collectionName)

	// Find the user by userID
	var user models.User
	filter := bson.M{"_id": userID}
	err := collection.FindOne(ctx, filter).Decode(&user)
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
	err = ioutil.WriteFile(imagePath, image, 0644)
	if err != nil {
		return err
	}

	// Update the user model with the image path
	user.ImagePath = imagePath

	// Update the user in the MongoDB collection
	update := bson.M{"$set": bson.M{"imagePath": user.ImagePath}}
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		// Clean up the saved image file if the database update fails
		_ = os.Remove(imagePath)
		return err
	}

	return nil
}
