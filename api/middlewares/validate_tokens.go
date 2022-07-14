package middlewares

import (
	pasetoutils "flutter_task_app_server/pkg/paseto_utils"
	"fmt"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateTokens(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	FP, secureFPErr := ctx.Cookie("_secureFP")
	if secureFPErr != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Cannot find FP: Unauthorized Access",
		})
		return
	}
	if len(authorization) <= 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Missing authorization header",
		})
		return
	}
	authToken := strings.Split(authorization, " ")

	if len(authToken) < 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Missing Tokens",
		})
		return
	}

	claims, validateErr := pasetoutils.ValidateToken(authToken[1])
	if validateErr != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": validateErr.Error(),
		})
		return
	}

	if claims.TokenID != FP {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid FP: Unauthorized Access",
		})
		return
	}
	fmt.Println(claims)
	fmt.Println(FP)
	fmt.Println("VALIDATED")
	ctx.Next()
}
