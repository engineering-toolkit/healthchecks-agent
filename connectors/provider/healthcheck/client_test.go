package healthcheck

import (
	"testing"
	"time"
)

func newClient(t *testing.T) (Client, error) {
	c := New(Config{
		ServerURL:        "https://hc-ping.com",
		RetryCount:       3,
		RetryWaitTime:    5 * time.Second,
		RetryMaxWaitTime: 20 * time.Second})

	return c, nil
}

func TestStart(t *testing.T) {
	c, err := newClient(t)
	if err != nil {
		t.Errorf("newClient: %s", err.Error())
		return
	}

	err = c.Start("50998592-b407-4a6e-a853-c4520e2d5766")
	if err != nil {
		t.Errorf("Start: %s", err.Error())
		return
	}
}

func TestFail(t *testing.T) {
	c, err := newClient(t)
	if err != nil {
		t.Errorf("newClient: %s", err.Error())
		return
	}

	err = c.Fail("50998592-b407-4a6e-a853-c4520e2d5766")
	if err != nil {
		t.Errorf("Start: %s", err.Error())
		return
	}
}

func TestSuccess(t *testing.T) {
	c, err := newClient(t)
	if err != nil {
		t.Errorf("newClient: %s", err.Error())
		return
	}

	err = c.Success("50998592-b407-4a6e-a853-c4520e2d5766")
	if err != nil {
		t.Errorf("Start: %s", err.Error())
		return
	}
}

func TestLog(t *testing.T) {
	c, err := newClient(t)
	if err != nil {
		t.Errorf("newClient: %s", err.Error())
		return
	}

	err = c.Log("50998592-b407-4a6e-a853-c4520e2d5766",
		Log{LogTypes: map[string]string{"Content-Type": "application/json"},
			LogValues: `{"sample_ata":"sample-data-value"}`})
	if err != nil {
		t.Errorf("Start: %s", err.Error())
		return
	}
}
