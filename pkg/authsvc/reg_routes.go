package authsvc

import (
	"gateway/config"
	"gateway/middleware"
	"gateway/pkg/authsvc/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/user")
	routes.POST("/", svc.Register)
	routes.GET("/get/:id", middleware.AuthMiddleware(), svc.GetUser)
	routes.POST("/login", svc.Login)
	routes.POST("/setimg/:id", svc.SetImage)
	routes.GET("/gettimg/:id", svc.GetImage)
	routes.POST("/refresh", svc.RefreshAccessToken)
	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	routes.Register(ctx, svc.Client)
}

func (svc *ServiceClient) GetUser(ctx *gin.Context) {
	routes.GetUser(ctx, svc.Client)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	routes.Login(ctx, svc.Client)
}

func (svc *ServiceClient) SetImage(ctx *gin.Context) {
	routes.SetImage(ctx, svc.Client)
}

func (svc *ServiceClient) GetImage(ctx *gin.Context) {
	routes.GetImage(ctx, svc.Client)
}

func (svc *ServiceClient) RefreshAccessToken(ctx *gin.Context) {
	routes.RefreshAccessToken(ctx, svc.Client)
}
