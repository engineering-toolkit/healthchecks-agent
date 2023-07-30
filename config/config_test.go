package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	ch, err := NewConfigHandler(true)
	if err != nil {
		t.Errorf("NewConfigHandler: %s", err.Error())
		return
	}

	c, err := ch.Read()
	if err != nil {
		t.Errorf("CH.Read: %s", err.Error())
		return
	}

	if c.HC.ServerURL != "http://localhost:8000" {
		t.Error("Invalid HC.ServerULR")
		return
	}

	assert.Equal(t, c.SGClientChecks, ch.ClientChecks, "should be same")

}
