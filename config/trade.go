package config

import (
	"strings"
	"time"
)

type Trade struct {
	Symbol    string        `yaml:"symbol"`
	Interval  time.Duration `yaml:"interval"`
	BuyPrice  float64       `yaml:"buyPrice"`
	SellPrice float64       `yaml:"sellPrice"`
	Limit     float64       `yaml:"limit"`
}

// GetSymbol split to get symbol that want to buy.
func (s *Trade) GetSymbol() string {
	return strings.Split(s.Symbol, "/")[0]
}

// BuyWith split to get symbol used to buy.
func (s *Trade) BuyWith() string {
	return strings.Split(s.Symbol, "/")[1]
}
