package main

import (
	"gateway/config"
	"gateway/pkg/authsvc"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"log"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	r := gin.Default()

	authsvc.RegisterRoutes(r, &c)

	r.Run(c.Port)
}
