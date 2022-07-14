package routes

import (
	"flutter_task_app_server/api/middlewares"
	model "flutter_task_app_server/api/model"
	"flutter_task_app_server/api/repository"
	pasetoutils "flutter_task_app_server/pkg/paseto_utils"
	pwdutils "flutter_task_app_server/pkg/pwd_util"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func InitializeUsersEndpoints(router *gin.RouterGroup, repos *repository.Repositories) {
	userRoute := router.Group("users")

	userRoute.POST("/login", middlewares.ValidateRequestBody(&model.LoginUserBody{}), func(ctx *gin.Context) {
		loginCredentials := model.LoginUserBody{}
		ctx.ShouldBindBodyWith(&loginCredentials, binding.JSON)

		user, selectErr := repos.UserRepo.FindUserByEmail(loginCredentials.Email)
		if selectErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid username or password",
			})
			return
		}
		comparePassErr := pwdutils.ComparePassword(loginCredentials.Password, user.Password)
		if comparePassErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{

				"message": "Invalid username or password",
			})
			return
		}
		now := time.Now()
		paseto := pasetoutils.PasetoClaims{
			Issuer:     "task_app_server",
			Subject:    strconv.Itoa(user.Id),
			Expiration: (now.Add(4 * time.Hour)),
			NotBefore:  now,
			IssuedAt:   now,
		}

		token, tokenId, tokenErr := paseto.NewToken()

		if tokenErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{

				"message": "There was a problem in token creation",
			})
		}

		ctx.SetCookie("_secureFP", tokenId, 3600*24, "/", "http://localhost:8080", false, false)
		ctx.SetCookie("accessToken", token, 3600*24, "/", "http://localhost:8080", false, false)
		ctx.Header("access-token", token)
		ctx.Header("user-id", strconv.Itoa(user.Id))
		ctx.Header("secure-fp", tokenId)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User successfully logged in",
		})
		return
	})

	userRoute.POST("/register", middlewares.ValidateRequestBody(&model.User{}), func(ctx *gin.Context) {

		body := model.User{}
		ctx.ShouldBindBodyWith(&body, binding.JSON)

		hashedPassword, hashingErr := pwdutils.HashPassword(body.Password)
		body.Password = hashedPassword
		registerUserErr := repos.UserRepo.RegisterUser(&body)
		if hashingErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Your request cannot be process right now",
			})
			return
		}

		if registerUserErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": registerUserErr,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User successfully registered",
		})
		return
	})

	userRoute.Use(middlewares.ValidateTokens)
	userRoute.GET("/:userId/tasks", func(ctx *gin.Context) {
		id := ctx.Param("userId")
		intId, idErr := strconv.Atoi(id)
		if idErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid user id",
			})
		}
		tasks, fetchErr := repos.TaskRepo.GetTodosByUserId(intId)
		if fetchErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": fetchErr.Error(),
			})
		}
		fmt.Println(tasks)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "List of tasks for specific user",
			"data":    tasks,
		})
	})
	userRoute.POST("/:userId/tasks", middlewares.ValidateRequestBody(&model.CreateTaskBody{}), func(ctx *gin.Context) {
		task := model.CreateTaskBody{}
		id := ctx.Param("userId")
		intId, idErr := strconv.Atoi(id)
		if idErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid user id",
			})
			return
		}
		ctx.ShouldBindBodyWith(&task, binding.JSON)
		task.UserId = intId
		insertErr := repos.TaskRepo.NewTodo(task)
		if insertErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": insertErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "New to-do has been created",
		})
		return
	})

	userRoute.DELETE("/:userId/tasks/:taskId", func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		taskID := ctx.Param("taskId")
		intUserId, userIdErr := strconv.Atoi(userId)
		intTaskId, taskIdErr := strconv.Atoi(taskID)
		if userIdErr != nil || taskIdErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid URL Params",
			})
			return
		}
		deleteErr := repos.UserRepo.DeleteTaskByIdAndUserId(intTaskId, intUserId)
		if deleteErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": deleteErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Task has been deleted",
		})
		return
	})
	userRoute.PUT("/:userId/tasks/:taskId/", middlewares.ValidateRequestBody(&model.UpdateTaskBody{}), func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		taskID := ctx.Param("taskId")
		intUserId, userIdErr := strconv.Atoi(userId)
		intTaskId, taskIdErr := strconv.Atoi(taskID)
		updatedTask := model.UpdateTaskBody{}
		ctx.ShouldBindBodyWith(&updatedTask, binding.JSON)

		if userIdErr != nil || taskIdErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid URL Params",
			})
			return
		}

		updateErr := repos.TaskRepo.UpdateTask(updatedTask, intUserId, intTaskId)
		if updateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": updateErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Task has been updated",
		})

	})
	userRoute.PUT("/:userId/tasks/:taskId/status/:status", func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		taskID := ctx.Param("taskId")
		status := ctx.Param("status")
		intUserId, userIdErr := strconv.Atoi(userId)
		intTaskId, taskIdErr := strconv.Atoi(taskID)
		intStatus, statusErr := strconv.Atoi(status)
		if userIdErr != nil || taskIdErr != nil || statusErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid URL Params",
			})
			return
		}

		updateErr := repos.TaskRepo.UpdateTaskStatus(intTaskId, intUserId, intStatus)

		if updateErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": updateErr.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Task status has been updated",
		})

	})

}
