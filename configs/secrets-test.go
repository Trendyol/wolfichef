package configs

// Secrets struct holds the secret values
type Secrets struct {
	Username        string
	Password        string
	slack_api_token string
	AccessToken     string
	APIKey          string
}
SLACK_API_KEY = "xoxb-263594206564-FGqddMF8t08v8N7Oq4i57vs1"
// GetSecrets returns the Secrets struct with secret values
func GetSecrets() Secrets {
	return Secrets{
		Username:        "dummy_username",
		Password:        "dummy_password",
		slack_api_token: "dummy_api_key",
		AccessToken:     "dummy_access_token",
		APIKey:          "dummy_api_key",
	}
}
