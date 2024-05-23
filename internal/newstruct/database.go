package newstruct

import (
	"fmt"
)

func generateExampleRepo(data *Data) (string, error) {
	var iniExampleRepoTemplate = `package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"{{.ProjectPath}}/internal/models"
)

// {{.StructName}} type
type {{.StructName}} interface {
	GetAll(ctx context.Context) ([]models.{{.StructName}}, error)
	GetById(ctx context.Context, {{.LowerStructName}}Id int64) (models.{{.StructName}}, error)
	Update(ctx context.Context, {{.LowerStructName}} *models.{{.StructName}}) error
	Save(ctx context.Context, username,password string) error
	Delete(ctx context.Context, {{.LowerStructName}}Id int64) error
}

type {{.LowerStructName}}Repo struct {
	db *sqlx.DB
}

// New{{.StructName}}Repo - return new lead repo
func New{{.StructName}}Repo{{.StructName}}Repo(db *sqlx.DB) {{.StructName}} {
	return &{{.LowerStructName}}Repo{
		db: db,
	}
}

// GetAll - get all leads
func (l *{{.LowerStructName}}Repo) GetAll(ctx context.Context) ([]models.{{.StructName}}, error) {
	var {{.LowerStructName}}s []models.{{.StructName}}
	err := l.db.SelectContext(ctx, &{{.LowerStructName}}s, "SELECT * FROM {{.LowerStructName}}")
	return {{.LowerStructName}}s, err
}

// GetById - get {{.LowerStructName}}
func (l *{{.LowerStructName}}Repo) GetById(ctx context.Context, {{.LowerStructName}}Id int64) (models.{{.StructName}}, error) {
	var {{.LowerStructName}} models.{{.StructName}}
	err := l.db.GetContext(ctx, &{{.LowerStructName}}, "SELECT * FROM {{.LowerStructName}} WHERE id = $1", {{.LowerStructName}}Id)
	return {{.LowerStructName}}, err
}

// Update - update {{.LowerStructName}}
func (l *{{.LowerStructName}}Repo) Update(ctx context.Context, {{.LowerStructName}} *models.{{.StructName}}) error {
	query := "UPDATE {{.LowerStructName}} SET username = $1, password = $2 WHERE id = $4"
	_, err := l.db.ExecContext(ctx, query, {{.LowerStructName}}.Username, {{.LowerStructName}}.Password, {{.LowerStructName}}.Id)
	return err
}

// Delete - delete {{.LowerStructName}}
func (l *{{.LowerStructName}}Repo) Delete(ctx context.Context, {{.LowerStructName}}Id int64) error {
	query := "DELETE FROM {{.LowerStructName}} WHERE id = $1"
	_, err := l.db.ExecContext(ctx, query, {{.LowerStructName}}Id)
	return err
}

// Save - save {{.LowerStructName}}
func (l *{{.LowerStructName}}Repo) Save(ctx context.Context, username,password string) error {
	query := "INSERT INTO {{.LowerStructName}} (username,password) VALUES ($1, $2)"
	_, err := l.db.ExecContext(ctx, query, username,password)

	return err
}

`
	databaseRepoFile := fmt.Sprintf("repo_%s.go", data.FileName)
	return databaseRepoFile, generateFile(fmt.Sprintf("generated_%s_internal/database/repo_%s.go", data.FileName, data.LowerStructName), iniExampleRepoTemplate, data)
}
