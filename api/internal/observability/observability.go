package observability

import (
	"api/config"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/grafana/loki-client-go/loki"
	"github.com/grafana/loki-client-go/pkg/urlutil"
	promConfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	Logger        *slog.Logger
	meterProvider *metric.MeterProvider
	lokiClient    *loki.Client
)

// LokiWriter implements io.Writer to send logs to Loki
type LokiWriter struct {
	client      *loki.Client
	serviceName string
	environment string
}

func (w *LokiWriter) Write(p []byte) (n int, err error) {
	// Parse log level from JSON if possible (basic parsing)
	logLevel := "info"
	logStr := string(p)
	if strings.Contains(logStr, `"level":"error"`) || strings.Contains(logStr, `"level":"ERROR"`) {
		logLevel = "error"
	} else if strings.Contains(logStr, `"level":"warn"`) || strings.Contains(logStr, `"level":"WARN"`) {
		logLevel = "warn"
	} else if strings.Contains(logStr, `"level":"debug"`) || strings.Contains(logStr, `"level":"DEBUG"`) {
		logLevel = "debug"
	}

	labels := model.LabelSet{
		"service":     model.LabelValue(w.serviceName),
		"environment": model.LabelValue(w.environment),
		"level":       model.LabelValue(logLevel),
	}

	if err := w.client.Handle(labels, time.Now(), logStr); err != nil {
		// Log to stderr if Loki fails, but don't block
		fmt.Fprintf(os.Stderr, "Failed to send to Loki: %v\n", err)
	}

	return len(p), nil
}

