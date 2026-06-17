package controller

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/shortlink-backend/internal/binder"
	"github.com/shortlink-backend/internal/dto"
	"github.com/shortlink-backend/internal/jwttoken"
	"github.com/shortlink-backend/internal/response"
	"github.com/shortlink-backend/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.Response "Created"
// @Failure 400 {object} dto.Response "Bad Request"
// @Failure 422 {object} dto.Response "Unprocessable Entity"
// @Failure 409 {object} dto.Response "Conflict"
// @Failure 500 {object} dto.Response "Internal Server Error"
// @Router /auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	if err := ac.authService.RegisterUser(ctx.Request.Context(), body); err != nil {
		log.Println("Error: ", err.Error())
		if strings.Contains(err.Error(), "users_email_key") {
			response.JSONDuplicate(ctx, "Email Already Used")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONCreated(ctx, nil, "Register Successfully")
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.Response "OK"
// @Failure 400 {object} dto.Response "Bad Request"
// @Failure 401 {object} dto.Response "Unauthorized"
// @Failure 422 {object} dto.Response "Unprocessable Entity"
// @Failure 500 {object} dto.Response "Internal Server Error"
// @Router /auth [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var body dto.LoginRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	res, err := ac.authService.LoginUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		if strings.Contains(err.Error(), "wrong password") || strings.Contains(err.Error(), "no rows") {
			response.JSONUnauthorized(ctx, "Invalid email or password")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Login Successfully")
}

// Logout godoc
// @Summary Logout a user
// @Description Logout a user by invalidating the token
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 204 {object} nil "No Content"
// @Failure 401 {object} dto.Response "Unauthorized"
// @Failure 500 {object} dto.Response "Internal Server Error"
// @Router /auth/logout [delete]
func (ac *AuthController) Logout(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil || expiresAt == nil {
		log.Println("Error: ", err.Error())
		response.JSONUnauthorized(ctx, "Token expired, please login again!")
		return
	}

	if err := ac.authService.LogoutUser(ctx.Request.Context(), claims.UserId); err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONNoContent(ctx)
}
