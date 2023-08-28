package rpc_transport

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io/ioutil"
	"refactored-robot/internal/models"
	"refactored-robot/internal/pb"
	"refactored-robot/utils"
	//"refactored-robot/internal/pb"
)

type IUserService interface {
	Register(user *models.User) error
	Delete(userID int) error
	Get(userID int) (*models.User, error)
	GetUserByName(Name string) (*models.User, error)
	ComparePasswordHash(hash, pass string) error
	SetImage(userID int, image []byte) error
}

type GRPCUserController struct {
	userService IUserService
	pb.UnimplementedUserControllerServer
}

func NewGRPCUserController(userService IUserService) *GRPCUserController {
	return &GRPCUserController{
		userService: userService,
	}
}

func (ctrl *GRPCUserController) AddUser(ctx context.Context, request *pb.RegisterUserRequest) (*emptypb.Empty, error) {
	user := &models.User{
		Name:      request.Name,
		Password:  request.Password,
		ImagePath: request.Image,
	}

	err := ctrl.userService.Register(user)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (ctrl *GRPCUserController) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.UserResponse, error) {
	userID := request.GetId()
	user, err := ctrl.userService.Get(int(userID))
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	userResponse := &pb.UserResponse{
		// Assuming you have a mapping function to convert your User model to UserResponse
		Id:       request.Id,
		Name:     user.Name,
		Password: user.Password,
		Image:    user.ImagePath,
	}

	return userResponse, nil
}

func (ctrl *GRPCUserController) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	userID := request.GetId()
	err := ctrl.userService.Delete(int(userID))
	if err != nil {
		return nil, status.Error(codes.NotFound, "Couldnt delete")
	}
	return &emptypb.Empty{}, nil
}

func (ctrl *GRPCUserController) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	name := request.Name
	pass := request.Password

	user, err := ctrl.userService.GetUserByName(name)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Couldnt get the Name")
	}

	err = ctrl.userService.ComparePasswordHash(user.Password, pass)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Couldnt compare hash")
	}
	const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
	token, err := utils.CreateToken(user.Id, hmacSampleSecret) // You need to define this function
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Error in making token")
	}
	response := &pb.LoginResponse{
		Token: token,
	}

	return response, nil
}

func (ctrl *GRPCUserController) UploadImage(ctx context.Context, request *pb.UploadImageRequest) (*emptypb.Empty, error) {
	userID := request.GetId()
	imageData := request.GetImageBytes()

	// You can now call your userService's SetImage method with the userID and imageData
	// Assuming ctrl.userService is your actual service interface
	err := ctrl.userService.SetImage(int(userID), imageData)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil

}

func (ctrl *GRPCUserController) GetImage(ctx context.Context, request *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	userID := request.GetId()

	// Convert userID to int
	userIDInt := int(userID)

	user, err := ctrl.userService.Get(userIDInt)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Read image file
	imageData, err := ioutil.ReadFile(user.ImagePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error reading image file")
	}

	response := &pb.GetImageResponse{
		ImageData: imageData,
	}

	return response, nil
}

func (ctrl *GRPCUserController) RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (ctrl *GRPCUserController) mustEmbedUnimplementedUserControllerServer() {
	//TODO implement me
	panic("implement me")
}
