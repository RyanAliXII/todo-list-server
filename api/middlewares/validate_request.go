package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidateRequestBody(this interface{}) gin.HandlerFunc { //YOU SHOULD PASS A STRUCT
	return func(c *gin.Context) {

		c.ShouldBindBodyWith(this, binding.JSON)

		validator := validator.New()
		validateErr := validator.Struct(this)
		if validateErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": validateErr.Error(),
				"data":    gin.H{},
			})
			return
		}

	}
}
