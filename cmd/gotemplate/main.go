package main

import (
	"flag"
	"fmt"
	"github.com/daulet140/gotemplate/internal/cfg"
	"github.com/daulet140/gotemplate/internal/files"
	"github.com/daulet140/gotemplate/internal/folder"
	"github.com/daulet140/gotemplate/internal/newstruct"
	"os"
)

func main() {
	help := flag.String("help", "", "Help command")
	projectName := flag.String("name", "", "Project name")
	dbType := flag.String("db", "-", "Database type (postgres or mysql or sqlite3 or mssql)")
	withAuth := flag.String("auth", "false", "Use auth (true or false)")
	jsonStruct := flag.String("struct", "", "File path to JSON struct")

	flag.Parse()

	if (*projectName == "" && *jsonStruct == "") || *help != "" {
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
	if jsonStruct != nil && *jsonStruct != "" {
		newStruct := newstruct.NewJSONStruct(config)
		structMsg, err := newStruct.GenerateStructFromJSON(*jsonStruct)
		if err != nil {
			fmt.Printf("%v Error generating struct from JSON: %v", structMsg, err)
			return
		}
	} else {

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
}
