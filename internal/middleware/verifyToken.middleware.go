package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shortlink-backend/internal/jwttoken"
	"github.com/shortlink-backend/pkg"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) { // closure function
		token, ok := jwttoken.VerifyClientToken(ctx)
		if !ok {
			ctx.AbortWithStatus(401)
			return
		}

		var claims pkg.Claims
		if err := claims.VerifyJWT(token); err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
