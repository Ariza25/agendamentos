package configuration

type WorkingHours struct {
	DayOfWeek int    `json:"dayOfWeek"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	Enabled   bool   `json:"enabled"`
}

type Configuration struct {
	TenantID                   string         `json:"tenantId"`
	EstablishmentName          string         `json:"establishmentName"`
	WorkingHours               []WorkingHours `json:"workingHours"`
	AppointmentIntervalMinutes int            `json:"appointmentIntervalMinutes"`
	DefaultDurationMinutes     int            `json:"defaultDurationMinutes"`
	AllowOnlineBooking         bool           `json:"allowOnlineBooking"`
	AdvanceNoticeHours         int            `json:"advanceNoticeHours"`
	UpdatedAt                  string         `json:"updatedAt"`
}

type Repository interface {
	Get(tenantID string) (Configuration, error)
	Update(input Configuration) (Configuration, error)
}
