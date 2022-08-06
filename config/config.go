package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name    string  `yaml:"name"`
	Binance Binance `yaml:"binance"`
	Notify  Notify  `yaml:"notify"`
	Trades  []Trade `yaml:"trades"`
}

type Binance struct {
	ApiKey    string `yaml:"api_key"`
	ApiSecret string `yaml:"api_secret"`
}

type Notify struct {
	WebhookURL string `yaml:"webhook_url"`
}

type Trade struct {
	Symbol      string        `yaml:"symbol"`
	Interval    time.Duration `yaml:"interval"`
	Amount      float64       `yaml:"amount"`
	StopLoss    float64       `yaml:"stop_loss"`
	RSIBuy      float64       `yaml:"rsi_buy"`
	RSISell     float64       `yaml:"rsi_sell"`
	RSILimit    int           `yaml:"rsi_limit"`
	RSIInterval string        `yaml:"rsi_interval"`
}

func (s *Trade) GetAsset() string {
	return fmt.Sprintf("%s%s", s.GetSymbol(), s.BuyWith())
}

// GetSymbol split to get symbol that want to buy.
func (s *Trade) GetSymbol() string {
	return strings.Split(s.Symbol, "/")[0]
}

// BuyWith split to get symbol used to buy.
func (s *Trade) BuyWith() string {
	return strings.Split(s.Symbol, "/")[1]
}

// Load configuration file.
func Load(configFile string) (*Config, error) {
	cfg, err := readConfigurationFile(configFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// readConfigurationFile read configuration file.
func readConfigurationFile(fileName string) (config *Config, err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(fileName); err == nil {
		return parseAndValidateConfig(bytes)
	}
	return
}

// parseAndValidateConfig parse configuration file and validates.
func parseAndValidateConfig(yamlBytes []byte) (config *Config, err error) {
	yamlBytes = []byte(os.ExpandEnv(string(yamlBytes)))
	if err = yaml.Unmarshal(yamlBytes, &config); err != nil {
		return
	}

	return
}
