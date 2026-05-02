package config

import (
	"fmt"
	"time"
)

type ObservabilityConfig struct{
	ServiceName		string				`koanf:"service_name" validate:"required"`
	Environment		string				`koanf:"environment" validate:"required"`
	Logging			LoggingConfig		`koanf:"logging" validate:"required"`
	NewRelic		NewRelicConfig		`koanf:"new_relic" validate:"required"`
	HealthChecks	HealthChecksConfig	`koanf:"health_checks" validate:"required"`
}

type LoggingConfig struct{
	Level				string			`koanf:"level" validate:"required"`
	Format				string			`koanf:"format" validate:"required"`
	SlowQueryThreshold 	time.Duration	`koanf:"slow_query_threshold"`
}

type NewRelicConfig struct{
	LicenseKey                string 	`koanf:"license_key" validate:"required"`
	AppLogForwardingEnabled   bool   	`koanf:"app_log_forwarding_enabled"`
	DistributedTracingEnabled bool   	`koanf:"distributed_tracing_enabled"`
	DebugLogging              bool   	`koanf:"debug_logging"`
}

type HealthChecksConfig struct{
	Enabled		bool			`koanf:"enabled"`
	Interval	time.Duration	`koanf:"interval" validate:"min=1s"`
	Timeout		time.Duration	`koanf:"enabled" validate:"min=1s"`
	Checks		[]string		`koanf:"checks"`
}


// These are the default values for the observability config. In case no values are given in .env
func DefaultObservabilityConfig() *ObservabilityConfig{
	return &ObservabilityConfig{
		ServiceName: "balloon",
		Environment: "local",
		Logging: LoggingConfig{
			Level:  			"debug",
			Format: 			"json",
			SlowQueryThreshold: 100* time.Millisecond,
		},

		NewRelic: NewRelicConfig{
			LicenseKey: "",
			AppLogForwardingEnabled: 	true,
			DistributedTracingEnabled: 	true,
			DebugLogging:				false,

		},

		HealthChecks: HealthChecksConfig{
			Enabled: true,
			Interval: 30 * time.Second,
			Timeout: 5 * time.Second,
			Checks: []string{"database", "redis"},
		},
	}
}

//to validate all the observability configs
func (o *ObservabilityConfig) Validate() error {
	if o.ServiceName == ""{
		return fmt.Errorf("The service name is missing")
	}

	// these are valid log levels
	ValidLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}

	if !ValidLogLevels[o.Logging.Level]{
		return fmt.Errorf("Logging level invalid: %s (Must be one of: debug, info, warn, error)", o.Logging.Level)
	}

	if o.Logging.SlowQueryThreshold < 0 {
		return fmt.Errorf("The slow_query_threshold must be non-negative")
	}

	return nil
}

func (o *ObservabilityConfig) GetLogLevel() string{
	switch o.Environment {

		case "production":
			if o.Logging.Level == "" {
				return "info"
			}
		case "development":
			if o.Logging.Level == "" {
				return "debug"
			}
	}

	return o.Logging.Level
}

func (o *ObservabilityConfig) isProduction() bool {
	return o.Environment == "Production"
}