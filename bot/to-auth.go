package bot

import (
	"encoding/json"
	"log/slog"
	"os"
)

// ToAuth performs the authentication process for the Strava API.
// It makes a POST request to the Strava authentication URL with the user's credentials.
func (s *StravaBot) toAuth() {

	missing := []string{}
	requiredEnvVars := []string{"USER_EMAIL", "USER_PASSWORD", "CLIENT_SECRET"}
	for _, envVar := range requiredEnvVars {
		_, exists := os.LookupEnv(envVar)
		if !exists {
			missing = append(missing, envVar)
		}
	}

	if len(missing) > 0 {
		slog.Error("required environment variables not found", "missing", missing)
		os.Exit(3)
	}

	var headers = map[string]string{}

	headers["Content-Type"] = "application/json"

	email := os.Getenv("USER_EMAIL")
	password := os.Getenv("USER_PASSWORD")
	clientSecret := os.Getenv("CLIENT_SECRET")

	type AuthReqBody struct {
		ClientId     int    `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
		Email        string `json:"email,omitempty"`
		Password     string `json:"password,omitempty"`
	}

	authReqBody := &AuthReqBody{
		ClientId:     2,
		ClientSecret: clientSecret,
		Email:        email,
		Password:     password,
	}

	authReqBodyJson, err := json.Marshal(authReqBody)
	if err != nil {
		slog.Error("error marshalling auth request body", "error", slog.String("error", err.Error()))
		return
	}

	html, statusCode := s.Client.MakeRequest(s.MapUrls["auth_url"], "POST", string(authReqBodyJson), headers)

	if statusCode != 200 {
		slog.Error("error getting auth response", "statusCode", statusCode, "body", html)
		os.Exit(5)
	}

	var result map[string]interface{}

	err = json.Unmarshal([]byte(html), &result)
	if err != nil {
		slog.Error("error unmarshalling auth response", "error", slog.String("error", err.Error()))
		return
	}

	if accessToken, ok := result["access_token"].(string); ok {
		s.authToken = accessToken
		s.Client.SetAuthToken(accessToken)
		err = s.saveAuthToken()
		if err != nil {
			slog.Error("error saving auth token", "error", err)
		}
	} else {
		slog.Error("access_token not found or not a string in auth response")
	}
}

func (s *StravaBot) ReadAuthToken() {

	tokenFile, tokenEnvExists := os.LookupEnv("AUTH_TOKEN")

	if tokenEnvExists {

		token, err := os.ReadFile(tokenFile)

		if err != nil {
			slog.Error("couldn't read token file", "error", err)
		}

		s.authToken = string(token)
	}
}

func (s *StravaBot) saveAuthToken() (err error) {
	err = os.WriteFile(os.Getenv("AUTH_TOKEN"), []byte(s.authToken), 0644)
	return
}
