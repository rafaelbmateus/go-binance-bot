package config_test

import (
	"testing"

	"github.com/rafaelbmateus/binance-bot/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := config.Config{Name: "binance-bot"}
	assert.Equal(t, "binance-bot", cfg.Name)
}

func TestLoad(t *testing.T) {
	t.Run("config ok", func(t *testing.T) {
		_, err := config.Load("../config.yaml")
		assert.NoError(t, err)
	})

	t.Run("config not found", func(t *testing.T) {
		_, err := config.Load("config.yml")
		assert.Error(t, err)
	})
}
