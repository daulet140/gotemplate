package main

import (
	"flag"
	"fmt"
	"os"
	"projectgenerator/internal/cfg"
	"projectgenerator/internal/files"
	"projectgenerator/internal/folder"
)

func main() {
	projectName := flag.String("name", "gotemplate", "Project name")
	dbType := flag.String("db", "postgresql", "Database type (postgres or mysql)")
	withAuth := flag.String("auth", "true", "Use auth (true or false)")

	flag.Parse()

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
