package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"os"
	http_transport "refactored-robot/internal/controller/http-transport"
	"refactored-robot/internal/middleware"
	"refactored-robot/internal/models"
	"refactored-robot/internal/repository/mongodbrepo"
	"refactored-robot/internal/repository/postgresrepo"
	"refactored-robot/internal/service"
	"refactored-robot/package/database/mongodb"
	"refactored-robot/package/database/postgres"
)

func configInit() models.Config {
	yamlFile := "C:\\Users\\Max\\refactored-robot\\config\\config.yaml"

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	var config models.Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}
	return config

}

var db *gorm.DB
var mongo_db *mongo.Client

func init() {
	mongo_db = mongodb.NewDBConnection()
	mongodb.MigrateCollections()
	db = postgres.Init()
}

func main() {
	config := configInit()
	router := gin.Default()
	router.Use(gin.Recovery())

	userRepo := configRepo()

	userService := service.NewUserService(userRepo)
	userController := http_transport.NewUserController(userService)

	userRouter := router.Group("/user")
	{
		userRouter.POST("/", userController.Register)
		userRouter.DELETE("/delete/:id", userController.Delete)
		userRouter.GET("/get/:id", middleware.RequireAuth, userController.Get)
		userRouter.POST("/login", userController.Login)
		userRouter.POST("/setimg/:id", middleware.RequireAuth, userController.SetImage)
		userRouter.GET("/getimg/:id", userController.GetImage)
		userRouter.GET("/validate", middleware.RequireAuth, http_transport.Validate)
		userRouter.POST("/refresh", userController.RefreshAccessToken)

	}
	_ = router.Run(config.Port)

	/*
		-------------------------GRPC----------------------------
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
	*/

	//web_socket.Run() /*------------------------------WEB-SOCKET---------------------*/
}

func configRepo() service.IUserRepository {
	config := configInit()
	if config.Environment == "dev" {
		return mongodbrepo.NewUserRepository(mongo_db)
	}

	return postgresrepo.NewUserRepository(db)
}
