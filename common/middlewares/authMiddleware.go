package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appconstant "github.com/salman-pathan/go-micro-arch/common/appConstant"
	"github.com/salman-pathan/go-micro-arch/common/auth"
)

const (
	BEARER_SCHEMA = "Bearer"
)

type AuthMiddleware struct {
	auth auth.Auth
}

func NewAuthMiddleware(auth auth.Auth) AuthMiddleware {
	return AuthMiddleware{auth: auth}
}

func (mw *AuthMiddleware) AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	tokenString := authHeader[len(BEARER_SCHEMA):]

	if len(tokenString) < 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"err": gin.H{
				"code":    403,
				"message": "Authentication token not found",
			},
		})
		return
	}

	payload, err := mw.auth.VerifyToken(tokenString)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"err": gin.H{
				"code":    403,
				"message": "Authentication failed",
			},
		})
		return
	}

	ctx.Set(appconstant.CTX_USER_ID, payload.UserId)
	ctx.Next()
}
