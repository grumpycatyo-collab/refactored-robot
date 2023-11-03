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

	user := r.Group("/user")
	user.POST("/", svc.Register)
	user.GET("/:id", middleware.AuthMiddleware(svc.Client), svc.GetUser)
	user.DELETE("/:id", svc.DeleteUser)
	user.POST("/login", svc.Login)
	user.POST("/refresh", svc.RefreshAccessToken)
	user.POST("/:id/photo", svc.SetImage)
	user.GET("/:id/photo", svc.GetImage)
	user.POST("/:id/admin", svc.AddAdmin)
	user.DELETE("/:id/admin", svc.DeleteAdmin)

	//user.GET("/:id/roles", svc.GetRoles)

	return svc
}

func (svc *ServiceClient) Register(ctx *gin.Context) {
	routes.Register(ctx, svc.Client)
}

func (svc *ServiceClient) GetUser(ctx *gin.Context) {
	routes.GetUser(ctx, svc.Client)
}

func (svc *ServiceClient) DeleteUser(ctx *gin.Context) {
	routes.DeleteUser(ctx, svc.Client)
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
func (svc *ServiceClient) AddAdmin(ctx *gin.Context) {
	routes.AddAdmin(ctx, svc.Client)
}

func (svc *ServiceClient) DeleteAdmin(ctx *gin.Context) {
	routes.DeleteAdmin(ctx, svc.Client)
}

//func (svc *ServiceClient) GetRoles(ctx *gin.Context) {
//	routes.GetRoles(ctx, svc.Client)
//}
