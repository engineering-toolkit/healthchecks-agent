package healthcheck

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	sysCfg "github.com/engineering-toolkit/healthchecks-agent/config"
	"github.com/engineering-toolkit/healthchecks-agent/shared/utils"

	"github.com/go-resty/resty/v2"
)

// Config represent healthcheck client configuration
type Config struct {
	ServerURL        string
	RetryCount       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
}

// Client represent healthcheck client
type Client struct {
	httpClient *resty.Client
	config     Config
}

// Log represent log structure
type Log struct {
	LogTypes  map[string]string
	LogValues interface{}
}

// LogData represent Log Data values
type LogData struct {
	RequestCommand string
	ResponseStatus string
	ErrorMessage   string
}

// New healthcheck client
func New(Config Config) Client {

	c := Client{}
	c.config = Config

	c.httpClient = resty.New()
	c.httpClient.SetRetryCount(c.config.RetryCount).
		SetRetryWaitTime(c.config.RetryWaitTime).
		SetRetryMaxWaitTime(c.config.RetryMaxWaitTime).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		}).
		AddRetryCondition(
			// RetryConditionFunc type is for retry condition function
			// input: non-nil Response OR request execution error
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == http.StatusTooManyRequests
			},
		)
	c.httpClient.SetHeader("User-Agent", "healthcheck-agent/"+sysCfg.Version+";go-resty/"+resty.Version)

	return c
}

// Success is called when your service check is healty
func (c *Client) Success(checkUUID string) error {
	r, e := c.httpClient.R().Get(c._createURL(checkUUID))

	if r.StatusCode() == 200 {
		return nil
	} else {

		switch r.StatusCode() {
		case 500:
			return fmt.Errorf("Error: %d", r.StatusCode())
		case 404:
			return fmt.Errorf("Error: %d; checkUUID: %s (%s)", r.StatusCode(), checkUUID, r.Status())
		default:
		}
	}

	return e
}

// Start is called when your service check is started
func (c *Client) Start(checkUUID string) error {
	r, e := c.httpClient.R().Get(fmt.Sprintf("%s/start", c._createURL(checkUUID)))

	if r.StatusCode() == 200 {
		return nil
	} else {

		switch r.StatusCode() {
		case 500:
			return fmt.Errorf("Error: %d", r.StatusCode())
		case 404:
			return fmt.Errorf("Error: %d; checkUUID: %s (%s)", r.StatusCode(), checkUUID, r.Status())
		default:
		}
	}

	return e
}

// Fail is called when your service check is failed
func (c *Client) Fail(checkUUID string) error {
	r, e := c.httpClient.R().Get(fmt.Sprintf("%s/fail", c._createURL(checkUUID)))

	if r.StatusCode() == 200 {
		return nil
	} else {

		switch r.StatusCode() {
		case 500:
			return fmt.Errorf("Error: %d", r.StatusCode())
		case 404:
			return fmt.Errorf("Error: %d; checkUUID: %s (%s)", r.StatusCode(), checkUUID, r.Status())
		default:
		}
	}

	return e
}

// Log is called when your service want to send log information
func (c *Client) Log(checkUUID string, log Log) error {
	for k, e := range log.LogTypes {
		// SetHeader("Content-Type", "application/json")
		c.httpClient.SetHeader(k, e)
	}

	r, e := c.httpClient.R().
		SetBody(log.LogValues).
		Post(fmt.Sprintf("%s/log", c._createURL(checkUUID)))

	if r.StatusCode() == 200 {
		return nil
	} else {

		switch r.StatusCode() {
		case 500:
			return fmt.Errorf("Error: %d", r.StatusCode())
		case 404:
			return fmt.Errorf("Error: %d; checkUUID: %s (%s)", r.StatusCode(), checkUUID, r.Status())
		default:
		}
	}

	return e
}

// _createURL create pring URL
func (c *Client) _createURL(checkUUID string) string {
	if strings.Contains(c.config.ServerURL, "https://hc-ping.com") {
		return utils.CreateURL(c.config.ServerURL, fmt.Sprintf("%s", checkUUID))
	}
	return utils.CreateURL(c.config.ServerURL, fmt.Sprintf("ping/%s", checkUUID))
}
