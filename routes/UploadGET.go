package routes

import "github.com/gin-gonic/gin"

func UploadGET(res *gin.Context) {
	res.JSON(200, gin.H{"message": "You have requested the GET method"})
}
