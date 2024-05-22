package files

import (
	"projectgenerator/internal/cfg"
)

func generateV1Example(config cfg.Config) (string, error) {
	var iniV1ExampleTemplate = `package example

import (
	"errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-playground/validator/v10"
	"{{.ProjectName}}/internal/http/validator"
	v1Manager "{{.ProjectName}}/internal/manager/v1"
	"{{.ProjectName}}/internal/http/middleware"
	"{{.ProjectName}}/internal/models"
	"net/http"
	"strconv"
)

type example struct {
	exampleManager v1Manager.ExampleManager
	validate       *validator.Validator
`
	if config.Auth {
		iniV1ExampleTemplate += `
		authMiddleware middleware.AuthController
		`
	}
	iniV1ExampleTemplate += `
	
}

func NewExampleResource(exampleManager v1Manager.ExampleManager, validate *validator.Validator`
	if config.Auth {
		iniV1ExampleTemplate += `, authMiddleware middleware.AuthController`
	}
	iniV1ExampleTemplate += `) *example {
	return &example{
		exampleManager: exampleManager,
		validate:       validate,
`
	if config.Auth {
		iniV1ExampleTemplate += ` authMiddleware: authMiddleware,`
	}
	iniV1ExampleTemplate += `
	}
}

func (e *example) Init(router *gin.RouterGroup) {
	router.GET("", e.getAll)
	router.GET("/:exampleId", e.getById)
	router.POST("", e.save)
	router.PUT("/:exampleId", e.update)
	router.DELETE("/:exampleId", e.delete)
}

func (e *example) getAll(ctx *gin.Context) {
	examples, err := e.exampleManager.GetAll(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, examples)
}

func (e *example) getById(ctx *gin.Context) {
	exampleIdParam := ctx.Param("exampleId")
	id, err := strconv.ParseInt(exampleIdParam, 10, 64)

	exampleResponse, err := e.exampleManager.GetById(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, exampleResponse)
}

func (e *example) save(ctx *gin.Context) {
	exampleBody := new(models.Example)

	if err := ctx.ShouldBindJSON(exampleBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := e.validate.Validator.Struct(exampleBody)
	if err != nil {
		var e validation.ValidationErrors

		if errors.As(err, &e) {
			for _, fe := range e {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, models.BaseResponse{
					Code:    http.StatusBadRequest,
					Message: validator.ValidatorErr(fe).Error(),
				})

				return
			}
		}
	}
	err = e.exampleManager.Save(ctx, exampleBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, exampleBody)
}

func (e *example) update(ctx *gin.Context) {
	exampleBody := new(models.Example)

	if err := ctx.ShouldBindJSON(exampleBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	err := e.validate.Validator.Struct(exampleBody)
	if err != nil {
		var e validation.ValidationErrors

		if errors.As(err, &e) {
			for _, fe := range e {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, models.BaseResponse{
					Code:    http.StatusBadRequest,
					Message: validator.ValidatorErr(fe).Error(),
				})

				return
			}
		}
	}

	err = e.exampleManager.Update(ctx, exampleBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, exampleBody)
}

func (e *example) delete(ctx *gin.Context) {
	exampleIdParam := ctx.Param("exampleId")
	id, err := strconv.ParseInt(exampleIdParam, 10, 64)
	err = e.exampleManager.Delete(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
`
	filename := config.ProjectName + "/internal/http/v1/example/example.go"
	err := generateFile(filename, iniV1ExampleTemplate, config)
	if err != nil {
		return filename, err
	}
	return "", nil
}
