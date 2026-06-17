package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shortlink-backend/internal/dto"
)

func JSONAbortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
		Success: false,
		Message: "Unauthorized",
	})
}

func JSONAbortInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
		Success: false,
		Message: "Internal Server Error",
	})
}
