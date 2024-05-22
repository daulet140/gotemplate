package files

import (
	"projectgenerator/internal/cfg"
)

func generateCommonRequest(config cfg.Config) (string, error) {
	var commonTemplate = `package http

import (
	"crypto/tls"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func MakeRequest(req *http.Request) ([]byte, error) {
	startTime := time.Now()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	cli := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	res, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		zap.S().Error(err.Error())

		return body, err
	}

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusNoContent {
		zap.S().Infof("Successfully made request to %s in %s",
			req.URL.String(), time.Since(startTime).String())

		return body, nil
	}

	zap.S().Errorf("Failed to make request to %s with status code %d and duration %s - %s",
		req.URL.String(), res.StatusCode, time.Since(startTime).String(), string(body))

	return body, errors.New(string(body))
}
`

	err := generateFile(config.ProjectName+"/internal/clients/http/common.go", commonTemplate, config)
	if err != nil {
		return "clients/http/common.go", err
	}
	return "clients/http/common.go", nil

}
func generateAuthClient(config cfg.Config) (string, error) {
	var authTemplate = `
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	request "{{.ProjectName}}/internal/clients/http"
	"{{.ProjectName}}/internal/models"
	"net/http"
)

func Login(ctx context.Context, url, authorization string) (*models.LoginResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authorization)

	return tokenRequest(req)
}

func Refresh(ctx context.Context, url string, tokenHash *models.TokenRequest) (*models.LoginResponse, error) {
	body, err := json.Marshal(tokenHash)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	return tokenRequest(req)
}

func TokenCheck(ctx context.Context, url string, tokenHash *models.TokenRequest) (any, error) {
	body, err := json.Marshal(tokenHash)
	if err != nil {
		return "", err
	}

	err = authRequest(ctx, url, body)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(tokenHash.TokenHash, claims, func(token *jwt.Token) (interface{}, error) {
		if err != nil {
			zap.S().Warn(err)
		}

		return token, nil
	})
	return claims["username"], nil
}

func tokenRequest(req *http.Request) (*models.LoginResponse, error) {
	body, err := request.MakeRequest(req)
	if err != nil {
		errMessage := new(models.BaseResponse)

		e := json.Unmarshal(body, errMessage)
		if e != nil {
			return nil, e
		}

		return nil, err
	}

	loginResponse := new(models.LoginResponse)

	err = json.Unmarshal(body, loginResponse)
	if err != nil {
		return nil, err
	}

	return loginResponse, nil
}

func authRequest(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	body, err = request.MakeRequest(req)
	if err != nil {
		errMessage := new(models.BaseResponse)

		e := json.Unmarshal(body, errMessage)
		if e != nil {
			return e
		}

		return errors.New(errMessage.Message)
	}

	return nil
}
`
	filename := config.ProjectName + "/internal/clients/http/auth/auth.go"
	err := generateFile(filename, authTemplate, config)
	if err != nil {
		return filename, err
	}
	return "", nil
}
