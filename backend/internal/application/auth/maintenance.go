package auth

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type MaintenanceConfig struct {
	Enabled   bool
	Interval  time.Duration
	Retention time.Duration
}

type MaintenanceWorker struct {
	useCase *UseCase
	config  MaintenanceConfig
}

func NewMaintenanceWorker(useCase *UseCase, config MaintenanceConfig) *MaintenanceWorker {
	if config.Interval <= 0 {
		config.Interval = time.Hour
	}
	if config.Retention <= 0 {
		config.Retention = 7 * 24 * time.Hour
	}
	return &MaintenanceWorker{useCase: useCase, config: config}
}

func (w *MaintenanceWorker) Start(ctx context.Context) {
	if w == nil || w.useCase == nil || !w.config.Enabled {
		return
	}
	go w.loop(ctx)
}

func (w *MaintenanceWorker) loop(ctx context.Context) {
	w.cleanup()
	ticker := time.NewTicker(w.config.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.cleanup()
		}
	}
}

func (w *MaintenanceWorker) cleanup() {
	summary, err := w.useCase.CleanupExpiredAuthData(w.config.Retention)
	entry := map[string]any{"message": "auth maintenance completed", "level": "INFO"}
	if err != nil {
		entry["message"] = "auth maintenance failed"
		entry["level"] = "ERROR"
		entry["error"] = err.Error()
	} else {
		entry["password_reset_tokens"] = summary.PasswordResetTokens
		entry["rate_limits"] = summary.RateLimits
		entry["sessions"] = summary.Sessions
	}
	payload, _ := json.Marshal(entry)
	log.Print(string(payload))
}
