package middleware

import (
	"api/internal/observability"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	httpRequestDuration metric.Float64Histogram
	httpRequestCounter  metric.Int64Counter
)

func init() {
	meter := otel.Meter("turniq-api")
	
	var err error
	httpRequestDuration, err = meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		panic(err)
	}

	httpRequestCounter, err = meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		panic(err)
	}
}

// ObservabilityMiddleware adds metrics and logging to HTTP requests
func ObservabilityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		// Record metrics
		attrs := []attribute.KeyValue{
			attribute.String("http.method", method),
			attribute.String("http.route", path),
			attribute.Int("http.status_code", status),
		}

		httpRequestDuration.Record(c.Request.Context(), duration, metric.WithAttributes(attrs...))
		httpRequestCounter.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))

		// Log request
		if observability.Logger != nil {
			observability.Logger.Info("HTTP request",
				slog.String("method", method),
				slog.String("path", path),
				slog.Int("status", status),
				slog.Float64("duration_seconds", duration),
				slog.String("client_ip", c.ClientIP()),
			)
		}
	}
}