// Initialize sets up observability with Loki logs and Prometheus/OTLP metrics
func Initialize(cfg config.Config) error {
	// Initialize structured logging with Loki
	if err := initLogger(cfg); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize metrics
	if err := initMetrics(cfg); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down observability components
func Shutdown(ctx context.Context) error {
	var errs []error

	if lokiClient != nil {
		Logger.Info("Stopping Loki client")
		lokiClient.Stop()
	}

	if meterProvider != nil {
		Logger.Info("Shutting down meter provider")
		if err := meterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}

func initLogger(cfg config.Config) error {
	var writer io.Writer = os.Stdout

	// Si tenim configuració de Loki, enviem logs allà també
	if cfg.Observability.LokiURL != "" && 
	   cfg.Observability.LokiUser != "" && 
	   cfg.Observability.LokiPassword != "" {
		
		lokiCfg := loki.Config{}
		
		u, err := url.Parse(cfg.Observability.LokiURL)
		if err != nil {
			return fmt.Errorf("failed to parse loki url: %w", err)
		}
		lokiCfg.URL = urlutil.URLValue{URL: u}
		
		lokiCfg.TenantID = ""
		lokiCfg.BatchWait = 1 * time.Second
		lokiCfg.BatchSize = 1024 * 1024 // 1MB
		lokiCfg.Timeout = 10 * time.Second
		lokiCfg.BackoffConfig.MinBackoff = 500 * time.Millisecond
		lokiCfg.BackoffConfig.MaxBackoff = 5 * time.Minute
		lokiCfg.BackoffConfig.MaxRetries = 10
		
		lokiCfg.Client.BasicAuth = &promConfig.BasicAuth{
			Username: cfg.Observability.LokiUser,
			Password: promConfig.Secret(cfg.Observability.LokiPassword),
		}

		client, err := loki.New(lokiCfg)
		if err != nil {
			return fmt.Errorf("failed to create loki client: %w", err)
		}

		lokiClient = client

		// MultiWriter per enviar logs tant a stdout com a Loki
		lokiWriter := &LokiWriter{
			client:      client,
			serviceName: cfg.Observability.ServiceName,
			environment: cfg.Observability.Environment,
		}
		writer = io.MultiWriter(os.Stdout, lokiWriter)
		
		// Log temporal a stdout que Loki està configurat
		fmt.Printf("Loki configured: %s\n", cfg.Observability.LokiURL)
	} else {
		fmt.Println("Loki not configured, logging to stdout only")
	}

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Logger = slog.New(handler).With(
		slog.String("service", cfg.Observability.ServiceName),
		slog.String("environment", cfg.Observability.Environment),
	)

	slog.SetDefault(Logger)
	Logger.Info("Logger initialized",
		slog.Bool("loki_enabled", cfg.Observability.LokiURL != ""),
		slog.String("service", cfg.Observability.ServiceName),
		slog.String("environment", cfg.Observability.Environment),
	)

	return nil
}

func initMetrics(cfg config.Config) error {
	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.Observability.ServiceName),
			semconv.DeploymentEnvironment(cfg.Observability.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	var readers []metric.Reader

	// Always create local Prometheus exporter for /metrics endpoint
	promExporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}
	readers = append(readers, promExporter)

	// If Prometheus remote write is configured, add OTLP exporter
	if cfg.Observability.PrometheusURL != "" &&
		cfg.Observability.PrometheusUser != "" &&
		cfg.Observability.PrometheusPassword != "" {

		// Extract endpoint without /api/prom/push if present
		endpoint := strings.TrimPrefix(cfg.Observability.PrometheusURL, "https://")
		endpoint = strings.TrimPrefix(endpoint, "http://")
		endpoint = strings.Split(endpoint, "/")[0]

		otlpExporter, err := otlpmetrichttp.New(
			context.Background(),
			otlpmetrichttp.WithEndpoint(endpoint),
			otlpmetrichttp.WithURLPath("/api/prom/push"),
			otlpmetrichttp.WithHeaders(map[string]string{
				"Authorization": fmt.Sprintf("Basic %s", 
					encodeBasicAuth(cfg.Observability.PrometheusUser, cfg.Observability.PrometheusPassword)),
			}),
		)
		if err != nil {
			Logger.Warn("Failed to create OTLP exporter, metrics will only be available locally",
				slog.Any("error", err))
		} else {
			// Periodic reader to push metrics every 15 seconds
			periodicReader := metric.NewPeriodicReader(
				otlpExporter,
				metric.WithInterval(15*time.Second),
			)
			readers = append(readers, periodicReader)
			Logger.Info("Prometheus remote write configured", 
				slog.String("endpoint", endpoint))
		}
	} else {
		Logger.Info("Prometheus remote write not configured, metrics available only at /metrics endpoint")
	}

	// Create meter provider with all readers
	meterProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(readers[0]),
	)
	
	// Add additional readers if any
	for i := 1; i < len(readers); i++ {
		meterProvider = metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(readers[i]),
		)
	}

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	Logger.Info("Metrics initialized",
		slog.String("service", cfg.Observability.ServiceName),
		slog.String("environment", cfg.Observability.Environment),
		slog.Bool("remote_write_enabled", cfg.Observability.PrometheusURL != ""),
	)

	return nil
}

// encodeBasicAuth encodes username and password for Basic Auth
func encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64Encode(auth)
}

func base64Encode(s string) string {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder
	
	for i := 0; i < len(s); i += 3 {
		chunk := []byte{0, 0, 0}
		chunkLen := 0
		
		for j := 0; j < 3 && i+j < len(s); j++ {
			chunk[j] = s[i+j]
			chunkLen++
		}
		
		result.WriteByte(base64Table[chunk[0]>>2])
		result.WriteByte(base64Table[((chunk[0]&0x03)<<4)|(chunk[1]>>4)])
		
		if chunkLen > 1 {
			result.WriteByte(base64Table[((chunk[1]&0x0F)<<2)|(chunk[2]>>6)])
		} else {
			result.WriteByte('=')
		}
		
		if chunkLen > 2 {
			result.WriteByte(base64Table[chunk[2]&0x3F])
		} else {
			result.WriteByte('=')
		}
	}
	
	return result.String()
}