package otlpmetricgrpc

import (
	"context"
	"sync"

	"github.com/f5/otel-arrow-adapter/collector/gen/exporter/otlpexporter"
	"github.com/lightstep/otel-launcher-go/lightstep/sdk/metric"
	"github.com/lightstep/otel-launcher-go/lightstep/sdk/metric/data"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configcompression"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/otel"
	globalmetric "go.opentelemetry.io/otel/metric/global"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type Option func(*Config)

type Config struct {
	Batcher  batchprocessor.Config
	Exporter otlpexporter.Config
}

type client struct {
	once     sync.Once
	resource pcommon.Resource

	exporter exporter.Metrics
	batcher  processor.Metrics
	settings exporter.CreateSettings
}

func NewDefaultConfig() Config {
	return Config{
		Exporter: otlpexporter.Config{
			TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
			RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
			QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
			GRPCClientSettings: configgrpc.GRPCClientSettings{
				Headers:         map[string]configopaque.String{},
				Compression:     configcompression.Zstd,
				WriteBufferSize: 512 * 1024,
			},
			Arrow: otlpexporter.ArrowSettings{
				NumStreams: 1,
				Enabled:    true,
			},
		},
	}
}

func NewConfig(opts ...Option) Config {
	cfg := NewDefaultConfig()
	for _, option := range opts {
		option(&cfg)
	}
	return cfg
}

func WithEndpoint(addr string) Option {
	return func(cfg *Config) {
		cfg.Exporter.GRPCClientSettings.Endpoint = addr
	}
}

func WithHeaders(hdrs map[string]string) Option {
	return func(cfg *Config) {
		for key, val := range hdrs {
			cfg.Exporter.GRPCClientSettings.Headers[key] = configopaque.String(val)
		}
	}
}

func WithCompressor(comp string) Option {
	return func(cfg *Config) {
		cfg.Exporter.GRPCClientSettings.Compression = configcompression.CompressionType(comp)
	}
}

func WithInsecure() Option {
	return func(cfg *Config) {
		cfg.Exporter.GRPCClientSettings.TLSSetting.Insecure = true
	}
}

func WithTLSSetting(tlss configtls.TLSClientSetting) Option {
	return func(cfg *Config) {
		cfg.Exporter.GRPCClientSettings.TLSSetting = tlss
	}
}

func NewClient(ctx context.Context, cfg *Config) (metric.PushExporter, error) {
	c := &client{}

	if cfg.Exporter.Arrow.Enabled {
		c.settings.ID = component.NewID("otlp/arrow")
	} else {
		c.settings.ID = component.NewID("otlp/proto")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	c.settings.TelemetrySettings.Logger = logger
	c.settings.TelemetrySettings.TracerProvider = otel.GetTracerProvider()
	c.settings.TelemetrySettings.MeterProvider = globalmetric.MeterProvider() // Note: becomes otel.GetMeterProvider()
	c.settings.TelemetrySettings.MetricsLevel = configtelemetry.LevelNormal

	exp, err := otlpexporter.NewFactory().CreateMetricsExporter(ctx, c.settings, cfg.Exporter)
	if err != nil {
		return nil, err
	}

	bset := processor.CreateSettings{
		ID:                component.NewID("otlp/batch"),
		TelemetrySettings: c.settings.TelemetrySettings,
		BuildInfo:         c.settings.BuildInfo,
	}

	bat, err := batchprocessor.NewFactory().CreateMetricsProcessor(ctx, bset, cfg.Batcher, exp)
	if err != nil {
		return nil, err
	}

	c.exporter = exp
	c.batcher = bat

	err = exp.Start(ctx, c)
	if err != nil {
		return nil, err
	}

	err = bat.Start(ctx, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *client) String() string {
	return "otel-arrow-adapter/otlpexporter"
}

// ExportMetrics implements PushExporter.
func (c *client) ExportMetrics(ctx context.Context, data data.Metrics) error {
	return c.batcher.ConsumeMetrics(ctx, c.d2pd(data))
}

// ShutdownMetrics implements PushExporter.
func (c *client) ShutdownMetrics(ctx context.Context, data data.Metrics) error {
	var err error
	if err1 := c.ForceFlushMetrics(ctx, data); err1 != nil {
		err = multierr.Append(err, err1)
	}
	if err2 := c.batcher.Shutdown(ctx); err2 != nil {
		err = multierr.Append(err, err2)
	}
	if err3 := c.exporter.Shutdown(ctx); err3 != nil {
		err = multierr.Append(err, err3)
	}
	return err
}

// ForceFlushMetrics implements PushExporter.
func (c *client) ForceFlushMetrics(ctx context.Context, data data.Metrics) error {
	return c.ExportMetrics(ctx, data)
}

// ReportFatalError implements component.Host.
func (c *client) ReportFatalError(err error) {
	c.settings.Logger.Fatal("exporter fatal", zap.Error(err))
}

// GetFactory implements component.Host.
func (c *client) GetFactory(component.Kind, component.Type) component.Factory {
	return nil
}

// GetExtensions implements component.Host.
func (c *client) GetExtensions() map[component.ID]component.Component {
	return nil
}

// GetExporters implements component.Host.
func (c *client) GetExporters() map[component.DataType]map[component.ID]component.Component {
	return nil
}
