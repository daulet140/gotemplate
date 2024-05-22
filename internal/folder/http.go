package folder

import "projectgenerator/internal/cfg"

func createHttpDir(internalDir string, config cfg.Config) (string, error) {
	httpDir := internalDir + "/http"
	err := generateDirectory(httpDir)
	if err != nil {
		return httpDir, err
	}

	middlewareDir := httpDir + "/middleware"
	err = generateDirectory(middlewareDir)
	if err != nil {
		return middlewareDir, err
	}

	validatorDir := httpDir + "/validator"
	err = generateDirectory(validatorDir)
	if err != nil {
		return validatorDir, err
	}

	v1Dir := httpDir + "/v1"
	err = generateDirectory(v1Dir)
	if err != nil {
		return v1Dir, err
	}

	if config.Auth {
		authDir := v1Dir + "/auth"
		err = generateDirectory(authDir)
		if err != nil {
			return authDir, err
		}
	}

	exampleDir := v1Dir + "/example"
	err = generateDirectory(exampleDir)
	if err != nil {
		return exampleDir, err
	}
	return "", err
}
