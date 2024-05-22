package files

import (
	"fmt"
	"github.com/daulet140/gotemplate/internal/cfg"
	"log"
)

type files interface {
	GenerateFiles(config cfg.Config) error
}

type FileManager struct {
	config cfg.Config
}

func NewFileManager(config cfg.Config) *FileManager {
	return &FileManager{config: config}
}

func (f *FileManager) GenerateFiles() error {
	file, err := fillModFile(f.config)
	if err != nil {
		log.Printf("Error creating file: %s err: %v", file, err)
		return err
	}
	file, err = fillCmd(f.config)
	if err != nil {
		log.Printf("Error creating file: %s err: %v", file, err)
		return err
	}
	file, err = generateConfigJsonFile(f.config)
	if err != nil {
		log.Printf("Error creating file: %s err: %v", file, err)
		return err
	}
	file, err = generateConfigGoFile(f.config)
	if err != nil {
		log.Printf("Error creating file: %s err: %v", file, err)
		return err
	}

	if f.config.DBType != "" && f.config.DBType != "-" {
		file, err = generateDatabase(f.config)
		if err != nil {
			fmt.Println("Error generating files:", err)
			return err
		}
		file, err = generateDBModels(f.config)
		if err != nil {
			fmt.Println("Error generating files:", err)
			return err
		}
	}

	file, err = generateCommonRequest(f.config)
	if err != nil {
		fmt.Println("Error generating files:", err)
		return err
	}
	if f.config.Auth {
		file, err = generateAuthClient(f.config)
		if err != nil {
			fmt.Println("Error generating files:", err)
			return err
		}

		file, err = generateAuthManager(f.config)
		if err != nil {
			fmt.Println("Error generating files:", err)
			return err
		}

		file, err = generateV1Auth(f.config)
		if err != nil {
			fmt.Println("Error generating files:", err)
			return err
		}

	}

	file, err = generateV1Example(f.config)
	if err != nil {
		fmt.Println("Error generating files:", err)
		return err
	}
	file, err = generateRequestModels(f.config)
	if err != nil {
		fmt.Println("Error generating files:", err)
		return err
	}
	file, err = generateResponseModels(f.config)
	if err != nil {
		fmt.Println("Error generating files:", err)
		return err
	}
	file, err = generateManager(f.config)
	if err != nil {
		fmt.Println("Error generating files:", err)
		return err
	}
	file, err = generateSrvr(f.config)
	if err != nil {
		fmt.Println("Error generating files:", err)
		return err
	}

	return nil
}
