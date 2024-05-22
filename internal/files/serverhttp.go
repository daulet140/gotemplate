package files

import (
	"projectgenerator/internal/cfg"
)

func generateSrvr(config cfg.Config) (string, error) {
	file, err := generateHealthz(config)
	if err != nil {
		return file, err
	}

	file, err = generateServer(config)
	if err != nil {
		return file, err
	}

	file, err = generateValidator(config)
	if err != nil {
		return file, err
	}

	if config.Auth {
		file, err = generateMiddleware(config)
		if err != nil {
			return file, err
		}
	}
	return "", nil
}

func generateHealthz(config cfg.Config) (string, error) {
	var healthzTemplate = `package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (srv *server) healthz(ctx *gin.Context) {
	ctx.String(http.StatusOK, http.StatusText(http.StatusOK))
}
`
	err := generateFile(config.ProjectName+"/internal/http/healthz.go", healthzTemplate, config)
	if err != nil {
		return "healthz.go", err
	}
	return "", nil
}

func generateServer(config cfg.Config) (string, error) {
	var serverTemplate = `package http

import (
	exampleResourceV1 "{{.ProjectName}}/internal/http/v1/example"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	//swaggerfiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"`
	if config.Auth {
		serverTemplate += `
	authResourceV1 "{{.ProjectName}}/internal/http/v1/auth"
	"{{.ProjectName}}/internal/http/middleware"
	"{{.ProjectName}}/config"
`
	}

	serverTemplate += `
	"{{.ProjectName}}/internal/http/validator"
	managerV1 "{{.ProjectName}}/internal/manager/v1"
	"log"
	"net/http"
	"time"
)
const (
	maxAge   = 300
	v1Prefix = "/v1"
)

type server struct {
	appPort        string
	router         *gin.Engine	
`
	if config.Auth {
		serverTemplate += `	authConfig     config.Auth
	authManager             managerV1.AuthManager`
	}
	serverTemplate += `
	exampleManager managerV1.ExampleManager
}

func NewServer(appPort string, 
exampleManager managerV1.ExampleManager,`
	if config.Auth {
		serverTemplate += `authConfig config.Auth, manager managerV1.AuthManager`
	}

	serverTemplate += `) *server {
	return &server{
		appPort:        appPort,
		router:         gin.New(),
`
	if config.Auth {
		serverTemplate += `		authConfig:     authConfig,
		authManager: manager,`
	}
	serverTemplate += `
		exampleManager: exampleManager,
	}
}

func (srv *server) setupRouter() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println(err.Error())
	}

	srv.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Origin"},
		AllowCredentials: true,
		MaxAge:           maxAge,
	}))

	srv.router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{TimeFormat: time.RFC3339, UTC: false, SkipPaths: []string{"/", "/health", "/healthz"}}))
	srv.router.GET("/healthz", srv.healthz)

	v1 := srv.router.Group(v1Prefix)

	validate := validator.NewValidator()

`
	if config.Auth {
		serverTemplate += `
	authMiddleware := middleware.NewAuthController(srv.authConfig)
	authResourceV1.NewAuthResource(srv.authManager).Init(v1)
	exampleResourceV1.NewExampleResource(srv.exampleManager, validate, authMiddleware).Init(v1)`

	} else {
		serverTemplate += `
	exampleResourceV1.NewExampleResource(srv.exampleManager,validate).Init(v1)`
	}
	serverTemplate += `
	//srv.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func (srv *server) Run() {
	srv.setupRouter()
	if err := srv.router.Run(":" + srv.appPort); err != nil {
		zap.S().Fatal("Couldn't run HTTP server")
	}
}
`
	err := generateFile(config.ProjectName+"/internal/http/server.go", serverTemplate, config)
	if err != nil {
		return "server.go", err
	}
	return "", nil
}

func generateValidator(config cfg.Config) (string, error) {

	var validatorTemplate = `package validator

import (
	"github.com/go-playground/validator/v10"
	"{{.ProjectName}}/internal/models"
)

type Validator struct {
	Validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	return &Validator{
		Validator: v,
	}
}

func ValidatorErr(fe validator.FieldError) error {
	switch fe.Tag() {
	case "required":
		return models.ErrInvalidRequestBody
	case "email":
		return models.ErrInvalidEmail
	case "min", "max":
		switch fe.Field() {
		case "Username":
			return models.ErrInvalidUsername
		case "Password":
			return models.ErrInvalidPasswordFormat
		default:
		}

		return fe
	}

	return fe
}
`

	err := generateFile(config.ProjectName+"/internal/http/validator/validator.go", validatorTemplate, config)
	if err != nil {
		return "validator.go", err
	}
	return "", nil
}

func generateMiddleware(config cfg.Config) (string, error) {
	var middlewareTemplate = `package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"{{.ProjectName}}/config"
	authHttpRequest "{{.ProjectName}}/internal/clients/http/auth"
	"{{.ProjectName}}/internal/models"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	config config.Auth
}

type AuthController interface {
	CheckToken(ctx *gin.Context)
}

func NewAuthController(config config.Auth) AuthController {
	return &AuthMiddleware{
		config: config,
	}
}

func (a AuthMiddleware) CheckToken(ctx *gin.Context) {
	if !strings.Contains(ctx.GetHeader("Authorization"), "Bearer") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.BaseResponse{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)})

		return
	}

	username, err := authHttpRequest.TokenCheck(ctx, fmt.Sprintf("%s%s", a.config.ServiceUrl, a.config.ValidateTokenUrl),
		&models.TokenRequest{TokenHash: strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.BaseResponse{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)})
		return
	}

	ctx.Set(models.UsernameKey, username)
}
`

	err := generateFile(config.ProjectName+"/internal/http/middleware/middleware.go", middlewareTemplate, config)
	if err != nil {
		return "middleware.go", err
	}
	return "", nil
}
