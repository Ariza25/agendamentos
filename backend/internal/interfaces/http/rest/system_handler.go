package rest

import (
	"net/http"
	"runtime"
	"time"

	appoutbox "erpagendamentos/backend/internal/application/outbox"
)

type SystemHandler struct {
	readinessChecker func() error
	dbStatsGetter    func() map[string]any
	outboxUseCase    *appoutbox.UseCase
	startedAt        time.Time
}

func NewSystemHandler(readinessChecker func() error, dbStatsGetter func() map[string]any, outboxUseCase *appoutbox.UseCase) *SystemHandler {
	return &SystemHandler{
		readinessChecker: readinessChecker,
		dbStatsGetter:    dbStatsGetter,
		outboxUseCase:    outboxUseCase,
		startedAt:        time.Now(),
	}
}

func (h *SystemHandler) Status(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dbConnected := true
	var dbError string
	if h.readinessChecker != nil {
		if err := h.readinessChecker(); err != nil {
			dbConnected = false
			dbError = err.Error()
		}
	}

	var dbStats map[string]any
	if h.dbStatsGetter != nil {
		dbStats = h.dbStatsGetter()
	}

	outboxStats := map[string]int{"pending": 0, "published": 0, "failed": 0}
	if h.outboxUseCase != nil {
		if stats, err := h.outboxUseCase.Stats(); err == nil {
			outboxStats = stats
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "healthy",
		"database": map[string]any{
			"connected": dbConnected,
			"error":     dbError,
			"stats":     dbStats,
		},
		"outbox":        outboxStats,
		"uptimeSeconds": int(time.Since(h.startedAt).Seconds()),
		"goroutines":    runtime.NumGoroutine(),
		"memoryAllocMb": float64(m.Alloc) / 1024 / 1024,
		"memorySysMb":   float64(m.Sys) / 1024 / 1024,
		"serverTime":    time.Now().UTC().Format(time.RFC3339),
	})
}
