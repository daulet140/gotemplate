package newstruct

import (
	"fmt"
)

func generateManager(data *Data) (string, error) {
	var managerTemplate = `package v1

import (
	"context"
	"{{.ProjectPath}}/internal/database"
	"{{.ProjectPath}}/internal/models"
)

type {{.StructName}}Manager interface {
	GetAll(ctx context.Context) ([]models.{{.StructName}}, error)
	GetById(ctx context.Context, {{.LowerStructName}}Id int64) (models.{{.StructName}}, error)
	Update(ctx context.Context, {{.LowerStructName}} *models.{{.StructName}}) error
	Save(ctx context.Context, {{.LowerStructName}} *models.{{.StructName}}) error
	Delete(ctx context.Context, {{.LowerStructName}}Id int64) error
}

type {{.LowerStructName}}Manager struct {
	{{.LowerStructName}}Repo database.{{.StructName}}
}

func New{{.StructName}}Manager({{.LowerStructName}}Repo database.{{.StructName}}) {{.StructName}}Manager {
	return &{{.LowerStructName}}Manager{
		{{.LowerStructName}}Repo: {{.LowerStructName}}Repo,
	}
}

func (l *{{.LowerStructName}}Manager) GetAll(ctx context.Context) ([]models.{{.StructName}}, error) {
	return l.{{.LowerStructName}}Repo.GetAll(ctx)
}

func (l *{{.LowerStructName}}Manager) GetById(ctx context.Context, {{.LowerStructName}}Id int64) (models.{{.StructName}}, error) {
	return l.{{.LowerStructName}}Repo.GetById(ctx, {{.LowerStructName}}Id)
}

func (l *{{.LowerStructName}}Manager) Update(ctx context.Context, {{.LowerStructName}} *models.{{.StructName}}) error {
	return l.{{.LowerStructName}}Repo.Update(ctx, {{.LowerStructName}})
}

func (l *{{.LowerStructName}}Manager) Save(ctx context.Context, {{.LowerStructName}} *models.{{.StructName}}) error {
	return l.{{.LowerStructName}}Repo.Save(ctx, {{.LowerStructName}}.Username, {{.LowerStructName}}.Password)
}

func (l *{{.LowerStructName}}Manager) Delete(ctx context.Context, {{.LowerStructName}}Id int64) error {
	return l.{{.LowerStructName}}Repo.Delete(ctx, {{.LowerStructName}}Id)
}

`

	databaseRepoFile := fmt.Sprintf("manager/v1/%s.go", data.FileName)
	return databaseRepoFile, generateFile(fmt.Sprintf("generated_%s_internal/manager/v1/%s.go", data.FileName, data.LowerStructName), managerTemplate, data)
}
