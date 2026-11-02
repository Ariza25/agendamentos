package notification

import (
	"strings"
	"testing"
)

func sampleData() TemplateData {
	return TemplateData{
		EstablishmentName: "Barbearia Central",
		CustomerName:      "Ana Souza",
		ServiceName:       "Corte de Cabelo",
		ProfessionalName:  "Bruno Lima",
		StartsAt:          "14/08/2026 às 14:00",
		ManagementCode:    "A1-2-3",
		ManageURL:         "https://example.com/#/manage/default",
		PortalURL:         "https://example.com/#/portal/default",
	}
}

func TestRenderConfirmation(t *testing.T) {
	subject, body := RenderConfirmation(sampleData())
	if !strings.Contains(subject, "Agendamento confirmado") {
		t.Errorf("subject should mention confirmation, got %q", subject)
	}
	if !strings.Contains(subject, "Barbearia Central") {
		t.Errorf("subject should include establishment, got %q", subject)
	}
	for _, want := range []string{"Ana Souza", "Corte de Cabelo", "Bruno Lima", "14/08/2026 às 14:00", "A1-2-3", "https://example.com/#/manage/default", "https://example.com/#/portal/default"} {
		if !strings.Contains(body, want) {
			t.Errorf("confirmation body should contain %q", want)
		}
	}
}

func TestRenderReminder(t *testing.T) {
	subject, body := RenderReminder(sampleData())
	if !strings.Contains(subject, "Lembrete de agendamento") {
		t.Errorf("subject should mention reminder, got %q", subject)
	}
	for _, want := range []string{"Ana Souza", "Corte de Cabelo", "14/08/2026 às 14:00"} {
		if !strings.Contains(body, want) {
			t.Errorf("reminder body should contain %q", want)
		}
	}
}

func TestRenderCancellation(t *testing.T) {
	subject, body := RenderCancellation(sampleData())
	if !strings.Contains(subject, "Agendamento cancelado") {
		t.Errorf("subject should mention cancellation, got %q", subject)
	}
	for _, want := range []string{"Ana Souza", "Corte de Cabelo", "https://example.com/#/portal/default"} {
		if !strings.Contains(body, want) {
			t.Errorf("cancellation body should contain %q", want)
		}
	}
}

func TestRenderReschedule(t *testing.T) {
	subject, body := RenderReschedule(sampleData())
	if !strings.Contains(subject, "Agendamento remarcado") {
		t.Errorf("subject should mention reschedule, got %q", subject)
	}
	for _, want := range []string{"Ana Souza", "Corte de Cabelo", "Novo horário", "14/08/2026 às 14:00"} {
		if !strings.Contains(body, want) {
			t.Errorf("reschedule body should contain %q", want)
		}
	}
}

func TestRenderConfirmationOmitsEmptyFields(t *testing.T) {
	data := sampleData()
	data.ManagementCode = ""
	data.ManageURL = ""
	data.ProfessionalName = ""
	_, body := RenderConfirmation(data)
	if strings.Contains(body, "Profissional") {
		t.Errorf("body should not include empty professional line")
	}
	if strings.Contains(body, "Código do agendamento") {
		t.Errorf("body should not include empty management code section")
	}
	if strings.Contains(body, "Cancelar ou reagendar") {
		t.Errorf("body should not include manage footer without code")
	}
}

func TestGreetingWithoutName(t *testing.T) {
	subject, body := RenderReminder(TemplateData{EstablishmentName: "X", StartsAt: "14/08/2026 às 14:00"})
	if !strings.Contains(body, "Olá,") {
		t.Errorf("body should use generic greeting, got %q", body)
	}
	if subject == "" {
		t.Errorf("subject should not be empty")
	}
}
