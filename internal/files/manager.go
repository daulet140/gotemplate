package files

import "projectgenerator/internal/cfg"

func generateManager(config cfg.Config) (string, error) {
	var managerTemplate = `package v1

import (
	"context"
	"{{.ProjectName}}/internal/database"
	"{{.ProjectName}}/internal/models"
)

type ExampleManager interface {
	// GetAll
	GetAll(ctx context.Context) ([]models.Example, error)
	// GetById
	GetById(ctx context.Context, exampleId int64) (models.Example, error)
	// Update
	Update(ctx context.Context, example *models.Example) error
	// Save
	Save(ctx context.Context, example *models.Example) error
	// Delete
	Delete(ctx context.Context, exampleId int64) error
}

type exampleManager struct {
	exampleRepo database.Example
}

func NewExampleManager(exampleRepo database.Example) ExampleManager {
	return &exampleManager{
		exampleRepo: exampleRepo,
	}
}

func (l *exampleManager) GetAll(ctx context.Context) ([]models.Example, error) {
	return l.exampleRepo.GetAll(ctx)
}

func (l *exampleManager) GetById(ctx context.Context, exampleId int64) (models.Example, error) {
	return l.exampleRepo.GetById(ctx, exampleId)
}

func (l *exampleManager) Update(ctx context.Context, example *models.Example) error {
	return l.exampleRepo.Update(ctx, example)
}

func (l *exampleManager) Save(ctx context.Context, example *models.Example) error {
	return l.exampleRepo.Save(ctx, example.Username, example.Password)
}

func (l *exampleManager) Delete(ctx context.Context, exampleId int64) error {
	return l.exampleRepo.Delete(ctx, exampleId)
}

`
	if config.DBType == "" || config.DBType == "-" {
		managerTemplate = `
package v1

type ExampleManager interface {
	}
type exampleManager struct {
}
func NewExampleManager() ExampleManager {
	return &exampleManager{}
}
`

	}
	err := generateFile(config.ProjectName+"/internal/manager/v1/example.go", managerTemplate, config)
	if err != nil {
		return "/manager/v1/example.go", err
	}
	return "", nil
}

func generateAuthManager(config cfg.Config) (string, error) {
	var authManagerTemplate = `package v1

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"{{.ProjectName}}/config"
	httpRequest "{{.ProjectName}}/internal/clients/http/auth"
	"{{.ProjectName}}/internal/models"
	"strings"
)

type AuthManager interface {
	Login(ctx context.Context, authorization string) (*models.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken *models.TokenRequest) (*models.LoginResponse, error)
}

type authorizationManager struct {
	authConfig config.Auth
}

func NewAuthorizationManager(authConfig config.Auth) AuthManager {
	return &authorizationManager{
		authConfig: authConfig,
	}
}

func (a *authorizationManager) Login(ctx context.Context, authorization string) (*models.LoginResponse, error) {
	tokenPair, err := httpRequest.Login(ctx, fmt.Sprintf("%s%s", a.authConfig.ServiceUrl, a.authConfig.LoginUrl), authorization)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return tokenPair, nil
}

func (a *authorizationManager) Refresh(ctx context.Context, refreshToken *models.TokenRequest) (*models.LoginResponse, error) {
	tokenPair, err := httpRequest.Refresh(ctx, fmt.Sprintf("%s%s", a.authConfig.ServiceUrl, a.authConfig.RefreshUrl), refreshToken)
	if err != nil {
		zap.S().Error(err)
		return nil, err

	}

	return tokenPair, nil
}

func parseBasicAuth(header string) (name, secret string, err error) {
	const prefix = "Basic "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", "", errors.New("Wrong auth length")
	}

	c, err := base64.StdEncoding.DecodeString(header[len(prefix):])
	if err != nil {
		return "", "", err
	}

	creds := strings.SplitN(string(c), ":", 2)
	if len(creds) < 2 {
		return "", "", errors.New("Wrong auth format")
	}
	return creds[0], creds[1], nil
}
`

	err := generateFile(config.ProjectName+"/internal/manager/v1/auth.go", authManagerTemplate, config)
	if err != nil {
		return "/manager/v1/auth.go", err
	}
	return "", nil
}
