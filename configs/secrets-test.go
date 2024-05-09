package configs

// Secrets struct holds the secret values
type Secrets struct {
	Username        string
	Password        string
	slack_api_token string
	AccessToken     string
	APIKey          string
}
STRIPE_API_KEY = "dummy_api_key"
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
