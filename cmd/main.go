package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	httptransport "refactored-robot/internal/controller/http-transport"
	rpctransport "refactored-robot/internal/controller/rpc_transport"
	"refactored-robot/internal/middleware"
	"refactored-robot/internal/pb"
	"refactored-robot/internal/repository"
	"refactored-robot/internal/service"
	"refactored-robot/package/database/postgres"
)

func main() {

	DB := postgres.Init()

	router := gin.Default()
	router.Use(gin.Recovery())

	userRepo := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepo)
	userController := httptransport.NewUserController(userService)

	userRouter := router.Group("/user")
	{
		userRouter.POST("/", userController.Register)
		userRouter.DELETE("/delete/:id", userController.Delete)
		userRouter.GET("/get/:id", middleware.RequireAuth, userController.Get)
		userRouter.POST("/login", userController.Login)
		userRouter.POST("/setimg/:id", middleware.RequireAuth, userController.SetImage)
		userRouter.GET("/getimg/:id", userController.GetImage)
		userRouter.GET("/validate", middleware.RequireAuth, httptransport.Validate)
		userRouter.POST("/refresh", userController.RefreshAccessToken)

	}

	grpcPort := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	userRPCService := rpctransport.NewGRPCUserController(userService)

	pb.RegisterUserControllerServer(grpcServer, userRPCService)

	fmt.Printf("gRPC server listening on port %d", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
	_ = router.Run(":8888")
}
