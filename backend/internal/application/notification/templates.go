package notification

import (
	"os"
	"strings"
)

type TemplateData struct {
	EstablishmentName string
	CustomerName      string
	ServiceName       string
	ProfessionalName  string
	StartsAt          string
	ManagementCode    string
	ManageURL         string
	PortalURL         string
}

func PublicBookingOrigin() string {
	for _, envKey := range []string{"WIDGET_BOOKING_ORIGIN", "BOOKING_BASE_URL", "BASE_URL"} {
		if raw := strings.TrimSpace(os.Getenv(envKey)); raw != "" {
			return strings.TrimSuffix(raw, "/")
		}
	}
	return ""
}

func RenderConfirmation(data TemplateData) (string, string) {
	subject := "Agendamento confirmado"
	if name := strings.TrimSpace(data.EstablishmentName); name != "" {
		subject += " - " + name
	}
	body := strings.Builder{}
	body.WriteString(greeting(data.CustomerName))
	body.WriteString("\n\nSeu agendamento foi confirmado:\n\n")
	writeLine(&body, "Serviço", data.ServiceName)
	writeLine(&body, "Profissional", data.ProfessionalName)
	writeLine(&body, "Horário", data.StartsAt)
	writeLine(&body, "Estabelecimento", data.EstablishmentName)
	writeLine(&body, "Código do agendamento", data.ManagementCode)
	writeFooter(&body, data)
	return subject, body.String()
}

func RenderReminder(data TemplateData) (string, string) {
	subject := "Lembrete de agendamento"
	if name := strings.TrimSpace(data.EstablishmentName); name != "" {
		subject += " - " + name
	}
	body := strings.Builder{}
	body.WriteString(greeting(data.CustomerName))
	body.WriteString("\n\nPassando para lembrar do seu agendamento:\n\n")
	writeLine(&body, "Serviço", data.ServiceName)
	writeLine(&body, "Profissional", data.ProfessionalName)
	writeLine(&body, "Horário", data.StartsAt)
	writeLine(&body, "Estabelecimento", data.EstablishmentName)
	writeLine(&body, "Código do agendamento", data.ManagementCode)
	writeFooter(&body, data)
	return subject, body.String()
}

func RenderCancellation(data TemplateData) (string, string) {
	subject := "Agendamento cancelado"
	if name := strings.TrimSpace(data.EstablishmentName); name != "" {
		subject += " - " + name
	}
	body := strings.Builder{}
	body.WriteString(greeting(data.CustomerName))
	body.WriteString("\n\nSeu agendamento foi cancelado:\n\n")
	writeLine(&body, "Serviço", data.ServiceName)
	writeLine(&body, "Profissional", data.ProfessionalName)
	writeLine(&body, "Horário", data.StartsAt)
	writeLine(&body, "Estabelecimento", data.EstablishmentName)
	body.WriteString("\nO horário voltou a ficar disponível. Para fazer um novo agendamento, acesse o portal do cliente:")
	if portal := strings.TrimSpace(data.PortalURL); portal != "" {
		body.WriteString("\n" + portal)
	}
	return subject, body.String()
}

func RenderReschedule(data TemplateData) (string, string) {
	subject := "Agendamento remarcado"
	if name := strings.TrimSpace(data.EstablishmentName); name != "" {
		subject += " - " + name
	}
	body := strings.Builder{}
	body.WriteString(greeting(data.CustomerName))
	body.WriteString("\n\nSeu agendamento foi remarcado. Confira o novo horário:\n\n")
	writeLine(&body, "Serviço", data.ServiceName)
	writeLine(&body, "Profissional", data.ProfessionalName)
	writeLine(&body, "Novo horário", data.StartsAt)
	writeLine(&body, "Estabelecimento", data.EstablishmentName)
	writeFooter(&body, data)
	return subject, body.String()
}

func greeting(name string) string {
	if strings.TrimSpace(name) != "" {
		return "Olá " + strings.TrimSpace(name) + ","
	}
	return "Olá,"
}

func writeLine(builder *strings.Builder, label string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	builder.WriteString("• " + label + ": " + value + "\n")
}

func writeFooter(builder *strings.Builder, data TemplateData) {
	if code := strings.TrimSpace(data.ManagementCode); code != "" {
		builder.WriteString("\nCancelar ou reagendar usando o código:\n")
		if manage := strings.TrimSpace(data.ManageURL); manage != "" {
			builder.WriteString(manage + "\n")
		}
	}
	if portal := strings.TrimSpace(data.PortalURL); portal != "" {
		builder.WriteString("\nPortal do cliente:\n" + portal + "\n")
	}
	builder.WriteString("\nAté lá!")
}
