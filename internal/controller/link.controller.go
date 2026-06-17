package controller

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/shortlink-backend/internal/binder"
	"github.com/shortlink-backend/internal/dto"
	"github.com/shortlink-backend/internal/jwttoken"
	"github.com/shortlink-backend/internal/response"
	"github.com/shortlink-backend/internal/service"
)

type LinkController struct {
	linkService *service.LinkService
}

func NewLinkController(linkService *service.LinkService) *LinkController {
	return &LinkController{
		linkService: linkService,
	}
}

// CreateLink godoc
// @Summary Create a shortlink
// @Description Create a shortlink with an optional custom slug
// @Tags Links
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param createLinkRequest body dto.CreateLinkRequest true "Create Link Request"
// @Success 201 {object} dto.Response "Created"
// @Failure 400 {object} dto.Response "Bad Request"
// @Failure 401 {object} dto.Response "Unauthorized"
// @Failure 409 {object} dto.Response "Conflict"
// @Failure 422 {object} dto.Response "Unprocessable Entity"
// @Failure 500 {object} dto.Response "Internal Server Error"
// @Router /links [post]
func (lc *LinkController) CreateLink(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var body dto.CreateLinkRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	res, err := lc.linkService.CreateLink(ctx.Request.Context(), claims.UserId, body)
	if err != nil {
		log.Println("Error: ", err.Error())
		if errors.Is(err, service.ErrDuplicateSlug) {
			response.JSONDuplicate(ctx, "Slug Already Used")
			return
		}
		if errors.Is(err, service.ErrDuplicateOriginalURL) {
			response.JSONDuplicate(ctx, "Original URL Already Used")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	res.ShortURL = buildShortURL(ctx, res.Slug)
	response.JSONCreated(ctx, res, "Link created successfully")
}

// ListLinks godoc
// @Summary Get user links
// @Description Get the authenticated user's links with pagination
// @Tags Links
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" minimum(1)
// @Param limit query int false "Items per page" minimum(1) maximum(100)
// @Success 200 {object} dto.Response "OK"
// @Failure 400 {object} dto.Response "Bad Request"
// @Failure 401 {object} dto.Response "Unauthorized"
// @Failure 422 {object} dto.Response "Unprocessable Entity"
// @Failure 500 {object} dto.Response "Internal Server Error"
// @Router /links [get]
func (lc *LinkController) ListLinks(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var query dto.ListLinksQuery
	if err := binder.BindFormat(ctx, &query, binding.Query); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	res, err := lc.linkService.ListLinks(ctx.Request.Context(), claims.UserId, query)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	for i := range res.Data {
		res.Data[i].ShortURL = buildShortURL(ctx, res.Data[i].Slug)
	}

	response.JSONSuccess(ctx, res, "Links retrieved successfully")
}

// DeleteLink godoc
// @Summary Delete a link
// @Description Soft delete one authenticated user's link
// @Tags Links
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Link ID" minimum(1)
// @Success 204 "No Content"
// @Failure 400 {object} dto.Response "Bad Request"
// @Failure 401 {object} dto.Response "Unauthorized"
// @Failure 404 {object} dto.Response "Not Found"
// @Failure 422 {object} dto.Response "Unprocessable Entity"
// @Failure 500 {object} dto.Response "Internal Server Error"
// @Router /links/{id} [delete]
func (lc *LinkController) DeleteLink(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var uri dto.LinkURI
	if err := ctx.ShouldBindUri(&uri); err != nil {
		log.Println("Error: ", err.Error())
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	linkId, err := strconv.ParseInt(uri.ID, 10, 64)
	if err != nil || linkId < 1 {
		response.JSONUnprocessableEntity(ctx, map[string]string{
			"id": "The id field must be greater than or equal to 1.",
		})
		return
	}

	if err := lc.linkService.DeleteLink(ctx.Request.Context(), claims.UserId, linkId); err != nil {
		log.Println("Error: ", err.Error())
		if errors.Is(err, service.ErrLinkNotFound) {
			response.JSONNotFound(ctx, "Link Not Found")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONNoContent(ctx)
}

func buildShortURL(ctx *gin.Context, slug string) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	if forwardedProto := ctx.GetHeader("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = strings.TrimSpace(strings.Split(forwardedProto, ",")[0])
	}

	return fmt.Sprintf("%s://%s/%s", scheme, ctx.Request.Host, slug)
}
