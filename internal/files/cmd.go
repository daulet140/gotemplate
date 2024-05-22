package files

import (
	"projectgenerator/internal/cfg"
)

func fillCmd(folderCfg cfg.Config) (string, error) {
	var mainTemplate = `package main

import (
	"gitlab.com/golang-libs/mosk.git"
	"{{.ProjectName}}/config"
`
	if folderCfg.DBType != "" && folderCfg.DBType != "-" {
		mainTemplate += `
	"{{.ProjectName}}/internal/database"`
	}
	mainTemplate += `
	"{{.ProjectName}}/internal/http"
	managerV1 "{{.ProjectName}}/internal/manager/v1"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

func encoderConf() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func InitConfiguredLogger(externalConfig zap.Config) error {
	logger, err := externalConfig.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return nil
}

func initDefaultLogger(debug bool) error {
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}

	const samplingCongValue = 100

	var defaultJsonConfig = zap.Config{
		Level:         zap.NewAtomicLevelAt(level),
		Development:   debug,
		DisableCaller: true,
		Sampling: &zap.SamplingConfig{
			Initial:    samplingCongValue,
			Thereafter: samplingCongValue,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConf(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zap.AddStacktrace(zap.ErrorLevel)

	return InitConfiguredLogger(defaultJsonConfig)
}

// @title           "{{.ProjectName}}-backend
// @version         1.0
// @description     'desc"
func main() {
	err := initDefaultLogger(os.Getenv("DEBUG") != "")
	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}

	configuration := config.Config()

	mosk.LoadLocalConfig(configuration)
`
	if folderCfg.DBType != "" && folderCfg.DBType != "-" {
		mainTemplate += `
	db := database.NewDb(configuration.DbConfig.DriverName, configuration.DbConfig.GetConnectionStr())

	err = db.Connect()
	if err != nil {
		zap.S().Fatalf("[ERROR] Can't initialize logger: %v", err)
	}

	err = db.InitDB()
	if err != nil {
		zap.S().Fatalf("[ERROR] Can't initialize database: %v", err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			zap.S().Fatalf("[ERROR] Couldn't close connection %v", err)
		}

		zap.S().Debug("[DEBUG] Database connection closed!")
	}()
	`
	}

	if folderCfg.DBType != "" && folderCfg.DBType != "-" {
		mainTemplate += `
	exampleManager := managerV1.NewExampleManager(db.ExampleRepo())
`
	} else {
		mainTemplate += `
	exampleManager := managerV1.NewExampleManager()
`
	}

	if folderCfg.Auth {
		mainTemplate += `	authManager := managerV1.NewAuthorizationManager(configuration.Auth)
	http.NewServer(configuration.AppPort, exampleManager, configuration.Auth,authManager)
}`
	} else {
		mainTemplate += `
	http.NewServer(configuration.AppPort, exampleManager)
}`
	}
	err := generateFile(folderCfg.ProjectName+"/cmd/main.go", mainTemplate, folderCfg)
	return "main.go", err
}
