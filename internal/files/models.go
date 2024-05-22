package files

import "projectgenerator/internal/cfg"

func generateDBModels(config cfg.Config) (string, error) {
	var modelTemplate = "package models\n\ntype Example struct {\n\tId int64 `json:\"id\" db:\"id\"`\n\tUsername string `json:\"username\" db:\"username\"`\n\tPassword string `json:\"password\" db:\"password\"`\n\tCreatedAt string `json:\"created_at\" db:\"created_at\"`\n\tUpdatedAt string `json:\"updated_at\" db:\"updated_at\"`\n}"
	err := generateFile(config.ProjectName+"/internal/models/db.go", modelTemplate, config)
	if err != nil {
		return "models/db.go", err
	}
	return "models/db.go", nil
}
func generateRequestModels(config cfg.Config) (string, error) {
	var modelTemplate = "package models\n\nimport \"errors\"\n\ntype ExampleRequest struct {\n\tId int64 `json:\"id\" db:\"id\"`\n\tUsername string `json:\"username\" db:\"username\"`\n\tPassword string `json:\"password\" db:\"password\"`\n\tCreatedAt string `json:\"created_at\" db:\"created_at\"`\n\tUpdatedAt string `json:\"updated_at\" db:\"updated_at\"`\n}\n\nvar (\n\tErrInvalidEmail              = errors.New(\"invalid email\")\n\tErrInvalidCredentials        = errors.New(\"invalid credentials\")\n\tErrInvalidRequestBody        = errors.New(\"malformed JSON\")\n\tErrProfileNotFound           = errors.New(\"profile not found\")\n\tErrInvalidRegistrationCode   = errors.New(\"invalid registration code\")\n\tErrConfirmationCodeExpired   = errors.New(\"confirmation code expired\")\n\tErrInvalidPasswordFormat     = errors.New(\"invalid password format\")\n\tErrProfileAlreadyExist       = errors.New(\"profile already exists\")\n\tErrInvalidUsername           = errors.New(\"invalid username\")\n\tErrInvalidUsernameOrEmail    = errors.New(\"invalid username or email\")\n\tErrProfileAlreadyExistsInBaf = errors.New(\"profile already exists in BAF\")\n\tErrProfileBlockedInBaf       = errors.New(\"profile blocked in BAF\")\n\tErrBiometryIncomplete        = errors.New(\"bio authentication was not completed\")\n\tErrEmailCodeWasNotConfirmed  = errors.New(\"email code was not confirmed\")\n\tErrInvalidEmailConfirmation  = errors.New(\"invalid confirmation code\")\n\tErrAttemptNotFound           = errors.New(\"attempt not found\")\n\tErrBioAttemptFailed          = errors.New(\"bio authentication attempt failed\")\n\tErrProfileBlocked            = errors.New(\"profile blocked\")\n\tErrLoginAttemptFailed        = errors.New(\"authentication attempt failed\")\n\tErrLoginBlocked              = errors.New(\"authentication blocked\")\n\tErrInvalidToken              = errors.New(\"invalid token\")\n\tErrBudgetNotChecked          = errors.New(\"budget not checked yet\")\n)" +
		"\n\nconst (\n\tUsernameKey                     = \"username\"\n\tBasicPrefix                     = \"Basic \"\n)\n"
	if config.Auth {
		modelTemplate += "\n\ntype TokenRequest struct {\n\tTokenHash string `json:\"token_hash\"`\n}\n" +
			"\ntype LoginResponse struct {\n\tExpiresIn        int    `json:\"expires_in\"`\n\tRefreshExpiresIn int    `json:\"refresh_expires_in\"`\n\tRefreshToken     string `json:\"refresh_token\"`\n\tTokenType        string `json:\"token_type\"`\n\tAccessToken      string `json:\"access_token\"`\n\tScope            string `json:\"scope\"`\n}\n"
	}
	err := generateFile(config.ProjectName+"/internal/models/requests.go", modelTemplate, config)
	if err != nil {
		return "models/requests.go", err
	}
	return "models/requests.go", nil
}
func generateResponseModels(config cfg.Config) (string, error) {
	var modelTemplate = "package models\n\ntype BaseResponse struct {\n\tCode    int    `json:\"code\"`\n\tMessage string `json:\"message\"`\n}\n" +
		""
	err := generateFile(config.ProjectName+"/internal/models/response.go", modelTemplate, config)
	if err != nil {
		return "models/response.go", err
	}
	return "models/response.go", nil
}
