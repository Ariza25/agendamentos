package appointment

type Appointment struct {
	ID                 string `json:"id"`
	TenantID           string `json:"tenantId"`
	ProfessionalID     string `json:"professionalId"`
	ServiceID          string `json:"serviceId"`
	ResourceID         string `json:"resourceId,omitempty"`
	CustomerID         string `json:"customerId"`
	StartsAt           string `json:"startsAt"`
	EndsAt             string `json:"endsAt"`
	Status             string `json:"status"`
	Notes              string `json:"notes"`
	CreatedBy          string `json:"createdBy"`
	ManagementCodeHash string `json:"-"`
	SeriesID           string `json:"seriesId,omitempty"`
	PaymentStatus      string `json:"paymentStatus,omitempty"`
	NoShowFeeCents     int    `json:"noShowFeeCents,omitempty"`
	Source             string `json:"source,omitempty"`
	CancellationReason string `json:"cancellationReason,omitempty"`
	CheckInAt          string `json:"checkInAt,omitempty"`
	NoShowAt           string `json:"noShowAt,omitempty"`
	ReminderSentAt     string `json:"reminderSentAt,omitempty"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type Analytics struct {
	TotalAppointments int                  `json:"totalAppointments"`
	Scheduled         int                  `json:"scheduled"`
	Confirmed         int                  `json:"confirmed"`
	Completed         int                  `json:"completed"`
	Cancelled         int                  `json:"cancelled"`
	NoShow            int                  `json:"noShow"`
	NoShowRate        float64              `json:"noShowRate"`
	RevenueCents      int64                `json:"revenueCents"`
	ByService         []ServiceMetric      `json:"byService"`
	ByProfessional    []ProfessionalMetric `json:"byProfessional"`
	ByDay             []DayMetric          `json:"byDay"`
}

type ServiceMetric struct {
	ServiceID    string `json:"serviceId"`
	ServiceName  string `json:"serviceName"`
	Count        int    `json:"count"`
	RevenueCents int64  `json:"revenueCents"`
}

type ProfessionalMetric struct {
	ProfessionalID   string `json:"professionalId"`
	ProfessionalName string `json:"professionalName"`
	Count            int    `json:"count"`
}

type DayMetric struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type Repository interface {
	Create(input Appointment) (Appointment, error)
	Update(input Appointment) (Appointment, error)
	GetByID(tenantID string, id string) (Appointment, error)
	GetByManagementCodeHash(tenantID string, hash string) (Appointment, error)
	ListByRange(tenantID string, from string, to string, professionalID string, customerID string, status string) ([]Appointment, error)
	ListDueReminders(from string, to string) ([]Appointment, error)
	UpdateStatus(tenantID string, id string, status string, updatedBy string) (Appointment, error)
	SetReminderSent(tenantID string, id string) error
	HasConflict(tenantID string, professionalID string, startsAt string, endsAt string, excludeID string) (bool, error)
	CountResourceOverlaps(tenantID string, resourceID string, startsAt string, endsAt string, excludeID string) (int, error)
	CountByStatus(tenantID string, from string, to string) ([]StatusCount, error)
	Analytics(tenantID string, from string, to string) (Analytics, error)
	Delete(tenantID string, id string) error
}
