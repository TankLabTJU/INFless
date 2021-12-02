package types

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	GatewayURL         string
	PrometheusHost     string
	InactivityDuration time.Duration
	ReconcileInterval  time.Duration
	PrometheusPort     int
}

//ReadConfig reads configuration files
func ReadConfig() (Config, error) {
	config := Config{}

	config.GatewayURL = "http://192.168.0.111:31112/"
	//config.GatewayURL = os.Getenv("gateway_url")
	if len(config.GatewayURL) == 0 {
		return config, fmt.Errorf("env-var gateway_url must be set\n")
	}

	config.PrometheusHost = "192.168.0.111"
	//config.PrometheusHost = os.Getenv("prometheus_host")
	if len(config.PrometheusHost) == 0 {
		return config, fmt.Errorf("env-var prometheus_host must be set\n")
	}

	config.InactivityDuration = time.Minute * 5
	if val, exists := os.LookupEnv("inactivity_duration"); exists {
		parsedVal, parseErr := time.ParseDuration(val)
		if parseErr != nil {
			return config, parseErr
		}
		config.InactivityDuration = parsedVal
	}

	//config.PrometheusPort = 9090
	config.PrometheusPort = 31113
	if val, exists := os.LookupEnv("prometheus_port"); exists {
		port, parseErr := strconv.Atoi(val)
		if parseErr != nil {
			return config, parseErr
		}
		config.PrometheusPort = port
	}

	config.ReconcileInterval = time.Second * 30
	if val, exists := os.LookupEnv("reconcile_interval"); exists {
		parsedVal, parseErr := time.ParseDuration(val)
		if parseErr != nil {
			return config, parseErr
		}
		config.ReconcileInterval = parsedVal
	}
	return config, nil
}
