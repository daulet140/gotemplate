package newstruct

import "fmt"

func generateV1Struct(data *Data) (string, error) {
	var iniV1Template = `package {{.LowerStructName}}

import (
	"errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-playground/validator/v10"
	"{{.ProjectPath}}/internal/http/validator"
	v1Manager "{{.ProjectPath}}/internal/manager/v1"
	"{{.ProjectPath}}/internal/http/middleware"
	"{{.ProjectPath}}/internal/models"
	"net/http"
	"strconv"
)

type {{.LowerStructName}} struct {
	{{.LowerStructName}}Manager v1Manager.{{.StructName}}Manager
	validate       *validator.Validator	
}

func New{{.StructName}}Resource({{.LowerStructName}}Manager v1Manager.{{.StructName}}Manager, validate *validator.Validator) *{{.LowerStructName}} {
	return &{{.LowerStructName}}{
		{{.LowerStructName}}Manager: {{.LowerStructName}}Manager,
		validate:       validate,
	}
}

func (e *{{.LowerStructName}}) Init(router *gin.RouterGroup) {
	router.GET("", e.getAll)
	router.GET("/:{{.LowerStructName}}Id", e.getById)
	router.POST("", e.save)
	router.PUT("/:{{.LowerStructName}}Id", e.update)
	router.DELETE("/:{{.LowerStructName}}Id", e.delete)
}

func (e *{{.LowerStructName}}) getAll(ctx *gin.Context) {
	{{.LowerStructName}}s, err := e.{{.LowerStructName}}Manager.GetAll(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, {{.LowerStructName}}s)
}

func (e *{{.LowerStructName}}) getById(ctx *gin.Context) {
	{{.LowerStructName}}IdParam := ctx.Param("{{.LowerStructName}}Id")
	id, err := strconv.ParseInt({{.LowerStructName}}IdParam, 10, 64)

	{{.LowerStructName}}Response, err := e.{{.LowerStructName}}Manager.GetById(ctx, id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, {{.LowerStructName}}Response)
}

func (e *{{.LowerStructName}}) save(ctx *gin.Context) {
	{{.LowerStructName}}Body := new(models.{{.StructName}})

	if err := ctx.ShouldBindJSON({{.LowerStructName}}Body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := e.validate.Validator.Struct({{.LowerStructName}}Body)
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
	err = e.{{.LowerStructName}}Manager.Save(ctx, {{.LowerStructName}}Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, {{.LowerStructName}}Body)
}

func (e *{{.LowerStructName}}) update(ctx *gin.Context) {
	{{.LowerStructName}}Body := new(models.{{.StructName}})

	if err := ctx.ShouldBindJSON({{.LowerStructName}}Body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.BaseResponse{
			Code:    http.StatusBadRequest,
			Message: models.ErrInvalidRequestBody.Error(),
		})
		return
	}

	err := e.validate.Validator.Struct({{.LowerStructName}}Body)
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

	err = e.{{.LowerStructName}}Manager.Update(ctx, {{.LowerStructName}}Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.BaseResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, {{.LowerStructName}}Body)
}

func (e *{{.LowerStructName}}) delete(ctx *gin.Context) {
	{{.LowerStructName}}IdParam := ctx.Param("{{.LowerStructName}}Id")
	id, err := strconv.ParseInt({{.LowerStructName}}IdParam, 10, 64)
	err = e.{{.LowerStructName}}Manager.Delete(ctx, id)
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
	databaseRepoFile := fmt.Sprintf("generated_%s_internal/http/v1/%s/%s.go", data.FileName, data.LowerStructName, data.LowerStructName)
	err := generateFile(databaseRepoFile, iniV1Template, data)
	if err != nil {
		return databaseRepoFile, err
	}
	return "", nil
}
