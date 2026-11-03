package bootstrap

import (
	"log"
	"strings"
	"time"

	appointmentdomain "erpagendamentos/backend/internal/domain/appointment"
	authdomain "erpagendamentos/backend/internal/domain/auth"
	cashregisterdomain "erpagendamentos/backend/internal/domain/cashregister"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
	customerdomain "erpagendamentos/backend/internal/domain/customer"
	financedomain "erpagendamentos/backend/internal/domain/finance"
	professionaldomain "erpagendamentos/backend/internal/domain/professional"
	resourcedomain "erpagendamentos/backend/internal/domain/resource"
	reviewdomain "erpagendamentos/backend/internal/domain/review"
	servicedomain "erpagendamentos/backend/internal/domain/service"
	stockdomain "erpagendamentos/backend/internal/domain/stock"
	tenantdomain "erpagendamentos/backend/internal/domain/tenant"
	"erpagendamentos/backend/internal/application/auth"

	"github.com/google/uuid"
)

func seedHomologationData(repos repositorySet, hasher auth.PasswordHasher) {
	log.Println(`{"level":"INFO","message":"seeding homologation data (dom navalha case)..."}`)

	const tenantID = "default"

	// 1. Tenant
	tenant, err := repos.tenant.GetByID(tenantID)
	tenantPayload := tenantdomain.Tenant{
		ID:                 tenantID,
		Name:               "Barbearia & Studio Dom Navalha",
		Slug:               "dom-navalha",
		Segment:            "barbearia",
		Active:             true,
		Plan:               "pro",
		SubscriptionStatus: "active",
		ActiveModules: []string{
			"agenda", "appointments", "customers", "services", "professionals",
			"finance", "commissions", "inventory", "pos_cashregister", "notifications", "reviews",
		},
		MaxProfessionals: 10,
	}
	if err != nil {
		if _, err := repos.tenant.Create(tenantPayload); err != nil {
			log.Printf(`{"level":"WARN","message":"seed tenant create failed","error":"%v"}`, err)
		}
	} else {
		tenantPayload.Active = tenant.Active
		if _, err := repos.tenant.Update(tenantPayload); err != nil {
			log.Printf(`{"level":"WARN","message":"seed tenant update failed","error":"%v"}`, err)
		}
	}

	// 2. Configuration
	workingHours := []configurationdomain.WorkingHours{
		{DayOfWeek: 0, OpenTime: "09:00", CloseTime: "18:00", Enabled: false}, // Domingo fechado
		{DayOfWeek: 1, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Segunda
		{DayOfWeek: 2, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Terça
		{DayOfWeek: 3, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Quarta
		{DayOfWeek: 4, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Quinta
		{DayOfWeek: 5, OpenTime: "09:00", CloseTime: "19:00", Enabled: true},  // Sexta
		{DayOfWeek: 6, OpenTime: "09:00", CloseTime: "18:00", Enabled: true},  // Sábado
	}
	_, _ = repos.configuration.Update(configurationdomain.Configuration{
		TenantID:                   tenantID,
		EstablishmentName:          "Barbearia & Studio Dom Navalha",
		WorkingHours:               workingHours,
		AppointmentIntervalMinutes: 30,
		DefaultDurationMinutes:     30,
		AllowOnlineBooking:         true,
		AdvanceNoticeHours:         2,
		BufferMinutes:              0,
		NotificationsEnabled:       true,
		NotificationReminderHours:  24,
		PrimaryColor:               "#1e3a8a",
	})

	// 3. User Accounts (Owner, Admin, Test)
	ownerUserID := seedAccount(repos.auth, hasher, "Carlos Eduardo Oliveira", "carlos.proprietario@domnavalha.com", "(11) 98765-4321", "DomNavalha@2026", "owner", tenantID)
	seedAccount(repos.auth, hasher, "Administrador Dom Navalha", "admin@erpagendamentos.com", "(11) 98765-4321", "Admin@123456", "owner", tenantID)
	seedAccount(repos.auth, hasher, "Carlos Proprietário", "teste@gmail.com", "(11) 98765-4321", "Teste@123", "owner", tenantID)

	// 4. Professionals
	pros, _ := repos.professional.ListByTenant(tenantID)
	var carlosProID, rafaelProID, lucasProID string
	if len(pros) == 0 {
		carlosPro, _ := repos.professional.Create(professionaldomain.Professional{
			TenantID:             tenantID,
			UserID:               ownerUserID,
			Name:                 "Carlos Oliveira",
			Email:                "carlos.proprietario@domnavalha.com",
			Phone:                "(11) 98765-4321",
			Color:                "#0284c7",
			CommissionPercentage: 70.0,
			Active:               true,
		})
		carlosProID = carlosPro.ID

		rafaelPro, _ := repos.professional.Create(professionaldomain.Professional{
			TenantID:             tenantID,
			Name:                 "Rafael Santos",
			Email:                "rafael.barber@domnavalha.com",
			Phone:                "(11) 97654-3210",
			Color:                "#b45309",
			CommissionPercentage: 50.0,
			Active:               true,
		})
		rafaelProID = rafaelPro.ID

		lucasPro, _ := repos.professional.Create(professionaldomain.Professional{
			TenantID:             tenantID,
			Name:                 "Lucas Mendes",
			Email:                "lucas.estilo@domnavalha.com",
			Phone:                "(11) 96543-2109",
			Color:                "#059669",
			CommissionPercentage: 50.0,
			Active:               true,
		})
		lucasProID = lucasPro.ID
	} else {
		for _, p := range pros {
			if strings.Contains(p.Name, "Carlos") {
				carlosProID = p.ID
			} else if strings.Contains(p.Name, "Rafael") {
				rafaelProID = p.ID
			} else if strings.Contains(p.Name, "Lucas") {
				lucasProID = p.ID
			}
		}
		if carlosProID == "" && len(pros) > 0 {
			carlosProID = pros[0].ID
		}
		if rafaelProID == "" && len(pros) > 1 {
			rafaelProID = pros[1].ID
		} else if rafaelProID == "" {
			rafaelProID = carlosProID
		}
		if lucasProID == "" && len(pros) > 2 {
			lucasProID = pros[2].ID
		} else if lucasProID == "" {
			lucasProID = carlosProID
		}
	}

	// 5. Services
	svcs, _ := repos.service.ListByTenant(tenantID)
	var corteSvcID, barbaSvcID, comboSvcID, acabamentoSvcID, selagemSvcID string
	if len(svcs) == 0 {
		corte, _ := repos.service.Create(servicedomain.Service{
			TenantID:        tenantID,
			Name:            "Corte de Cabelo Masculino",
			Description:     "Corte moderno na tesoura ou máquina, lavagem com shampoo mentolado e finalização impecável com pomada matte.",
			DurationMinutes: 30,
			PriceCents:      5000,
			Color:           "#0284c7",
			Active:          true,
		})
		corteSvcID = corte.ID
		_ = repos.service.SetProfessionals(tenantID, corte.ID, []string{carlosProID, rafaelProID, lucasProID})

		barba, _ := repos.service.Create(servicedomain.Service{
			TenantID:        tenantID,
			Name:            "Barboterapia Completa",
			Description:     "Alinhamento e desenho com toalha quente aromatizada, óleos essenciais, navalhete descartável e hidratação facial.",
			DurationMinutes: 30,
			PriceCents:      4000,
			Color:           "#b45309",
			Active:          true,
		})
		barbaSvcID = barba.ID
		_ = repos.service.SetProfessionals(tenantID, barba.ID, []string{carlosProID, rafaelProID})

		combo, _ := repos.service.Create(servicedomain.Service{
			TenantID:        tenantID,
			Name:            "Combo Corte + Barba",
			Description:     "Experiência completa Dom Navalha com corte personalizado, barboterapia relaxante e finalização premium.",
			DurationMinutes: 60,
			PriceCents:      8500,
			Color:           "#059669",
			Active:          true,
		})
		comboSvcID = combo.ID
		_ = repos.service.SetProfessionals(tenantID, combo.ID, []string{carlosProID, rafaelProID, lucasProID})

		acabamento, _ := repos.service.Create(servicedomain.Service{
			TenantID:        tenantID,
			Name:            "Acabamento & Sobrancelha",
			Description:     "Pézinho navalhado, alinhamento de costeletas e design de sobrancelha com navalha/pinça.",
			DurationMinutes: 20,
			PriceCents:      2500,
			Color:           "#6366f1",
			Active:          true,
		})
		acabamentoSvcID = acabamento.ID
		_ = repos.service.SetProfessionals(tenantID, acabamento.ID, []string{carlosProID, rafaelProID, lucasProID})

		selagem, _ := repos.service.Create(servicedomain.Service{
			TenantID:        tenantID,
			Name:            "Selagem Capilar & Tratamento",
			Description:     "Tratamento capilar profundo para redução de frizz, alinhamento dos fios e brilho natural prolongado.",
			DurationMinutes: 60,
			PriceCents:      12000,
			Color:           "#d97706",
			Active:          true,
		})
		selagemSvcID = selagem.ID
		_ = repos.service.SetProfessionals(tenantID, selagem.ID, []string{carlosProID, lucasProID})
	} else {
		for _, s := range svcs {
			if strings.Contains(s.Name, "Corte") && !strings.Contains(s.Name, "Combo") {
				corteSvcID = s.ID
			} else if strings.Contains(s.Name, "Barba") && !strings.Contains(s.Name, "Combo") {
				barbaSvcID = s.ID
			} else if strings.Contains(s.Name, "Combo") {
				comboSvcID = s.ID
			} else if strings.Contains(s.Name, "Acabamento") {
				acabamentoSvcID = s.ID
			} else if strings.Contains(s.Name, "Selagem") {
				selagemSvcID = s.ID
			}
		}
		if corteSvcID == "" && len(svcs) > 0 {
			corteSvcID = svcs[0].ID
		}
		if barbaSvcID == "" && len(svcs) > 1 {
			barbaSvcID = svcs[1].ID
		} else if barbaSvcID == "" {
			barbaSvcID = corteSvcID
		}
		if comboSvcID == "" && len(svcs) > 2 {
			comboSvcID = svcs[2].ID
		} else if comboSvcID == "" {
			comboSvcID = corteSvcID
		}
		if acabamentoSvcID == "" {
			acabamentoSvcID = corteSvcID
		}
		if selagemSvcID == "" {
			selagemSvcID = corteSvcID
		}
	}

	// 6. Resources
	resList, _ := repos.resource.ListByTenant(tenantID)
	if len(resList) == 0 {
		_, _ = repos.resource.Create(resourcedomain.Resource{
			TenantID: tenantID,
			Name:     "Cadeira Master 01 (Visagismo)",
			Type:     "chair",
			Capacity: 1,
			Active:   true,
			Notes:    "Bancada principal para cortes e finalizações",
		})
		_, _ = repos.resource.Create(resourcedomain.Resource{
			TenantID: tenantID,
			Name:     "Cadeira Barba & Spa 02",
			Type:     "chair",
			Capacity: 1,
			Active:   true,
			Notes:    "Cadeira reclinável equipada com vaporizador de ozônio para barboterapia",
		})
		_, _ = repos.resource.Create(resourcedomain.Resource{
			TenantID: tenantID,
			Name:     "Bancada Estilo & Química 03",
			Type:     "chair",
			Capacity: 1,
			Active:   true,
			Notes:    "Área com lavatório e iluminação para colorimetria e selagem",
		})
	}

	// 7. Customers (5 Realistic Clients)
	custs, _ := repos.customer.ListByTenant(tenantID, "")
	var guiID, felipeID, rodrigoID, brunoID, matheusID string
	if len(custs) == 0 {
		c1, _ := repos.customer.Create(customerdomain.Customer{
			TenantID: tenantID,
			Name:     "Guilherme Arantes",
			Phone:    "(11) 99123-4567",
			Email:    "guilherme.arantes@gmail.com",
			Birthday: "1992-04-15",
			Notes:    "Prefere degradê navalhado baixo (skin fade). Gosta de café puro sem açúcar.",
		})
		guiID = c1.ID

		c2, _ := repos.customer.Create(customerdomain.Customer{
			TenantID: tenantID,
			Name:     "Felipe Alcântara",
			Phone:    "(11) 98234-5678",
			Email:    "felipe.alcantara@hotmail.com",
			Birthday: "1988-08-22",
			Notes:    "Cliente semanal assíduo. Gosta de barboterapia caprichada com toalha bem quente.",
		})
		felipeID = c2.ID

		c3, _ := repos.customer.Create(customerdomain.Customer{
			TenantID: tenantID,
			Name:     "Rodrigo Viana",
			Phone:    "(11) 97345-6789",
			Email:    "rodrigo.viana@outlook.com",
			Birthday: "1995-12-03",
			Notes:    "Cliente VIP. Faz combo e compra pomada matte e óleo para barba frequentemente.",
		})
		rodrigoID = c3.ID

		c4, _ := repos.customer.Create(customerdomain.Customer{
			TenantID: tenantID,
			Name:     "Bruno Carvalho",
			Phone:    "(11) 96456-7890",
			Email:    "bruno.carvalho@gmail.com",
			Birthday: "1990-01-29",
			Notes:    "Barba desenhada e acabamento navalhado. Prefere horário de almoço ou fim de tarde.",
		})
		brunoID = c4.ID

		c5, _ := repos.customer.Create(customerdomain.Customer{
			TenantID: tenantID,
			Name:     "Matheus Silveira",
			Phone:    "(11) 95567-8901",
			Email:    "matheus.silveira@uol.com.br",
			Birthday: "2001-06-18",
			Notes:    "Estudante de Engenharia. Corte tesoura social clássico e hidratação periódica.",
		})
		matheusID = c5.ID
	} else {
		for _, c := range custs {
			if strings.Contains(c.Name, "Guilherme") {
				guiID = c.ID
			} else if strings.Contains(c.Name, "Felipe") {
				felipeID = c.ID
			} else if strings.Contains(c.Name, "Rodrigo") {
				rodrigoID = c.ID
			} else if strings.Contains(c.Name, "Bruno") {
				brunoID = c.ID
			} else if strings.Contains(c.Name, "Matheus") {
				matheusID = c.ID
			}
		}
		if guiID == "" && len(custs) > 0 {
			guiID = custs[0].ID
		}
		if felipeID == "" && len(custs) > 1 {
			felipeID = custs[1].ID
		} else if felipeID == "" {
			felipeID = guiID
		}
		if rodrigoID == "" && len(custs) > 2 {
			rodrigoID = custs[2].ID
		} else if rodrigoID == "" {
			rodrigoID = guiID
		}
		if brunoID == "" && len(custs) > 3 {
			brunoID = custs[3].ID
		} else if brunoID == "" {
			brunoID = guiID
		}
		if matheusID == "" && len(custs) > 4 {
			matheusID = custs[4].ID
		} else if matheusID == "" {
			matheusID = guiID
		}
	}

	// 8. Dynamic Appointments (Hoje, Amanhã, Dias Recentes)
	seedAppointments(repos.appointment, tenantID, carlosProID, rafaelProID, lucasProID, corteSvcID, barbaSvcID, comboSvcID, acabamentoSvcID, selagemSvcID, guiID, felipeID, rodrigoID, brunoID, matheusID)

	// 9. Stock / Products
	products, _ := repos.stock.ListProducts(tenantID, false)
	if len(products) == 0 {
		_, _ = repos.stock.CreateProduct(stockdomain.Product{
			TenantID:      tenantID,
			Name:          "Pomada Modeladora Efeito Matte (100g)",
			SKU:           "POM-MAT-01",
			Category:      "Finalizadores",
			PriceCents:    4500,
			CostCents:     2000,
			StockQuantity: 18,
			MinStockAlert: 5,
			Unit:          "unidade",
			Active:        true,
		})
		_, _ = repos.stock.CreateProduct(stockdomain.Product{
			TenantID:      tenantID,
			Name:          "Óleo Hidratante para Barba Premium (30ml)",
			SKU:           "OIL-BRB-02",
			Category:      "Barba",
			PriceCents:    3800,
			CostCents:     1500,
			StockQuantity: 12,
			MinStockAlert: 4,
			Unit:          "unidade",
			Active:        true,
		})
		_, _ = repos.stock.CreateProduct(stockdomain.Product{
			TenantID:      tenantID,
			Name:          "Shampoo Antiqueda & Fortificante (250ml)",
			SKU:           "SHA-FRT-03",
			Category:      "Capilar",
			PriceCents:    5500,
			CostCents:     2400,
			StockQuantity: 8,
			MinStockAlert: 3,
			Unit:          "unidade",
			Active:        true,
		})
		_, _ = repos.stock.CreateProduct(stockdomain.Product{
			TenantID:      tenantID,
			Name:          "Balm Pós-Barba Refrescante Mentol (120ml)",
			SKU:           "BLM-REF-04",
			Category:      "Barba",
			PriceCents:    3500,
			CostCents:     1400,
			StockQuantity: 14,
			MinStockAlert: 4,
			Unit:          "unidade",
			Active:        true,
		})
	}

	// 10. Cash Register
	if currentOpen, err := repos.cashRegister.GetCurrentOpen(tenantID); err != nil || currentOpen.ID == "" {
		todayStr := time.Now().Format("2006-01-02")
		openTime := todayStr + "T08:30:00Z"
		_, _ = repos.cashRegister.Create(cashregisterdomain.CashRegister{
			TenantID:            tenantID,
			OpenedByUserID:      ownerUserID,
			Status:              "open",
			InitialBalanceCents: 15000, // R$ 150,00 troco inicial
			OpenedAt:            openTime,
			Notes:               "Abertura de caixa padrão do turno da manhã",
		})
	}

	// 11. Finance
	txs, _ := repos.finance.List(tenantID, "", "", "", "")
	if len(txs) == 0 {
		todayStr := time.Now().Format("2006-01-02")
		yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		_, _ = repos.finance.Create(financedomain.Transaction{
			TenantID:      tenantID,
			Type:          "income",
			Category:      "service_sale",
			AmountCents:   5000,
			Description:   "Corte de Cabelo Masculino - Guilherme Arantes",
			PaymentMethod: "pix",
			DueDate:       todayStr,
			PaidAt:        todayStr + "T09:35:00Z",
			Status:        "paid",
		})
		_, _ = repos.finance.Create(financedomain.Transaction{
			TenantID:      tenantID,
			Type:          "income",
			Category:      "service_sale",
			AmountCents:   8500,
			Description:   "Combo Corte + Barba - Rodrigo Viana",
			PaymentMethod: "credit_card",
			DueDate:       todayStr,
			PaidAt:        todayStr + "T11:05:00Z",
			Status:        "paid",
		})
		_, _ = repos.finance.Create(financedomain.Transaction{
			TenantID:      tenantID,
			Type:          "income",
			Category:      "product_sale",
			AmountCents:   4500,
			Description:   "Venda 1x Pomada Modeladora Matte (Balcão)",
			PaymentMethod: "pix",
			DueDate:       todayStr,
			PaidAt:        todayStr + "T10:15:00Z",
			Status:        "paid",
		})
		_, _ = repos.finance.Create(financedomain.Transaction{
			TenantID:      tenantID,
			Type:          "expense",
			Category:      "supplies",
			AmountCents:   12000,
			Description:   "Reposição de lâminas descartáveis e toalhas higienizadas",
			PaymentMethod: "pix",
			DueDate:       yesterdayStr,
			PaidAt:        yesterdayStr + "T16:00:00Z",
			Status:        "paid",
		})
		_, _ = repos.finance.Create(financedomain.Transaction{
			TenantID:      tenantID,
			Type:          "expense",
			Category:      "rent",
			AmountCents:   220000,
			Description:   "Aluguel Comercial da Barbearia (Mês Atual)",
			PaymentMethod: "boleto",
			DueDate:       time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
			Status:        "pending",
		})
	}

	// 12. Reviews
	reviews, _ := repos.review.ListByTenant(tenantID, "")
	if len(reviews) == 0 {
		nowISO := time.Now().Format(time.RFC3339)
		_, _ = repos.review.Create(reviewdomain.Review{
			TenantID:     tenantID,
			TenantSlug:   "dom-navalha",
			CustomerID:   rodrigoID,
			CustomerName: "Rodrigo Viana",
			Status:       "submitted",
			Rating:       5,
			Comment:      "Melhor barbearia da região! O atendimento do Carlos e do Rafael é impecável, toalha quente sensacional e pontualidade exemplar.",
			InvitedAt:    nowISO,
			SubmittedAt:  nowISO,
			CreatedAt:    nowISO,
		})
		_, _ = repos.review.Create(reviewdomain.Review{
			TenantID:     tenantID,
			TenantSlug:   "dom-navalha",
			CustomerID:   guiID,
			CustomerName: "Guilherme Arantes",
			Status:       "submitted",
			Rating:       5,
			Comment:      "O degradê navalhado mais preciso que já fiz. Ambiente limpo, café de primeira e profissionais extremamente educados.",
			InvitedAt:    nowISO,
			SubmittedAt:  nowISO,
			CreatedAt:    nowISO,
		})
		_, _ = repos.review.Create(reviewdomain.Review{
			TenantID:     tenantID,
			TenantSlug:   "dom-navalha",
			CustomerID:   felipeID,
			CustomerName: "Felipe Alcântara",
			Status:       "submitted",
			Rating:       5,
			Comment:      "Atendimento nota 10! Sou cliente fiel há anos, a barboterapia renova qualquer um. Recomendo de olhos fechados!",
			InvitedAt:    nowISO,
			SubmittedAt:  nowISO,
			CreatedAt:    nowISO,
		})
	}

	log.Println(`{"level":"INFO","message":"homologation data seed completed successfully"}`)
}

func seedAccount(repo authdomain.Repository, hasher auth.PasswordHasher, name, email, phone, password, role, tenantID string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	if existing, err := repo.GetStaffAccountByEmail(email); err == nil {
		if existing.Role != role {
			_, _ = repo.UpdateRole(existing.ID, role)
		}
		if hasher.Compare(existing.PasswordHash, password) != nil {
			if hash, err := hasher.Hash(password); err == nil {
				_ = repo.UpdatePassword(existing.ID, hash)
			}
		}
		return existing.ID
	}

	hash, err := hasher.Hash(password)
	if err != nil {
		return ""
	}

	acc, err := repo.CreateAccount(authdomain.Account{
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: hash,
		Role:         role,
		TenantID:     tenantID,
	})
	if err != nil {
		return ""
	}
	return acc.ID
}

func seedAppointments(
	repo appointmentdomain.Repository,
	tenantID string,
	carlosID, rafaelID, lucasID string,
	corteID, barbaID, comboID, acabamentoID, selagemID string,
	guiID, felipeID, rodrigoID, brunoID, matheusID string,
) {
	// Look up existing appointments for tenant
	existing, _ := repo.ListByRange(tenantID, "", "", "", "", "")
	if len(existing) > 0 {
		return
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	makeTime := func(base time.Time, hour, minute int) string {
		t := time.Date(base.Year(), base.Month(), base.Day(), hour, minute, 0, 0, time.Local)
		return t.Format(time.RFC3339)
	}

	type apptItem struct {
		proID     string
		svcID     string
		custID    string
		startsAt  string
		endsAt    string
		status    string
		payStatus string
		notes     string
		checkInAt string
	}

	items := []apptItem{
		// --- HOJE (Today) ---
		{
			proID:     carlosID,
			svcID:     corteID,
			custID:    guiID,
			startsAt:  makeTime(today, 9, 0),
			endsAt:    makeTime(today, 9, 30),
			status:    "completed",
			payStatus: "paid",
			notes:     "Degradê navalhado skin fade",
		},
		{
			proID:     rafaelID,
			svcID:     comboID,
			custID:    rodrigoID,
			startsAt:  makeTime(today, 10, 0),
			endsAt:    makeTime(today, 11, 0),
			status:    "completed",
			payStatus: "paid",
			notes:     "Combo completo Dom Navalha com barboterapia",
		},
		{
			proID:     carlosID,
			svcID:     barbaID,
			custID:    brunoID,
			startsAt:  makeTime(today, 11, 30),
			endsAt:    makeTime(today, 12, 0),
			status:    "in_progress",
			payStatus: "pending",
			checkInAt: makeTime(today, 11, 25),
			notes:     "Cliente em atendimento na Cadeira 02",
		},
		{
			proID:     rafaelID,
			svcID:     corteID,
			custID:    felipeID,
			startsAt:  makeTime(today, 14, 0),
			endsAt:    makeTime(today, 14, 30),
			status:    "confirmed",
			payStatus: "pending",
			notes:     "Confirmou presença pelo WhatsApp",
		},
		{
			proID:     lucasID,
			svcID:     comboID,
			custID:    matheusID,
			startsAt:  makeTime(today, 15, 30),
			endsAt:    makeTime(today, 16, 30),
			status:    "confirmed",
			payStatus: "pending",
			notes:     "Horário agendado online",
		},
		{
			proID:     carlosID,
			svcID:     acabamentoID,
			custID:    guiID,
			startsAt:  makeTime(today, 17, 0),
			endsAt:    makeTime(today, 17, 20),
			status:    "scheduled",
			payStatus: "pending",
			notes:     "Acabamento pézinho e sobrancelha",
		},
		{
			proID:     lucasID,
			svcID:     selagemID,
			custID:    rodrigoID,
			startsAt:  makeTime(today, 18, 0),
			endsAt:    makeTime(today, 19, 0),
			status:    "scheduled",
			payStatus: "pending",
			notes:     "Selagem capilar anti-frizz",
		},

		// --- AMANHÃ (Tomorrow) ---
		{
			proID:     carlosID,
			svcID:     comboID,
			custID:    rodrigoID,
			startsAt:  makeTime(tomorrow, 9, 30),
			endsAt:    makeTime(tomorrow, 10, 30),
			status:    "confirmed",
			payStatus: "pending",
			notes:     "Agendado pelo portal de clientes",
		},
		{
			proID:     rafaelID,
			svcID:     barbaID,
			custID:    brunoID,
			startsAt:  makeTime(tomorrow, 11, 0),
			endsAt:    makeTime(tomorrow, 11, 30),
			status:    "scheduled",
			payStatus: "pending",
		},
		{
			proID:     lucasID,
			svcID:     selagemID,
			custID:    matheusID,
			startsAt:  makeTime(tomorrow, 14, 30),
			endsAt:    makeTime(tomorrow, 15, 30),
			status:    "scheduled",
			payStatus: "pending",
		},
		{
			proID:     carlosID,
			svcID:     corteID,
			custID:    felipeID,
			startsAt:  makeTime(tomorrow, 16, 0),
			endsAt:    makeTime(tomorrow, 16, 30),
			status:    "scheduled",
			payStatus: "pending",
		},

		// --- ONTEM / ANTERIORES (Recent Past) ---
		{
			proID:     rafaelID,
			svcID:     comboID,
			custID:    rodrigoID,
			startsAt:  makeTime(yesterday, 10, 0),
			endsAt:    makeTime(yesterday, 11, 0),
			status:    "completed",
			payStatus: "paid",
			notes:     "Atendimento concluído com sucesso",
		},
		{
			proID:     carlosID,
			svcID:     corteID,
			custID:    felipeID,
			startsAt:  makeTime(yesterday, 14, 0),
			endsAt:    makeTime(yesterday, 14, 30),
			status:    "completed",
			payStatus: "paid",
			notes:     "Corte social e lavagem",
		},
		{
			proID:     rafaelID,
			svcID:     barbaID,
			custID:    brunoID,
			startsAt:  makeTime(yesterday, 16, 0),
			endsAt:    makeTime(yesterday, 16, 30),
			status:    "completed",
			payStatus: "paid",
		},
		{
			proID:     lucasID,
			svcID:     comboID,
			custID:    matheusID,
			startsAt:  makeTime(yesterday, 17, 30),
			endsAt:    makeTime(yesterday, 18, 30),
			status:    "completed",
			payStatus: "paid",
		},
	}

	for _, item := range items {
		_, _ = repo.Create(appointmentdomain.Appointment{
			ID:             uuid.NewString(),
			TenantID:       tenantID,
			ProfessionalID: item.proID,
			ServiceID:      item.svcID,
			CustomerID:     item.custID,
			StartsAt:       item.startsAt,
			EndsAt:         item.endsAt,
			Status:         item.status,
			PaymentStatus:  item.payStatus,
			Notes:          item.notes,
			CheckInAt:      item.checkInAt,
			CreatedBy:      "system_seed",
			CreatedAt:      time.Now().Format(time.RFC3339),
			UpdatedAt:      time.Now().Format(time.RFC3339),
		})
	}
}

