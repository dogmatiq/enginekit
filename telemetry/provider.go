package telemetry

import (
	"runtime/debug"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	nooplog "go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	noopmetric "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

// Provider provides [Recorder] instances scoped to particular subsystems.
//
// The zero value of a *Provider is equivalent to a provider configured with
// no-op tracer, meter and logging providers.
type Provider struct {
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	LoggerProvider log.LoggerProvider
}

// Recorder records traces, metrics and logs for a particular subsystem.
type Recorder struct {
	tracer  trace.Tracer
	meter   metric.Meter
	logger  log.Logger
	attrKVs attribute.Set
	logKVs  []log.KeyValue

	errorCount              Instrument[int64]
	operationCount          Instrument[int64]
	operationsInFlightCount Instrument[int64]
}

// Recorder returns a new Recorder instance.
//
// pkg is the path to the Go package that is performing the instrumentation. If
// it is an internal package, use the package path of the public parent package
// instead.
func (p *Provider) Recorder(pkg string, attrs ...Attr) *Recorder {
	var (
		tracerProvider trace.TracerProvider
		meterProvider  metric.MeterProvider
		loggerProvider log.LoggerProvider
	)

	if p != nil {
		tracerProvider = p.TracerProvider
		meterProvider = p.MeterProvider
		loggerProvider = p.LoggerProvider
	}

	if tracerProvider == nil {
		tracerProvider = nooptrace.NewTracerProvider()
	}

	if meterProvider == nil {
		meterProvider = noopmetric.NewMeterProvider()
	}

	if loggerProvider == nil {
		loggerProvider = nooplog.NewLoggerProvider()
	}

	version := moduleVersion(pkg)

	r := &Recorder{
		tracer:  tracerProvider.Tracer(pkg, trace.WithInstrumentationVersion(version)),
		meter:   meterProvider.Meter(pkg, metric.WithInstrumentationVersion(version)),
		logger:  loggerProvider.Logger(pkg, log.WithInstrumentationVersion(version)),
		attrKVs: attribute.NewSet(asAttrKeyValues(attrs)...),
		logKVs:  asLogKeyValues(attrs),
	}

	r.errorCount = r.Counter("errors", "{error}", "The number of errors that have occurred.")
	r.operationCount = r.Counter("operations", "{operation}", "The number of operations that have been performed.")
	r.operationsInFlightCount = r.UpDownCounter("operations.in_flight", "{operation}", "The number of operations that are currently in progress.")

	return r
}

func moduleVersion(pkg string) string {
	module := ""
	version := "unknown"

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == pkg {
				return dep.Version
			}

			if strings.HasPrefix(pkg, dep.Path+"/") {
				if len(dep.Path) > len(module) {
					module = dep.Path
					version = dep.Version
				}
			}
		}
	}

	return version
}
