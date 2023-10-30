package config

type GitlabConfig struct {
	AppId       string `yaml:"AppId" env:"APP_ID"`
	DeployToken string `yaml:"DeployToken" env:"DEPLOY_TOKEN"`
	DeployUser  string `yaml:"DeployUser" env:"DEPLOY_USER"`
	Domain      string `yaml:"Domain" env:"DOMAIN"`
	RedirectUri string `yaml:"RedirectUri" env:"REDIRECT_URI"`
	SecretKey   string `yaml:"SecretKey" env:"SECRET_KEY"`
}
