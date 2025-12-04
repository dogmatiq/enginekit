package telemetry

import (
	"context"
	"log/slog"
	"strconv"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

// NewSLogProvider returns a [log.LoggerProvider] thats writes logs to the
// provided [slog.Logger].
func NewSLogProvider(target *slog.Logger) log.LoggerProvider {
	return &slogProvider{
		Target: target,
	}
}

// slogProvider adapts an [slog.Logger] to the OpenTelemetry
// [log.LoggerProvider] interface.
type slogProvider struct {
	embedded.LoggerProvider

	// Target is the underlying [slog.Logger] that logs are written to.
	Target *slog.Logger
}

// Logger returns a new [Logger] with the provided name and configuration.
func (p *slogProvider) Logger(name string, _ ...log.LoggerOption) log.Logger {
	return &standardLogger{Target: p.Target.WithGroup(name)}
}

type standardLogger struct {
	embedded.Logger
	Target *slog.Logger
}

func (l *standardLogger) Emit(ctx context.Context, rec log.Record) {
	var (
		attrs   []slog.Attr
		message = "?"
	)

	if rec.Body().Kind() == log.KindString {
		message = rec.Body().AsString()
	} else if !rec.Body().Empty() {
		attrs = append(
			attrs,
			slogAttrFromLogValue("body", rec.Body()),
		)
	}

	if ev := rec.EventName(); ev != "" {
		attrs = append(attrs, slog.String("event", ev))
	}

	rec.WalkAttributes(
		func(kv log.KeyValue) bool {
			attrs = append(
				attrs,
				slogAttrFromLogValue(kv.Key, kv.Value),
			)
			return true
		},
	)

	if !rec.Timestamp().IsZero() {
		attrs = append(attrs, slog.Time("timestamp", rec.Timestamp()))
	}

	if !rec.ObservedTimestamp().IsZero() {
		attrs = append(attrs, slog.Time("observed_timestamp", rec.ObservedTimestamp()))
	}

	l.Target.LogAttrs(
		ctx,
		slogLevelFromLogSeverity(rec.Severity()),
		message,
		attrs...,
	)
}

func (l *standardLogger) Enabled(ctx context.Context, p log.EnabledParameters) bool {
	return l.Target.Enabled(ctx, slogLevelFromLogSeverity(p.Severity))
}

// slogLevelFromLogSeverity maps an OpenTelemetry [log.Severity] to a
// [slog.Level].
func slogLevelFromLogSeverity(sev log.Severity) slog.Level {
	level := slog.LevelDebug

	switch {
	case sev >= log.SeverityError:
		level = slog.LevelError
	case sev >= log.SeverityWarn:
		level = slog.LevelWarn
	case sev >= log.SeverityInfo:
		level = slog.LevelInfo
	}

	return level
}

// slogAttrFromLogValue converts an OpenTelemetry [log.Value] to an
// [slog.Attr].
func slogAttrFromLogValue(name string, v log.Value) slog.Attr {
	switch v.Kind() {
	case log.KindEmpty:
		return slog.Any(name, nil)

	case log.KindBool:
		return slog.Bool(name, v.AsBool())

	case log.KindFloat64:
		return slog.Float64(name, v.AsFloat64())

	case log.KindInt64:
		return slog.Int64(name, v.AsInt64())

	case log.KindString:
		return slog.String(name, v.AsString())

	case log.KindBytes:
		return slog.Any(name, v.AsBytes())

	case log.KindSlice:
		var attrs []slog.Attr
		for i, elem := range v.AsSlice() {
			attrs = append(
				attrs,
				slogAttrFromLogValue(
					strconv.Itoa(i),
					elem,
				),
			)
		}
		return slogGroupAttrs(name, attrs)

	case log.KindMap:
		var attrs []slog.Attr
		for _, pair := range v.AsMap() {
			attrs = append(
				attrs,
				slogAttrFromLogValue(
					pair.Key,
					pair.Value,
				),
			)
		}
		return slogGroupAttrs(name, attrs)

	default:
		return slog.String(name, v.String())
	}
}

// slogGroupAttrs is an implementation of [slog.GroupAttrs] for Go versions
// prior to 1.25.
//
// TODO: remove this when go.mod is updated to Go 1.25 or later.
func slogGroupAttrs(name string, attrs []slog.Attr) slog.Attr {
	return slog.Attr{
		Key:   name,
		Value: slog.GroupValue(attrs...),
	}
}
