package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/tomek7667/go-multi-logger-slog/internal/utils"
)

type LokiHandler struct {
	lokiURL      string
	username     string
	password     string
	defaultLevel slog.Level
	appName      string
}

func NewLokiHandler(lokiURL, username, password, appName string, defaultLevel slog.Level) *LokiHandler {
	return &LokiHandler{
		lokiURL:      lokiURL,
		username:     username,
		password:     password,
		defaultLevel: defaultLevel,
		appName:      appName,
	}
}

func getRecordMessage(record slog.Record) string {
	writer := &bytes.Buffer{}
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     record.Level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				tmS := a.Value.String()
				tm, err := time.Parse(time.RFC3339Nano, tmS)
				if err != nil {
					return a
				}
				return slog.Attr{
					Key:   "time",
					Value: slog.AnyValue(tm.Format(time.DateTime)),
				}
			} else {
				return a
			}
		},
	})
	handler.Handle(context.Background(), record)
	s := writer.String()
	return s
}

func (h *LokiHandler) Handle(ctx context.Context, record slog.Record) error {
	labels := map[string]string{
		"job":   h.appName,
		"level": record.Level.String(),
		"env":   utils.Getenv("ENVIRONMENT", "dev"),
	}

	entry := [][]string{
		{
			fmt.Sprintf("%d", time.Now().UnixNano()),
			getRecordMessage(record),
		},
	}
	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": labels,
				"values": entry,
			},
		},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal log payload: %w", err)
	}

	req, err := http.NewRequest("POST", h.lokiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if h.username != "" && h.password != "" {
		req.SetBasicAuth(h.username, h.password)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send log to Loki: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("loki responded with status: %s", resp.Status)
	}
	return nil
}

func (h *LokiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *LokiHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *LokiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level <= h.defaultLevel
}

func (h *LokiHandler) Level() slog.Level {
	return h.defaultLevel
}
