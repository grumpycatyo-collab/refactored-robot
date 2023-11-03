package routes

import (
	"context"
	"fmt"
	"gateway/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RegisterRequestBody struct {
	Name      string `json:"name"`
	Password  string `json:"pass"`
	ImagePath string `json:"image"`
}

func Register(ctx *gin.Context, c pb.UserControllerClient) {
	b := RegisterRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.AddUser(context.Background(), &pb.RegisterUserRequest{
		Name:     b.Name,
		Password: b.Password,
		Image:    b.ImagePath,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

func GetUser(ctx *gin.Context, c pb.UserControllerClient) {
	userID := ctx.Param("id") // Assuming the user ID is passed as a URL parameter

	// Create a gRPC client using the provided connection

	// Convert userID to integer
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the gRPC GetUser method
	userResponse, err := c.GetUser(context.Background(), &pb.GetUserRequest{
		Id: int32(intUserID),
	})

	if err != nil {
		grpcStatus, _ := status.FromError(err)
		if grpcStatus.Code() == codes.NotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Convert the gRPC UserResponse to the desired HTTP response format
	httpResponse := struct {
		Id       int32  `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Image    string `json:"image"`
	}{
		Id:       userResponse.Id,
		Name:     userResponse.Name,
		Password: userResponse.Password,
		Image:    userResponse.Image,
	}

	ctx.JSON(http.StatusOK, httpResponse)
}

func GetRoles(ctx *gin.Context, c pb.UserControllerClient) {
	userID := ctx.Param("id") // Assuming the user ID is passed as a URL parameter

	// Create a gRPC client using the provided connection

	// Convert userID to integer
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the gRPC GetUser method
	roleResponse, err := c.GetRoles(context.Background(), &pb.GetRolesRequest{
		Id: int32(intUserID),
	})

	if err != nil {
		grpcStatus, _ := status.FromError(err)
		if grpcStatus.Code() == codes.NotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to retrieve role"})
		return
	}

	// Convert the gRPC UserResponse to the desired HTTP response format
	httpResponse := struct {
		Id    int32  `json:"id"`
		Roles string `json:"roles"`
	}{
		Id:    roleResponse.Id,
		Roles: roleResponse.Roles,
	}

	ctx.JSON(http.StatusOK, httpResponse)
}

func DeleteUser(ctx *gin.Context, c pb.UserControllerClient) {
	userID := ctx.Param("id") // Assuming the user ID is passed as a URL parameter

	// Create a gRPC client using the provided connection

	// Convert userID to integer
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call the gRPC GetUser method
	res, err := c.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		Id: int32(intUserID),
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Couldnt delete user"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type LoginRequestBody struct {
	Name     string `json:"name"`
	Password string `json:"pass"`
}

func Login(ctx *gin.Context, c pb.UserControllerClient) {
	b := LoginRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.Login(context.Background(), &pb.LoginRequest{
		Name:     b.Name,
		Password: b.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	tokenString := res.Token
	refreshToken := res.RefreshToken
	role := res.Role
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.Writer.Header().Set("role", role)
	ctx.Writer.Header().Set("Authorization", tokenString)
	ctx.Writer.Header().Set("RefreshToken", refreshToken)

	ctx.JSON(http.StatusOK, gin.H{"Authorization": tokenString})
	ctx.Next()
}

func SetImage(ctx *gin.Context, c pb.UserControllerClient) {

	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Read the image data from the request body
	imageData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image data"})
		return
	}

	// Call the gRPC service to set the user's image
	_, err = c.UploadImage(context.Background(), &pb.UploadImageRequest{
		Id:         int32(userID),
		ImageBytes: imageData,
	})

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func GetImage(ctx *gin.Context, c pb.UserControllerClient) {
	// Extract user ID from the URL parameter
	userIDStr := ctx.Param("id")
	fmt.Printf(userIDStr)

	// Convert the user ID string to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	// Call the UserService to get the user by ID
	user, err := c.GetUser(context.Background(), &pb.GetUserRequest{Id: int32(userID)})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Construct the full path to the user's image
	//imagePath := user.Image
	//fullImagePath := filepath.Join("C:\\Users\\Max\\refactored-robot\\web\\static", imagePath) // Replace with your image directory path

	// Serve the image file
	ctx.Status(http.StatusOK)
	ctx.File(user.Image)
}

func RefreshAccessToken(ctx *gin.Context, c pb.UserControllerClient) {
	const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
	cookie, err := ctx.Cookie("RefreshToken")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail"})
		return
	}
	tokenString, err := c.RefreshToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: cookie})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": tokenString})
}

func AddAdmin(ctx *gin.Context, c pb.UserControllerClient) {
	userIDStr := ctx.Param("id")

	// Convert the user ID string to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}
	res, err := c.AddAdmin(context.Background(), &pb.AddAdminRequest{Id: int32(userID)})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Couldnt delete user"})
		return
	}

	ctx.JSON(http.StatusOK, res)

}

func DeleteAdmin(ctx *gin.Context, c pb.UserControllerClient) {
	userIDStr := ctx.Param("id")

	// Convert the user ID string to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	res, err := c.DeleteAdmin(context.Background(), &pb.DeleteAdminRequest{Id: int32(userID)})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Couldnt delete user"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
