package rest

import (
	"net/http"
	"strconv"
	"time"

	appappointment "erpagendamentos/backend/internal/application/appointment"
	appcustomer "erpagendamentos/backend/internal/application/customer"
	apppayment "erpagendamentos/backend/internal/application/payment"
	appprofessional "erpagendamentos/backend/internal/application/professional"
	appresource "erpagendamentos/backend/internal/application/resource"
	appreview "erpagendamentos/backend/internal/application/review"
	appservice "erpagendamentos/backend/internal/application/service"
)

type ExportHandler struct {
	appointments  *appappointment.UseCase
	customers     *appcustomer.UseCase
	payments      *apppayment.UseCase
	professionals *appprofessional.UseCase
	services      *appservice.UseCase
	reviews       *appreview.UseCase
	resources     *appresource.UseCase
}

func NewExportHandler(
	appointments *appappointment.UseCase,
	customers *appcustomer.UseCase,
	payments *apppayment.UseCase,
	professionals *appprofessional.UseCase,
	services *appservice.UseCase,
	reviews *appreview.UseCase,
	resources *appresource.UseCase,
) *ExportHandler {
	return &ExportHandler{
		appointments:  appointments,
		customers:     customers,
		payments:      payments,
		professionals: professionals,
		services:      services,
		reviews:       reviews,
		resources:     resources,
	}
}

func (h *ExportHandler) ExportAppointments(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	if from == "" {
		from = time.Now().AddDate(0, 0, -30).Format(time.RFC3339)
	}
	if to == "" {
		to = time.Now().AddDate(0, 0, 90).Format(time.RFC3339)
	}

	items, err := h.appointments.ListByRange(tenantID, from, to, query.Get("professionalId"), query.Get("customerId"), query.Get("status"))
	if err != nil {
		writeInternalError(w, "Could not export the appointments.")
		return
	}

	customerNames := h.customerNameMap(tenantID)
	serviceNames := h.serviceNameMap(tenantID)
	professionalNames := h.professionalNameMap(tenantID)
	resourceNames := h.resourceNameMap(tenantID)

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			item.StartsAt,
			item.EndsAt,
			customerNames[item.CustomerID],
			serviceNames[item.ServiceID],
			professionalNames[item.ProfessionalID],
			resourceNames[item.ResourceID],
			item.Status,
			item.Notes,
		})
	}
	writeCSV(w, "agendamentos.csv", []string{"inicio", "fim", "cliente", "servico", "profissional", "recurso", "status", "notas"}, rows)
}

func (h *ExportHandler) ExportCustomers(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	items, err := h.customers.ListByTenant(tenantID, r.URL.Query().Get("search"))
	if err != nil {
		writeInternalError(w, "Could not export the customers.")
		return
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			item.Name,
			item.Phone,
			item.Email,
			item.Birthday,
			item.Notes,
			item.CreatedAt,
		})
	}
	writeCSV(w, "clientes.csv", []string{"nome", "telefone", "email", "aniversario", "notas", "criado_em"}, rows)
}

func (h *ExportHandler) ExportPayments(w http.ResponseWriter, r *http.Request) {
	tenantID := resolveRequestTenantID(r, nil)
	query := r.URL.Query()
	items, err := h.payments.ListByRange(tenantID, query.Get("from"), query.Get("to"))
	if err != nil {
		writeInternalError(w, "Could not export the payments.")
		return
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			item.CreatedAt,
			item.Description,
			item.Method,
			item.Provider,
			item.Status,
			strconv.Itoa(item.AmountCents),
		})
	}
	writeCSV(w, "pagamentos.csv", []string{"criado_em", "descricao", "metodo", "provedor", "status", "valor_centavos"}, rows)
}

func (h *ExportHandler) ExportReviews(w http.ResponseWriter, r *http.Request) {
	if h.reviews == nil {
		writeInternalError(w, "Could not export the reviews.")
		return
	}
	tenantID := resolveRequestTenantID(r, nil)
	items, err := h.reviews.ListByTenant(tenantID, "")
	if err != nil {
		writeInternalError(w, "Could not export the reviews.")
		return
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			item.CreatedAt,
			item.CustomerName,
			item.Status,
			strconv.Itoa(item.Rating),
			item.Comment,
		})
	}
	writeCSV(w, "avaliacoes.csv", []string{"criado_em", "cliente", "status", "nota", "comentario"}, rows)
}

func (h *ExportHandler) customerNameMap(tenantID string) map[string]string {
	result := map[string]string{}
	if h.customers == nil {
		return result
	}
	items, err := h.customers.ListByTenant(tenantID, "")
	if err != nil {
		return result
	}
	for _, item := range items {
		result[item.ID] = item.Name
	}
	return result
}

func (h *ExportHandler) serviceNameMap(tenantID string) map[string]string {
	result := map[string]string{}
	if h.services == nil {
		return result
	}
	items, err := h.services.ListByTenant(tenantID)
	if err != nil {
		return result
	}
	for _, item := range items {
		result[item.ID] = item.Name
	}
	return result
}

func (h *ExportHandler) professionalNameMap(tenantID string) map[string]string {
	result := map[string]string{}
	if h.professionals == nil {
		return result
	}
	items, err := h.professionals.ListByTenant(tenantID)
	if err != nil {
		return result
	}
	for _, item := range items {
		result[item.ID] = item.Name
	}
	return result
}

func (h *ExportHandler) resourceNameMap(tenantID string) map[string]string {
	result := map[string]string{}
	if h.resources == nil {
		return result
	}
	items, err := h.resources.ListByTenant(tenantID)
	if err != nil {
		return result
	}
	for _, item := range items {
		result[item.ID] = item.Name
	}
	return result
}
