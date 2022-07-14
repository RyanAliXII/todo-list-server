package api

import (
	"database/sql"
	"flutter_task_app_server/api/repository"
	"flutter_task_app_server/api/routes"
	pasetoutils "flutter_task_app_server/pkg/paseto_utils"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func InitAPI(router *gin.Engine) {

	v1 := router.Group("api/1")
	{

		DB_DRIVER := os.Getenv("DB_DRIVER")
		DB_USERNAME := os.Getenv("DB_USERNAME")
		DB_PASSWORD := os.Getenv("DB_PASSWORD")
		DB_HOST := os.Getenv("DB_HOST")
		DB_PORT := os.Getenv("DB_PORT")
		DB_NAME := os.Getenv("DB_NAME")
		CONNECTION_STRING := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DB_USERNAME, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
		fmt.Println(CONNECTION_STRING)
		mySQL, err := sql.Open(DB_DRIVER, CONNECTION_STRING)

		if err != nil {
			fmt.Println("INIT API METHOD")
			fmt.Println(err.Error())
		}

		// routes.InitializeTaskEndpoints(v1, &services)
		// routes.InitializeTaskEndpoints(v1, &services)
		v1.POST("/auth/verify", func(ctx *gin.Context) {
			authorization := ctx.Request.Header.Get("Authorization")
			FP, secureFPErr := ctx.Cookie("_secureFP")
			if secureFPErr != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"message": "Cannot find FP: Unauthorized Access",
				})
				return
			}
			if len(authorization) <= 0 {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"message": "No Acess Token: Unauthorized Access ",
				})
				return
			}
			authToken := strings.Split(authorization, " ")

			if len(authToken) < 2 {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"message": "No Acess Token: Unauthorized Access ",
				})
				return
			}

			claims, validateErr := pasetoutils.ValidateToken(authToken[1])
			if validateErr != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"message": validateErr.Error(),
				})
				return
			}

			if claims.TokenID != FP {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"message": "Invalid FP: Unauthorized Access",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"message": "Authorized Access ",
			})
			return

		})

		var repos repository.Repositories = repository.Repositories{
			UserRepo: *repository.NewUserRepo(mySQL),
			TaskRepo: *repository.NewTaskRepo(mySQL),
		}

		routes.InitializeUsersEndpoints(v1, &repos)
	}

}
