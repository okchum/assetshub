package main

import (
	"assets-hub/controllers"
	"assets-hub/middlewares/config"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile(config.GetRoot("public_root",false), true)))

	// API
	authorized := router.Group("/404", gin.BasicAuth(config.GetAuth()))

	// Root
	authorized.GET("/", controllers.Test)
	// listAll
	authorized.GET("/listAll", controllers.ListAll)
	// list
	authorized.GET("/list", controllers.List)
	// upload
	router.POST("/upload", controllers.MultipleUploads)
	// move
	authorized.POST("/move", controllers.Move)
	// delete
	router.DELETE("/remove", controllers.Remove)

	router.Run(config.GetPort())
}
