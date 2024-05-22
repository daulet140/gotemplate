package main

import (
	"flag"
	"fmt"
	"github.com/daulet140/gotemplate/internal/cfg"
	"github.com/daulet140/gotemplate/internal/files"
	"github.com/daulet140/gotemplate/internal/folder"
	"os"
)

func main() {
	help := flag.String("help", "", "Help command")
	projectName := flag.String("name", "", "Project name")
	dbType := flag.String("db", "-", "Database type (postgres or mysql or sqlite3 or mssql)")
	withAuth := flag.String("auth", "false", "Use auth (true or false)")
	flag.Parse()

	if *projectName == "" || *help != "" {
		flag.PrintDefaults()
		return
	}

	config := cfg.Config{
		ProjectName: *projectName,
		DBType:      *dbType,
		Auth:        *withAuth == "true",
	}
	if config.DBType == "postgres" {
		config.DBType = "pgx"
	}
	if config.ProjectName == "" {
		fmt.Println("Project name is required")
		return
	}

	if _, err := os.Stat(config.ProjectName); err == nil {
		fmt.Println("Project directory already exists")
		return
	}

	err := os.Mkdir(config.ProjectName, 0755)
	if err != nil {
		fmt.Println("Error creating project directory:", err)
		return
	}

	folderManager := folder.NewFolderManager(config)
	fileManager := files.NewFileManager(config)

	err = folderManager.CreateDirectories()
	if err != nil {
		fmt.Println("Error creating project directories:", err)
		return
	}

	err = fileManager.GenerateFiles()
	if err != nil {
		fmt.Println("Error generating files:", err)
		return
	}

	fmt.Println("Project directory created successfully")
}
