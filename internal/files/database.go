package files

import (
	"github.com/daulet140/gotemplate/internal/cfg"
)

func generateDatabase(config cfg.Config) (string, error) {

	filename, err := generateDatabaseInit(config)
	if err != nil {
		return filename, err
	}

	filename, err = generateExampleRepo(config)
	if err != nil {
		return filename, err
	}

	return "", nil
}

func generateDatabaseInit(config cfg.Config) (string, error) {
	var iniDatabaseTemplate = `
package database

import (
	"database/sql"
`

	switch config.DBType {
	case "pgx":
		iniDatabaseTemplate = iniDatabaseTemplate + `_ "github.com/jackc/pgx/stdlib"`
	case "mysql":
		iniDatabaseTemplate = iniDatabaseTemplate + `_ "github.com/go-sql-driver/mysql"`
	case "mariadb":
		iniDatabaseTemplate = iniDatabaseTemplate + `_ "github.com/go-sql-driver/mysql"`
	case "sqlite":
		iniDatabaseTemplate = iniDatabaseTemplate + `_ "github.com/mattn/go-sqlite3"`
	default:
		iniDatabaseTemplate = iniDatabaseTemplate + `_ "database/sql"`
	}

	iniDatabaseTemplate = iniDatabaseTemplate + `
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DataStore interface {
	Base
	ExampleRepo() Example
}

type Base interface {
	Connect() error
	InitDB() error
	Close() error
}

type DB struct {
	driver, dbStr     string
	db                *sqlx.DB
	exampleRepo        Example
}

func NewDb(driver, dbStr string) DataStore {
	return &DB{
		driver: driver,
		dbStr:  dbStr,
	}
}

func (d *DB) ExampleRepo() Example {
	if d.exampleRepo == nil {
		d.exampleRepo = NewExampleRepo(d.db)
	}

	return d.exampleRepo
}


func (d *DB) Connect() error {
	var err error

	d.db, err = sqlx.Open(d.driver, d.dbStr)
	if err != nil {
		zap.S().Errorf("[ERROR] Failed to create database connection %s", err)

		return err
	}

	err = d.db.Ping()
	if err != nil {
		zap.S().Errorf("[ERROR] Could not ping database connection %s", err)

		return err
	}

	zap.S().Info("Established Database connection")

	return nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Stats() sql.DBStats {
	return d.db.Stats()
}

func (d *DB) InitDB() error {
	//file, err := os.Open("./config/db_script.sql")
	//if err != nil {
	//	log.Printf("err open file: %v", err)
	//	return err
	//}
	//
	//scanner := bufio.NewScanner(file)
	//scanner.Split(bufio.ScanLines)
	//
	//var lines []string
	//for scanner.Scan() {
	//	lines = append(lines, scanner.Text())
	//}
	//
	//file.Close()
	//
	//for _, line := range lines {
	//	_, err = d.db.Exec(line)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}
`
	return "database.go", generateFile(config.ProjectName+"/internal/database/database.go", iniDatabaseTemplate, config)
}

func generateExampleRepo(config cfg.Config) (string, error) {
	var iniExampleRepoTemplate = `package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"{{.ProjectName}}/internal/models"
)

// Example type
type Example interface {
	GetAll(ctx context.Context) ([]models.Example, error)
	GetById(ctx context.Context, exampleId int64) (models.Example, error)
	Update(ctx context.Context, example *models.Example) error
	Save(ctx context.Context, username,password string) error
	Delete(ctx context.Context, exampleId int64) error
}

type exampleRepo struct {
	db *sqlx.DB
}

// NewExampleRepo - return new lead repo
func NewExampleRepo(db *sqlx.DB) Example {
	return &exampleRepo{
		db: db,
	}
}

// GetAll - get all leads
func (l *exampleRepo) GetAll(ctx context.Context) ([]models.Example, error) {
	var examples []models.Example
	err := l.db.SelectContext(ctx, &examples, "SELECT * FROM example")
	return examples, err
}

// GetById - get example
func (l *exampleRepo) GetById(ctx context.Context, exampleId int64) (models.Example, error) {
	var example models.Example
	err := l.db.GetContext(ctx, &example, "SELECT * FROM example WHERE id = $1", exampleId)
	return example, err
}

// Update - update example
func (l *exampleRepo) Update(ctx context.Context, example *models.Example) error {
	query := "UPDATE example SET username = $1, password = $2 WHERE id = $4"
	_, err := l.db.ExecContext(ctx, query, example.Username, example.Password, example.Id)
	return err
}

// Delete - delete example
func (l *exampleRepo) Delete(ctx context.Context, exampleId int64) error {
	query := "DELETE FROM example WHERE id = $1"
	_, err := l.db.ExecContext(ctx, query, exampleId)
	return err
}

// Save - save example
func (l *exampleRepo) Save(ctx context.Context, username,password string) error {
	query := "INSERT INTO example (username,password) VALUES ($1, $2)"
	_, err := l.db.ExecContext(ctx, query, username,password)

	return err
}

`
	return "repo_example.go", generateFile(config.ProjectName+"/internal/database/repo_example.go", iniExampleRepoTemplate, config)
}
