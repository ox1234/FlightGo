package web

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(){
	router := gin.Default()
	router.Use(cors.Default())
	SetRouter(router)
	router.Run(":8767")
}
