package parser

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	userAgent    string
	WebClient    *http.Client
	requestBody  io.Reader
	request      *http.Request
	response     *http.Response
	responseData string
	err          error
	url          *url.URL
}

func NewClient() *Client {
	c := &Client{}
	c.WebClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableCompression:    true,
		},
	}

	c.userAgent = "Strava/33.0.0 (Linux; Android 8.0.0; Pixel 2 XL Build/OPD1.170816.004)"

	return c
}

func (c *Client) MakeRequest(
	siteUrl string,
	method string,
	body string,
	headers map[string]string,
) (string, int) {

	var err error
	c.url, err = url.Parse(siteUrl)
	if err != nil {
		slog.Error("couldn't parse url", "url", siteUrl)
	}

	if c.requestBody = nil; len(body) > 0 {
		c.requestBody = strings.NewReader(body)
	}

	c.request, err = http.NewRequest(method, c.url.String(), c.requestBody)
	if err != nil {
		slog.Error("couldn't create request", "err", slog.String("error", err.Error()))
	}

	c.request.Header.Set("user-agent", c.userAgent)

	if len(headers) > 0 {
		for name, values := range headers {
			c.request.Header.Set(name, values)
		}
	}

	c.response, err = c.WebClient.Do(c.request)
	if err != nil {
		slog.Error("couldn't make request", "err", slog.String("error", err.Error()))
	}

	if c.response != nil {
		if c.response.StatusCode == http.StatusOK {
			page, err := io.ReadAll(c.response.Body)
			if err != nil {
				slog.Error("couldn't read response body", "err", slog.String("error", err.Error()))
			}
			c.responseData = string(page)
		}

		defer func() {
			if err := c.response.Body.Close(); err != nil {
				slog.Error("couldn't close response body", "err", slog.String("error", err.Error()))
			}
		}()

	}

	return c.responseData, c.response.StatusCode
}
