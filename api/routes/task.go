// package routes
package routes

// import (
// 	"flutter_task_app_server/api/definitions"
// 	"flutter_task_app_server/api/middlewares"
// 	"flutter_task_app_server/api/service"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/gin-gonic/gin/binding"
// )

// func InitializeTaskEndpoints(router *gin.RouterGroup, services *service.Services) {
// 	taskRoute := router.Group("/tasks")
// 	taskRoute.GET("/", func(c *gin.Context) {

// 		var tasks []definitions.Task
// 		res, err := services.MySQL.Query("Select * from task")
// 		fmt.Println(res)
// 		if err != nil {
// 			fmt.Println(err.Error())

// 		}

// 		for res.Next() {
// 			var task definitions.Task
// 			res.Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted)

// 			tasks = append(tasks, task)
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "List of Tasks",
// 			"data":    tasks,
// 		})
// 	})
// 	taskRoute.POST("/", middlewares.ValidateRequestBody(&definitions.CreateAndUpdateTaskBody{}), func(c *gin.Context) {
// 		var body definitions.CreateAndUpdateTaskBody = definitions.CreateAndUpdateTaskBody{}

// 		bindErr := c.ShouldBindBodyWith(&body, binding.JSON)

// 		if bindErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": "JSON failed to bind",
// 				"data":    gin.H{},
// 			})
// 			return
// 		}

// 		_, insertErr := services.MySQL.Exec("INSERT INTO task(title,description,isCompleted)VALUES(?,?,?)", body.Title, body.Description, 0)

// 		if insertErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": "Task failed to be created",
// 				"data":    gin.H{},
// 			})
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Task has been created",
// 			"data":    gin.H{},
// 		})
// 		return
// 	})
// 	taskRoute.PUT("/:id/edit/status/:status", func(c *gin.Context) {
// 		status, statusConvErr := strconv.Atoi(c.Param("status"))
// 		taskId, taskIdConvErr := strconv.Atoi(c.Param("id"))

// 		if statusConvErr != nil || taskIdConvErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": "Invalid URL Params",
// 				"data":    gin.H{},
// 			})
// 			return
// 		}

// 		_, updateErr := services.MySQL.Exec("Update task set isCompleted = ? where id = ?", status, taskId)

// 		if updateErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"message": "Task failed to update",
// 				"data":    gin.H{},
// 			})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Task has been updated",
// 			"data":    gin.H{},
// 		})

// 	})
// }
