package audit

import (
	"testing"
	"time"

	auditdomain "erpagendamentos/backend/internal/domain/audit"
	"erpagendamentos/backend/internal/infrastructure/persistence/memory"
)

func TestAuditLogRecordingAndListing(t *testing.T) {
	repo := memory.NewAuditRepository()
	uc := NewUseCase(repo)

	err := uc.Log(auditdomain.AuditLog{
		TenantID:   "tenant-1",
		UserID:     "user-1",
		UserEmail:  "admin@tenant1.local",
		Action:     "cash_register.close",
		Resource:   "cash_register",
		ResourceID: "shift-123",
		Details: map[string]any{
			"closingBalanceCents": 15000,
		},
		IPAddress: "192.168.1.100",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("expected nil error recording audit log, got: %v", err)
	}

	// Record for tenant 2
	_ = uc.Log(auditdomain.AuditLog{
		TenantID:   "tenant-2",
		UserID:     "user-2",
		UserEmail:  "owner@tenant2.local",
		Action:     "tenant.update_plan",
		Resource:   "tenant",
		ResourceID: "tenant-2",
		CreatedAt:  time.Now(),
	})

	// Master admin listing (tenantID = "") should see all 2 logs
	allLogs, total, err := uc.List("", 10, 0)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if total != 2 || len(allLogs) != 2 {
		t.Fatalf("expected 2 logs for master, got %d (total: %d)", len(allLogs), total)
	}

	// Filter by tenant-1 should return only 1 log
	t1Logs, t1Total, err := uc.List("tenant-1", 10, 0)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if t1Total != 1 || len(t1Logs) != 1 {
		t.Fatalf("expected 1 log for tenant-1, got %d", len(t1Logs))
	}
	if t1Logs[0].Action != "cash_register.close" {
		t.Fatalf("expected action cash_register.close, got %s", t1Logs[0].Action)
	}
}

