package core

type globalConfig struct {
	ServerModuleID     string `yaml:"server_module" env:"SERVER_MODULE"`
	ServerlessModuleID string `yaml:"serverless_module" env:"SERVER_MODULE"`
}

// Config utilized for global config values
var Config globalConfig = globalConfig{
	ServerlessModuleID: "aws",
	ServerModuleID:     "aws",
}
