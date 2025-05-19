package parser

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Client struct {
	userAgent string
	WebClient *http.Client
	authToken string
	headerMu  sync.Mutex
}

// NewClient creates and initializes a new Client with a default HTTP client and user agent.
func NewClient(authToken ...string) *Client {
	c := &Client{}
	c.WebClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second, // TCP keep-alive
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			// DisableKeepAlives:     true, // Add this line to disable HTTP keep-alives
		},
	}

	c.userAgent = "Strava/33.0.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004)"

	if len(authToken) > 0 && authToken[0] != "" {
		c.authToken = authToken[0]
	}
	c.headerMu = sync.Mutex{}

	return c
}

// SetAuthToken sets the authentication token for the client.
func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

// MakeRequest performs an HTTP request to the specified URL with the given method, body, and headers.
// It returns the response body as a string and the HTTP status code.
func (c *Client) MakeRequest(
	siteUrl string,
	method string,
	body string,
	headers map[string]string,
) (string, int) {

	parsedURL, parseErr := url.Parse(siteUrl)
	if parseErr != nil {
		slog.Error("couldn't parse url", "url", siteUrl, "err", slog.String("error", parseErr.Error()))
		return "", 0
	}

	var localRequestBody io.Reader
	if len(body) > 0 {
		localRequestBody = strings.NewReader(body)
	}

	request, reqErr := http.NewRequest(method, parsedURL.String(), localRequestBody)
	if reqErr != nil {
		slog.Error("couldn't create request", "url", siteUrl, "method", method, "err", slog.String("error", reqErr.Error()))
		return "", 0
	}

	c.headerMu.Lock()
	request.Header.Set("user-agent", c.userAgent)
	if c.authToken != "" && !strings.Contains(siteUrl, "/oauth/internal/token") {
		request.Header.Set("Authorization", "access_token "+c.authToken)
	}
	if len(headers) > 0 {
		for name, values := range headers {
			// Ensure Authorization header isn't overwritten if already set by authToken,
			// unless it's explicitly passed in headers for a different purpose (e.g. basic auth for token endpoint)
			if strings.ToLower(name) == "authorization" && c.authToken != "" && !strings.Contains(siteUrl, "/oauth/internal/token") {
				// Potentially skip or log if trying to overwrite Strava token auth
				// For now, allow override if explicitly provided
			}
			request.Header.Set(name, values)
		}
	}
	c.headerMu.Unlock()

	response, doErr := c.WebClient.Do(request)
	if doErr != nil {
		slog.Error("couldn't make request", "url", siteUrl, "method", method, "err", slog.String("error", doErr.Error()))
		// Try to get status code even if there's an error
		statusCode := 0
		if response != nil {
			statusCode = response.StatusCode
			if response.Body != nil {
				response.Body.Close() // Ensure body is closed on error path
			}
		}
		return "", statusCode
	}

	// Defensive check, though Do() should return non-nil response if err is nil.
	if response == nil {
		slog.Error("received nil response from http client (and nil error)", "url", siteUrl, "method", method)
		return "", 0 // Or a more specific error status code
	}

	// Handle cases where response.Body might be nil (e.g., 204 No Content)
	if response.Body == nil {
		slog.Warn("response body is nil", "url", siteUrl, "method", method, "statusCode", response.StatusCode)
		// Return empty string; caller will likely encounter JSON unmarshal error if it expects a body.
		// This is often acceptable for non-2xx or 204 responses.
		return "", response.StatusCode
	}

	// Ensure the response body is closed when we're done.
	defer response.Body.Close()

	page, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		slog.Error("couldn't read response body", "url", siteUrl, "method", method, "statusCode", response.StatusCode, "err", slog.String("error", readErr.Error()))
		// response.Body will be closed by the defer.
		return "", response.StatusCode
	}

	responseData := string(page)
	return responseData, response.StatusCode
}
