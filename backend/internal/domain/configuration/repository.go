package configuration

type WorkingHours struct {
	DayOfWeek int    `json:"dayOfWeek"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	Enabled   bool   `json:"enabled"`
}

type ClosedDay struct {
	Date      string `json:"date"`
	Name      string `json:"name"`
	AllDay    bool   `json:"allDay"`
	OpenTime  string `json:"openTime,omitempty"`
	CloseTime string `json:"closeTime,omitempty"`
}

type Configuration struct {
	TenantID                   string         `json:"tenantId"`
	EstablishmentName          string         `json:"establishmentName"`
	WorkingHours               []WorkingHours `json:"workingHours"`
	AppointmentIntervalMinutes int            `json:"appointmentIntervalMinutes"`
	DefaultDurationMinutes     int            `json:"defaultDurationMinutes"`
	AllowOnlineBooking         bool           `json:"allowOnlineBooking"`
	AdvanceNoticeHours         int            `json:"advanceNoticeHours"`
	BufferMinutes              int            `json:"bufferMinutes"`
	ClosedDays                 []ClosedDay    `json:"closedDays"`
	NotificationsEnabled       bool           `json:"notificationsEnabled"`
	NotificationReminderHours  int            `json:"notificationReminderHours"`
	NotificationChannel        string         `json:"notificationChannel"`
	DepositEnabled             bool           `json:"depositEnabled"`
	DepositType                string         `json:"depositType"`
	DepositAmountCents         int            `json:"depositAmountCents"`
	DepositPercentage          int            `json:"depositPercentage"`
	NoShowFeeCents             int            `json:"noShowFeeCents"`
	CancellationPolicyHours    int            `json:"cancellationPolicyHours"`
	CalendarSyncEnabled        bool           `json:"calendarSyncEnabled"`
	LogoURL                    string         `json:"logoUrl,omitempty"`
	PrimaryColor               string         `json:"primaryColor,omitempty"`
	Address                    string         `json:"address,omitempty"`
	AddressNumber              string         `json:"addressNumber,omitempty"`
	Neighborhood               string         `json:"neighborhood,omitempty"`
	City                       string         `json:"city,omitempty"`
	State                      string         `json:"state,omitempty"`
	PostalCode                 string         `json:"postalCode,omitempty"`
	ContactPhone               string         `json:"contactPhone,omitempty"`
	ContactWhatsApp            string         `json:"contactWhatsApp,omitempty"`
	Instagram                  string         `json:"instagram,omitempty"`
	Facebook                   string         `json:"facebook,omitempty"`
	Website                    string         `json:"website,omitempty"`
	Description                string         `json:"description,omitempty"`
	UpdatedAt                  string         `json:"updatedAt"`
}

type Repository interface {
	Get(tenantID string) (Configuration, error)
	Update(input Configuration) (Configuration, error)
}
