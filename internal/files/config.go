package files

import "github.com/daulet140/gotemplate/internal/cfg"

func generateConfigJsonFile(config cfg.Config) (string, error) {
	var configJsonTemplate = `{
		"app_port": "8080"`

	if config.Auth {
		configJsonTemplate += `,
		"auth": {
			"service_url": "https://auth2-staging.1cb.kz",
				"check_method_url": "/v1/check",
				"client_check_method_url": "/v1/client/check",
				"is_valid_method_url": "/v1/isvalid",
				"login_method_url": "/v1/login",
				"refresh_method_url": "/v1/refresh",
				"client_change_pass_url": "/v2/changePass",
				"client_name": "changeme",
				"client_secret": "changeme",
				"zone_name": "changeme"
		},
		"client_scope": {
			"area": "changeme",
				"topic": "*",
				"operation": "*"
		}
	`
	}
	if config.DBType != "" && config.DBType != "-" {
		configJsonTemplate += `,
		"db": {
			"host": "changeme",
				"port": "changeme",
				"user": "changeme",
				"pass": "changeme",
				"name": "changeme",
				"idle_conns": 3,
				"open_conns": 3,
				"driver_name": "{{ .DBType }}"
		}`
	}
	configJsonTemplate += `
}
`
	return "default.conf.json", generateFile(config.ProjectName+"/config/default.conf.json", configJsonTemplate, config)

}

func generateConfigGoFile(config cfg.Config) (string, error) {
	var configJsonTemplate = `package config

`
	if config.DBType != "" && config.DBType != "-" {
		configJsonTemplate += `import (
		dbConf "gitlab.com/golang-libs/databases.git/config"
	)
`
	}
	configJsonTemplate += "\ntype Configuration struct {\n\tAppPort  string          `json:\"app_port\" validate:\"nonzero\"`\n"

	if config.DBType != "" && config.DBType != "-" {
		configJsonTemplate += "\tDbConfig dbConf.DBConfig `json:\"db\" validate:\"nonzero\"`"
	}
	if config.Auth {
		configJsonTemplate += "\n\tAuth        Auth         `json:\"auth\"`\n }\n\ntype Auth struct {\n\tServiceUrl           string `json:\"service_url\"`\n\tRegisterUrl          string `json:\"register_url\"`\n\tPasswordSaveUrl      string `json:\"password_save_url\"`\n\tPasswordUpdateUrl    string `json:\"password_update_url\"`\n\tPasswordlessLoginUrl string `json:\"passwordless_login_url\"`\n\tRefreshUrl           string `json:\"refresh_url\"`\n\tLoginUrl             string `json:\"login_url\"`\n\tValidateTokenUrl     string `json:\"validate_token_url\"`\n\tDeleteUserUrl        string `json:\"delete_user_url\"`\n\tPasswordCheckUrl     string `json:\"password_check_url\"`\n"
	}
	configJsonTemplate += `
}

func Config() *Configuration {
		return &Configuration{}
}
`
	return "config.go", generateFile(config.ProjectName+"/config/config.go", configJsonTemplate, config)

}
