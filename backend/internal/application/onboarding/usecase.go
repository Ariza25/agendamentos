package onboarding

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	appauth "erpagendamentos/backend/internal/application/auth"
	authdomain "erpagendamentos/backend/internal/domain/auth"
	configdomain "erpagendamentos/backend/internal/domain/configuration"
	prodomain "erpagendamentos/backend/internal/domain/professional"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"

	"github.com/google/uuid"
)

var (
	ErrBusinessNameRequired = errors.New("business name is required")
	ErrInvalidSlug          = errors.New("slug must be alphanumeric with hyphens (e.g. minha-empresa)")
	ErrSlugAlreadyTaken     = errors.New("o link/slug informado já está em uso por outro estabelecimento")
	ErrAdminNameRequired    = errors.New("admin name is required")
	ErrEmailRequired        = errors.New("email is required")
	ErrPasswordTooShort     = errors.New("password must have at least 8 characters")
)

var slugPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{1,58}[a-z0-9])?$`)

type SignupInput struct {
	BusinessName string `json:"businessName"`
	Slug         string `json:"slug"`
	AdminName    string `json:"adminName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
	Segment      string `json:"segment"`
}

type UseCase struct {
	tenantRepo       tenantdomain.Repository
	authUseCase      *appauth.UseCase
	configRepo       configdomain.Repository
	serviceRepo      servicedomain.Repository
	professionalRepo prodomain.Repository
}

func NewUseCase(
	tenantRepo tenantdomain.Repository,
	authUseCase *appauth.UseCase,
	configRepo configdomain.Repository,
	serviceRepo servicedomain.Repository,
	professionalRepo prodomain.Repository,
) *UseCase {
	return &UseCase{
		tenantRepo:       tenantRepo,
		authUseCase:      authUseCase,
		configRepo:       configRepo,
		serviceRepo:      serviceRepo,
		professionalRepo: professionalRepo,
	}
}

