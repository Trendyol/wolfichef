package config

type GitlabConfig struct {
	AppId       string `yaml:"AppId"`
	DeployToken string `yaml:"DeployToken"`
	DeployUser  string `yaml:"DeployUser"`
	Domain      string `yaml:"Domain"`
	RedirectUri string `yaml:"RedirectUri"`
	SecretKey   string `yaml:"SecretKey"`
}
