package types

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func Test_properReadConfig(t *testing.T) {
	defaultVals := []struct {
		Case               string
		title              string
		gatewayURL         string
		prometheusHost     string
		prometheusPort     int
		inactivityDuration time.Duration
		reconcileInterval  time.Duration
	}{
		{
			Case:               "default values",
			title:              "These are the current default values",
			gatewayURL:         "http://gateway:8080/", //Dont have defaults but needed
			prometheusHost:     "prometheus",           //Dont have defaults but needed
			prometheusPort:     9090,
			inactivityDuration: time.Duration(5) * time.Minute,
			reconcileInterval:  time.Duration(30) * time.Second,
		},
		{
			Case:               "manual values",
			title:              "These are values that someone would manually set",
			gatewayURL:         "http://som.random.domain/",
			prometheusHost:     "name",
			prometheusPort:     1234,
			inactivityDuration: time.Duration(1) * time.Minute,  //i.e. "1m"
			reconcileInterval:  time.Duration(45) * time.Second, //i.e. "45s"
		},
	}

	for _, test := range defaultVals {
		if test.Case == "default values" {
			//We need to set those two cause they dont have default values
			os.Setenv("gateway_url", test.gatewayURL)
			os.Setenv("prometheus_host", test.prometheusHost)
			config, _ := ReadConfig()
			if (test.prometheusPort) != (config.PrometheusPort) {
				t.Errorf("Default for prometheus port should be: %d got: %d.", test.prometheusPort, config.PrometheusPort)
			}
			if (test.inactivityDuration) != (config.InactivityDuration) {
				t.Errorf("Default time for inactivity duration should be: %s got :%s", test.inactivityDuration, config.InactivityDuration)
			}
			if (test.reconcileInterval) != (config.ReconcileInterval) {
				t.Errorf("Default time for reconcile interval should be: %s got :%s", test.reconcileInterval, config.ReconcileInterval)
			}
		}
		if test.Case == "manual values" {
			os.Setenv("gateway_url", test.gatewayURL)
			os.Setenv("prometheus_host", test.prometheusHost)
			os.Setenv("prometheus_port", strconv.Itoa(test.prometheusPort))
			os.Setenv("inactivity_duration", test.inactivityDuration.String())
			os.Setenv("reconcile_interval", test.reconcileInterval.String())
			config, _ := ReadConfig()
			if test.gatewayURL != config.GatewayURL {
				t.Errorf("Gateway wanted: %s got :%s", test.gatewayURL, config.GatewayURL)
			}
			if test.prometheusHost != config.PrometheusHost {
				t.Errorf("Prometheus host wanted: %s got :%s", test.prometheusHost, config.PrometheusHost)
			}
			if test.prometheusPort != config.PrometheusPort {
				t.Errorf("Prometheus port wanted: %s got :%s", strconv.Itoa(test.prometheusPort), strconv.Itoa((config.PrometheusPort)))
			}
			if test.inactivityDuration != config.InactivityDuration {
				t.Errorf("Inactivity duration wanted: %s got :%s", test.inactivityDuration.String(), config.InactivityDuration.String())
			}
			if test.reconcileInterval != config.ReconcileInterval {
				t.Errorf("Reconcile interval wanted: %s got :%s", test.reconcileInterval.String(), config.ReconcileInterval.String())
			}
		}
	}
}

func Test_unproperReadConfig(t *testing.T) {
	defaultVals := []struct {
		Case               string
		title              string
		gatewayURL         string
		prometheusHost     string
		prometheusPort     string
		inactivityDuration string
		reconcileInterval  string
	}{
		{
			Case:               "first case",
			title:              "These are valid values",
			gatewayURL:         "http://this.is.valid:8080/", //Not default value but needed
			prometheusHost:     "thename",                    //Not default value but needed
			prometheusPort:     "1234",
			inactivityDuration: "1m",
			reconcileInterval:  "1m",
		},
		{
			Case:               "second case",
			title:              "This is bad setup",
			gatewayURL:         "bad gateway name",
			prometheusHost:     "1234",
			prometheusPort:     "ports are good",
			inactivityDuration: "?",
			reconcileInterval:  "just random things",
		}, {
			Case:               "third case",
			title:              "Everything is unset and should fail",
			gatewayURL:         "",
			prometheusHost:     "",
			prometheusPort:     "",
			inactivityDuration: "",
			reconcileInterval:  "",
		},
	}
	//We need to set those two cause they dont have default values
	for _, test := range defaultVals {
		os.Setenv("gateway_url", test.gatewayURL)
		os.Setenv("prometheus_host", test.prometheusHost)
		os.Setenv("prometheus_port", test.prometheusPort)
		os.Setenv("inactivity_duration", test.inactivityDuration)
		os.Setenv("reconcile_interval", test.reconcileInterval)
		_, configErr := ReadConfig()
		if test.Case == "first case" {
			if configErr != nil {
				t.Errorf("Unexpected error :\n%s", configErr.Error())
			}
		}
		if test.Case == "second case" {
			if configErr == nil {
				t.Errorf("Had to have errors due to bad configuration")
			}
		}
		if test.Case == "third case" {
			if configErr == nil {
				t.Errorf("Had to have errors due configuration set to empty string")
			}
		}
	}
}
