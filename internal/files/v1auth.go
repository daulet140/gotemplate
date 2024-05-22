package files

import "projectgenerator/internal/cfg"

func generateV1Auth(config cfg.Config) (string, error) {
	var authTemplate = `package auth

import (
	"github.com/gin-gonic/gin"
	v1Manager "{{.ProjectName}}/internal/manager/v1"
	"{{.ProjectName}}/internal/models"
	"net/http"
	"strings"
)

type auth struct {
	authManager v1Manager.AuthManager
}

func NewAuthResource(authManager v1Manager.AuthManager) *auth {
	return &auth{
		authManager: authManager,
	}
}

func (a *auth) Init(router *gin.RouterGroup) {
	route := router.Group("/auth")
	route.POST("/login", a.login)
	route.POST("/refresh", a.refresh)

}

// refresh godoc
// @Summary Refresh token pair
// @Tags auth
// @Accept json
// @Produce json
// @Param token body models.TokenRequest true "refresh token"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /v1/auth/refresh [post]
func (a *auth) refresh(ctx *gin.Context) {
	body := new(models.TokenRequest)

	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: models.ErrInvalidRequestBody.Error(),
		})

		return
	}

	token, err := a.authManager.Refresh(ctx, body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.BaseResponse{
			Code:    http.StatusUnauthorized,
			Message: models.ErrInvalidToken.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, token)
}

// login godoc
// @Summary BAF user login endpoint
// @Tags auth
// @Accept json
// @Produce json
// @Security Authorization
// @param Authorization header string true "Client-Authorization"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.BaseResponse
// @Failure 500 {object} models.BaseResponse
// @Router /v1/auth/login [post]
func (a *auth) login(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")

	if authorization == "" || !strings.Contains(authorization, "Basic") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.BaseResponse{
			Code:    http.StatusUnauthorized,
			Message: http.StatusText(http.StatusUnauthorized),
		})

		return
	}

	token, err := a.authManager.Login(ctx, authorization)
	switch err {
	case nil:
	case models.ErrLoginBlocked:
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, models.BaseResponse{
			Code:    http.StatusTooManyRequests,
			Message: err.Error(),
		})

		return
	default:
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.BaseResponse{
			Code:    http.StatusUnauthorized,
			Message: models.ErrInvalidCredentials.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, token)
}

`
	filename := config.ProjectName + "/internal/http/v1/auth/auth.go"
	err := generateFile(filename, authTemplate, config)
	if err != nil {
		return filename, err
	}
	return "", nil
}
