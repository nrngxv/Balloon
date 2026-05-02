package logger

import (
	"os"
	"io"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nrngxv/go-boilerplate/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

//Creating a connection directly to new relic for observavility
type LoggerService struct {
	nrApp *newrelic.Application
}


func NewLoggerService(nls *config.ObservabilityConfig) *LoggerService {
	service := &LoggerService{}

	if nls.NewRelic.LicenseKey == "" {
		return service
	}

	var configOptions []newrelic.ConfigOption

	// A bunch of settings that are sent to new relic on connection
	configOptions = append(configOptions,
		newrelic.ConfigAppName(nls.ServiceName),
		newrelic.ConfigLicense(nls.NewRelic.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(nls.NewRelic.AppLogForwardingEnabled),
		newrelic.ConfigDistributedTracerEnabled(nls.NewRelic.DistributedTracingEnabled),
	)

	//Adding debug logging
	if nls.NewRelic.DebugLogging {
		configOptions = append(configOptions, newrelic.ConfigDebugLogger(os.Stdout))
	}

	app, err := newrelic.NewApplication(configOptions...) // this ... is a spread operator
	if err != nil {
		return service
	}

	service.nrApp = app
	return service
}

// THis is take the logger service implicitly and return the nrapp
func (ls *LoggerService) GetApplication() *newrelic.Application {
	return ls.nrApp
}

func NewLogger(level string, isProd bool) zerolog.Logger {
	return NewLoggerWithService(&config.ObservabilityConfig{
		Logging: config.LoggingConfig{Level: level}, Environment: func() string {
			if isProd {
				return "production"
			}
			return "development"
		}(),
	}, nil)
}


func NewLoggerWithService(cfg *config.ObservabilityConfig, loggerservice *LoggerService) zerolog.Logger {
	var logLevel zerolog.Level
	level := cfg.GetLogLevel()

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
		
	}

	zerolog.TimeFieldFormat = "2000-01-02 12:09:04"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var writer io.Writer

	//This is the base writer
	var baseWriter io.Writer
	if cfg.isProduction() && cfg.Logging.Format == "json"{

		baseWriter = os.Stdout

		if loggerservice != nil && loggerservice.nrApp != nil {
			nrWriter := zerologWriter.New(baseWriter, loggerService.nrApp)
			writer = nrWriter
			} else {
				writer = basebaseWriter
			}

			} else {
				consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2000-01-02 12:09:04"}
				writer = consoleWriter
		}
	}

	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("environment", cfg.environment).
		Logger()

		if !cfg.IsProduction() {
		logger = logger.With().Stack().Logger()
	}

	return logger
}


// This is to get trace contents with zerolog and then putting it in New relic
func WithTraceContext(logger zerolog.Logger, txn *newrelic.Transaction) zerolog.Logger {
	if txn == nil {
		return logger
	}

	// Get trace metadata from transaction
	metadata := txn.GetTraceMetadata()

	return logger.With().
		Str("trace.id", metadata.TraceID).
		Str("span.id", metadata.SpanID).
		Logger()
	}


// To get database logging with pgx driver
func NewPgxLogger(level zerolog.Level) zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		FormatFieldValue: func(i any) string {
			switch v := i.(type) {
			case string:
				// Clean and format SQL for better readability
				if len(v) > 200 {
					// Truncate very long SQL statements
					return v[:200] + "..."
				}
				return v
			case []byte:
				var obj interface{}
				if err := json.Unmarshal(v, &obj); err == nil {
					pretty, _ := json.MarshalIndent(obj, "", "    ")
					return "\n" + string(pretty)
				}
				return string(v)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}

	return zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Str("component", "database").
		Logger()
}

// GetPgxTraceLogLevel converts zerolog level to pgx tracelog level
func GetPgxTraceLogLevel(level zerolog.Level) int {
	switch level {
	case zerolog.DebugLevel:
		return 6 // tracelog.LogLevelDebug
	case zerolog.InfoLevel:
		return 4 // tracelog.LogLevelInfo
	case zerolog.WarnLevel:
		return 3 // tracelog.LogLevelWarn
	case zerolog.ErrorLevel:
		return 2 // tracelog.LogLevelError
	default:
		return 0 // tracelog.LogLevelNone
	}
}