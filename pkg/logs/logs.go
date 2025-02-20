package logs

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

type handler struct {
	opts   *slog.HandlerOptions
	logger *log.Logger
	attrs  []slog.Attr
}

func newLogHandler(w io.Writer, opts *slog.HandlerOptions) *handler {
	return &handler{
		opts:   opts,
		logger: log.New(w, "", log.LstdFlags),
	}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return level.Level() >= h.opts.Level.Level()
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	level := r.Level.String()
	msg := r.Message

	var service string
	attrs := make([]string, 0, r.NumAttrs()+len(h.attrs))

	// Add stored attributes
	for _, attr := range h.attrs {
		if attr.Key == "service" {
			service = attr.Value.String()
		} else {
			attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any()))
		}
	}

	// Add record attributes
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "service" {
			service = a.Value.String()
		} else {
			attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value.Any()))
		}
		return true
	})

	var output string
	if service != "" {
		output = fmt.Sprintf("%s [%s]: %s", level, service, msg)
	} else {
		output = fmt.Sprintf("%s %s", level, msg)
	}

	if len(attrs) > 0 {
		output += " " + strings.Join(attrs, " ")
	}

	h.logger.Print(output)
	return nil
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

func (h *handler) WithGroup(name string) slog.Handler {
	return h
}

func (h *handler) WithService(svcName string) {
	h.WithAttrs([]slog.Attr{{Key: "service", Value: slog.StringValue(svcName)}})
}

func NewLogger(serviceName string) *slog.Logger {
	debugMode := os.Getenv("LOG_LEVEL") == "debug"
	useJSON := os.Getenv("LOG_OUTPUT") == "json"
	logLevel := slog.LevelInfo

	if debugMode {
		logLevel = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}

	var h slog.Handler
	if useJSON {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = newLogHandler(os.Stdout, opts)
	}

	logger := slog.New(h).With(
		slog.String("service", serviceName),
	)
	return logger
}
