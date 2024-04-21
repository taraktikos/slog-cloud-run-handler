package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// LevelCritical is an extra log level supported by Cloud Logging
const (
	LevelCritical = slog.Level(12)
)

// CloudLoggingHandler Handler that outputs JSON understood by the structured log agent.
// See https://cloud.google.com/logging/docs/agent/logging/configuration#special-fields
type CloudLoggingHandler struct{ handler slog.Handler }

func NewCloudLoggingHandler(projectID string) *CloudLoggingHandler {
	return &CloudLoggingHandler{handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			} else if a.Key == slog.SourceKey {
				a.Key = "logging.googleapis.com/sourceLocation"
			} else if a.Key == slog.LevelKey {
				a.Key = "severity"
				if level, ok := a.Value.Any().(slog.Level); ok {
					if level == LevelCritical {
						a.Value = slog.StringValue("CRITICAL")
					}
				}
			} else if a.Key == "trace-id" {
				a.Key = "logging.googleapis.com/trace"
				if traceID, ok := a.Value.Any().(string); ok {
					a.Value = slog.StringValue(fmt.Sprintf("projects/%s/traces/%s", projectID, traceID))
				}
			} else if a.Key == "span-id" {
				a.Key = "logging.googleapis.com/spanId"
			}
			return a
		},
	})}
}

func (h *CloudLoggingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CloudLoggingHandler) Handle(ctx context.Context, rec slog.Record) error {
	//trace := traceFromContext(ctx)
	//if trace != "" {
	//	rec = rec.Clone()
	//	rec.
	//		// Add trace ID	to the record so it is correlated with the Cloud Run request log
	//		// See https://cloud.google.com/trace/docs/trace-log-integration
	//		rec.Add("logging.googleapis.com/trace", slog.StringValue(trace))
	//}
	return h.handler.Handle(ctx, rec)
}

func (h *CloudLoggingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *CloudLoggingHandler) WithGroup(name string) slog.Handler {
	return &CloudLoggingHandler{handler: h.handler.WithGroup(name)}
}
