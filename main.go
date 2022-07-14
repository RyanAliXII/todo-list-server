package main

import (
	"net/http"

	"flutter_task_app_server/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()
	godotenv.Load(".env")
	r.GET("/", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "index",
		})
		return
	})
	api.InitAPI(r)
	r.Run()
}
