package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hanshal101/fileUpload/routes"
)

func main() {
	api := gin.Default()
	api.GET("/upload", routes.UploadGET)
	api.POST("/upload", routes.UploadPOST)
	api.Run(":7894")
}
