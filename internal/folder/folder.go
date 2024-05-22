package folder

import (
	"fmt"
	"projectgenerator/internal/cfg"
)

type folder interface {
	CreateDirectories(config cfg.Config) error
}

type Manager struct {
	config cfg.Config
}

func NewFolderManager(config cfg.Config) *Manager { return &Manager{config: config} }

func (f *Manager) CreateDirectories() error {

	cmdDir := f.config.ProjectName + "/cmd"
	err := generateDirectory(cmdDir)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", cmdDir, err.Error())
	}

	configDir := f.config.ProjectName + "/config"
	err = generateDirectory(configDir)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", configDir, err.Error())
	}

	internalDir := f.config.ProjectName + "/internal"

	err = generateDirectory(internalDir)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", internalDir, err.Error())
	}

	if f.config.DBType != "" && f.config.DBType != "-" {
		databaseDir := internalDir + "/database"
		err = generateDirectory(databaseDir)
		if err != nil {

			return fmt.Errorf("Error creating project directory {%s} err:%s ", databaseDir, err.Error())
		}
	}

	clientsDir := internalDir + "/clients"
	err = generateDirectory(clientsDir)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", clientsDir, err.Error())
	}
	httpClientsDir := clientsDir + "/http"
	err = generateDirectory(httpClientsDir)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", httpClientsDir, err.Error())
	}
	if f.config.Auth {
		authDir := httpClientsDir + "/auth"
		err = generateDirectory(authDir)
		if err != nil {

			return fmt.Errorf("Error creating project directory {%s} err:%s ", authDir, err.Error())
		}
	}

	dirName, err := createHttpDir(internalDir, f.config)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", dirName, err.Error())
	}

	dirName, err = createManagerDir(internalDir)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", dirName, err.Error())
	}

	modelDir := internalDir + "/models"
	err = generateDirectory(modelDir)
	if err != nil {

		return fmt.Errorf("Error creating project directory {%s} err:%s ", modelDir, err.Error())
	}

	return nil
}
