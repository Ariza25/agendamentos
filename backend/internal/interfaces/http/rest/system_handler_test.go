package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appaudit "erpagendamentos/backend/internal/application/audit"
	appoutbox "erpagendamentos/backend/internal/application/outbox"
	auditdomain "erpagendamentos/backend/internal/domain/audit"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestReadinessEndpoint(t *testing.T) {
	checkerCalled := false
	checker := func() error {
		checkerCalled = true
		return nil
	}

	handler := readiness(checker)
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !checkerCalled {
		t.Fatalf("expected checker to be called")
	}

	var res map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	if res["status"] != "ready" || res["database"] != "connected" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestAuditAndSystemHandlers(t *testing.T) {
	auditRepo := memory.NewAuditRepository()
	auditUC := appaudit.NewUseCase(auditRepo)
	_ = auditUC.Log(auditdomain.AuditLog{
		TenantID:   "tenant-test",
		Action:     "checkout.completed",
		Resource:   "appointment",
		ResourceID: "apt-123",
		CreatedAt:  time.Now(),
	})

	outboxRepo := memory.NewOutboxRepository()
	outboxUC := appoutbox.NewUseCase(outboxRepo, nil)
	_ = outboxUC.Publish("tenant-test", "checkout.completed", map[string]any{"id": "apt-123"})

	// Test Audit Handler
	auditHandler := NewAuditHandler(auditUC)
	auditReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit-logs", nil)
	auditRec := httptest.NewRecorder()
	auditHandler.List(auditRec, auditReq)

	if auditRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on audit logs, got %d", auditRec.Code)
	}
	var auditRes map[string]any
	_ = json.Unmarshal(auditRec.Body.Bytes(), &auditRes)
	if auditRes["total"] != float64(1) {
		t.Fatalf("expected 1 audit log, got %v", auditRes["total"])
	}

	// Test System Handler
	systemHandler := NewSystemHandler(
		func() error { return nil },
		func() map[string]any { return map[string]any{"driver": "memory", "openConnections": 1} },
		outboxUC,
	)
	sysReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/system/status", nil)
	sysRec := httptest.NewRecorder()
	systemHandler.Status(sysRec, sysReq)

	if sysRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on system status, got %d", sysRec.Code)
	}
	var sysRes map[string]any
	_ = json.Unmarshal(sysRec.Body.Bytes(), &sysRes)
	if sysRes["status"] != "healthy" {
		t.Fatalf("expected healthy status, got %v", sysRes["status"])
	}
	outboxData, ok := sysRes["outbox"].(map[string]any)
	if !ok || outboxData["pending"] != float64(1) {
		t.Fatalf("expected 1 pending outbox item, got %+v", sysRes["outbox"])
	}
}
