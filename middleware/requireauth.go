package middleware

import (
	"context"
	"fmt"
	"gateway/pb"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
)

func AuthMiddleware(ctx pb.UserControllerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		//fmt.Println(tokenString)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			const hmacSampleSecret = "fjsdakfljsdfklasjfksdajlfa42134jkh4j23hfdsoaj"
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(hmacSampleSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		//fmt.Println(token.Raw)
		id, _ := claims["sub"].(float64)
		fmt.Println(id)
		//exp := claims["exp"].(float64)

		roleResponse, err := ctx.GetRoles(context.Background(), &pb.GetRolesRequest{
			Id: int32(id),
		})
		log.Println(roleResponse.Roles)
		//if roleResponse.Roles != "Admin" {
		//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Not admin"})
		//	c.Abort()
		//	return
		//}

		c.Next()
	}
}