func (u *UseCase) Signup(input SignupInput) (authdomain.Session, tenantdomain.Tenant, error) {
	businessName := strings.TrimSpace(input.BusinessName)
	if businessName == "" {
		return authdomain.Session{}, tenantdomain.Tenant{}, ErrBusinessNameRequired
	}

	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	if slug == "" || !slugPattern.MatchString(slug) {
		return authdomain.Session{}, tenantdomain.Tenant{}, ErrInvalidSlug
	}

	adminName := strings.TrimSpace(input.AdminName)
	if adminName == "" {
		return authdomain.Session{}, tenantdomain.Tenant{}, ErrAdminNameRequired
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" || !strings.Contains(email, "@") {
		return authdomain.Session{}, tenantdomain.Tenant{}, ErrEmailRequired
	}

	if len(strings.TrimSpace(input.Password)) < 8 {
		return authdomain.Session{}, tenantdomain.Tenant{}, ErrPasswordTooShort
	}

	segment := strings.ToLower(strings.TrimSpace(input.Segment))
	if segment == "" {
		segment = "general"
	}

	// 1. Verify slug uniqueness
	if _, err := u.tenantRepo.GetBySlug(slug); err == nil {
		return authdomain.Session{}, tenantdomain.Tenant{}, ErrSlugAlreadyTaken
	}

	// 2. Create Tenant with 14-Day Free Trial
	tenantID := uuid.NewString()
	trialEnds := time.Now().Add(14 * 24 * time.Hour).Format(time.RFC3339)
	createdTenant, err := u.tenantRepo.Create(tenantdomain.Tenant{
		ID:                 tenantID,
		Name:               businessName,
		Slug:               slug,
		Segment:            segment,
		Active:             true,
		Plan:               "trial",
		SubscriptionStatus: "trialing",
		TrialEndsAt:        trialEnds,
		ActiveModules: []string{
			"agenda", "appointments", "customers", "services", "professionals",
			"finance", "commissions", "inventory", "pos_cashregister", "notifications", "reviews",
		},
		MaxProfessionals: 5,
	})
	if err != nil {
		return authdomain.Session{}, tenantdomain.Tenant{}, fmt.Errorf("create tenant: %w", err)
	}

	// 3. Register Owner Account
	session, err := u.authUseCase.RegisterOwner(authdomain.RegisterInput{
		Name:     adminName,
		Email:    email,
		Password: strings.TrimSpace(input.Password),
		Phone:    strings.TrimSpace(input.Phone),
		TenantID: createdTenant.ID,
	})
	if err != nil {
		return authdomain.Session{}, tenantdomain.Tenant{}, fmt.Errorf("create owner account: %w", err)
	}

	// 4. Seed Default Working Hours Configuration
	workingHours := []configdomain.WorkingHours{
		{DayOfWeek: 0, OpenTime: "09:00", CloseTime: "18:00", Enabled: false}, // Domingo
		{DayOfWeek: 1, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Segunda
		{DayOfWeek: 2, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Terça
		{DayOfWeek: 3, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Quarta
		{DayOfWeek: 4, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Quinta
		{DayOfWeek: 5, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Sexta
		{DayOfWeek: 6, OpenTime: "09:00", CloseTime: "18:00", Enabled: true},  // Sábado
	}
	_, _ = u.configRepo.Update(configdomain.Configuration{
		TenantID:                   createdTenant.ID,
		EstablishmentName:          businessName,
		WorkingHours:               workingHours,
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
		AllowOnlineBooking:         true,
		AdvanceNoticeHours:         2,
		BufferMinutes:              0,
		NotificationsEnabled:       true,
		NotificationReminderHours:  24,
	})

	// 5. Seed Default Professional
	pro, _ := u.professionalRepo.Create(prodomain.Professional{
		TenantID: createdTenant.ID,
		UserID:   session.User.ID,
		Name:     adminName,
		Email:    email,
		Phone:    strings.TrimSpace(input.Phone),
		Color:    "#0284c7",
		Active:   true,
	})

	// 6. Seed Segment-Specific Initial Services
	type sampleSvc struct {
		name     string
		desc     string
		duration int
		price    int
		color    string
	}
	var sampleServices []sampleSvc

	switch segment {
	case "barbearia":
		sampleServices = []sampleSvc{
			{name: "Corte de Cabelo Masculino", desc: "Corte moderno na tesoura ou máquina com lavagem e finalização.", duration: 30, price: 4500, color: "#0284c7"},
			{name: "Barba Terapia", desc: "Alinhamento com toalha quente, navalha e óleo hidratante.", duration: 30, price: 3500, color: "#b45309"},
			{name: "Combo Corte + Barba", desc: "Corte completo e barba alinhada com desconto exclusivo.", duration: 60, price: 7000, color: "#059669"},
		}
	case "salao", "estetica":
		sampleServices = []sampleSvc{
			{name: "Consulta / Avaliação Facial", desc: "Análise personalizada de pele e planejamento estético.", duration: 40, price: 8000, color: "#db2777"},
			{name: "Limpeza de Pele Profunda", desc: "Remoção de impurezas, vapor de ozônio e máscara revitalizante.", duration: 60, price: 15000, color: "#9333ea"},
			{name: "Design de Sobrancelhas", desc: "Alinhamento geométrico e acabamento impecável.", duration: 30, price: 4500, color: "#d97706"},
		}
	default:
		sampleServices = []sampleSvc{
			{name: "Atendimento Padrão", desc: "Sessão de atendimento individual com hora marcada.", duration: 45, price: 8000, color: "#0284c7"},
			{name: "Consulta Inicial / Avaliação", desc: "Primeiro contato e diagnóstico de necessidades.", duration: 30, price: 5000, color: "#059669"},
		}
	}

	for _, s := range sampleServices {
		createdSvc, err := u.serviceRepo.Create(servicedomain.Service{
			TenantID:        createdTenant.ID,
			Name:            s.name,
			Description:     s.desc,
			DurationMinutes: s.duration,
			PriceCents:      s.price,
			Color:           s.color,
			Active:          true,
		})
		if err == nil && pro.ID != "" {
			_ = u.serviceRepo.SetProfessionals(createdTenant.ID, createdSvc.ID, []string{pro.ID})
		}
	}

	return session, createdTenant, nil
}

