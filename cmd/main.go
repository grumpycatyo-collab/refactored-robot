package main

import (
	"github.com/gin-gonic/gin"
	"refactored-robot/internal/controller"
	"refactored-robot/internal/package/database"
	"refactored-robot/internal/repository"
	"refactored-robot/internal/service"
)

func main() {

	DB := database.Init()

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
		userRouter.GET("/login", userController.Login)
		userRouter.POST("/uploadimg/:id", userController.UploadImage)
		userRouter.GET("/getimg/:id", userController.GetImage)

	}

	_ = router.Run(":8888")
}
