package main

import (
	"github.com/gin-gonic/gin"
	"refactored-robot/internal/controller"
	"refactored-robot/internal/middleware"
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
	userController := controller.NewUserController(userService)

	userRouter := router.Group("/user")
	{
		userRouter.POST("/", userController.Register)
		userRouter.DELETE("/delete/:id", userController.Delete)
		userRouter.GET("/get/:id", userController.Get)
		userRouter.POST("/login", userController.Login)
		userRouter.POST("/setimg/:id", userController.SetImage)
		userRouter.GET("/getimg/:id", userController.GetImage)
		userRouter.GET("/validate", middleware.RequireAuth, controller.Validate)

	}

	_ = router.Run(":8888")
}
