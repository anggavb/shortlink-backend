package jwttoken

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shortlink-backend/pkg"
)

func GetClaims(ctx *gin.Context) (pkg.Claims, bool) {
	claimsValue, ok := ctx.Get("claims")
	if !ok {
		log.Println("Error: Claims not found in context")
		return pkg.Claims{}, false
	}

	claims, ok := claimsValue.(pkg.Claims)
	if !ok {
		log.Println("Error: Invalid claims type")
		return pkg.Claims{}, false
	}

	return claims, true
}
