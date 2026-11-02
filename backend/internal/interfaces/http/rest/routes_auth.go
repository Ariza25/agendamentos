package rest

import "net/http"

func registerAuthRoutes(mux *http.ServeMux, config RouterConfig) {
	authHandler := NewAuthHandler(config.AuthUseCase)

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /live", health)
	mux.HandleFunc("GET /ready", readiness(config.ReadinessChecker))

	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/v1/auth/password-reset/request", authHandler.RequestPasswordReset)
	mux.HandleFunc("POST /api/v1/auth/password-reset/confirm", authHandler.ConfirmPasswordReset)
	mux.HandleFunc("POST /api/v1/auth/password-reset/code/request", authHandler.RequestPasswordResetCode)
	mux.HandleFunc("POST /api/v1/auth/password-reset/code/confirm", authHandler.ConfirmPasswordResetCode)

	mux.HandleFunc("GET /api/v1/auth/me", authHandler.Me)
	mux.HandleFunc("GET /api/v1/auth/sessions", authHandler.ListSessions)
	mux.HandleFunc("DELETE /api/v1/auth/sessions/{sessionID}", authHandler.RevokeSession)
	mux.HandleFunc("POST /api/v1/auth/sessions/revoke-others", authHandler.RevokeOtherSessions)
	mux.HandleFunc("PUT /api/v1/auth/profile", authHandler.UpdateProfile)
	mux.HandleFunc("PUT /api/v1/auth/password", authHandler.ChangePassword)

	mux.HandleFunc("GET /api/v1/auth/mfa", requirePrivilegedAdmin(config.AuthUseCase, authHandler.MFAStatus))
	mux.HandleFunc("POST /api/v1/auth/mfa/setup", requirePrivilegedAdmin(config.AuthUseCase, authHandler.BeginMFASetup))
	mux.HandleFunc("POST /api/v1/auth/mfa/enable", requirePrivilegedAdmin(config.AuthUseCase, authHandler.EnableMFA))
	mux.HandleFunc("DELETE /api/v1/auth/mfa", requirePrivilegedAdmin(config.AuthUseCase, authHandler.DisableMFA))
	mux.HandleFunc("POST /api/v1/auth/mfa/step-up", requirePrivilegedAdmin(config.AuthUseCase, authHandler.StepUpMFA))

	mux.HandleFunc("GET /api/v1/admin/users", requireOwnerOrAdmin(config.AuthUseCase, authHandler.ListAdminUsers))
	mux.HandleFunc("POST /api/v1/admin/users", requireOwnerOrAdmin(config.AuthUseCase, authHandler.CreateAdminUser))
	mux.HandleFunc("PATCH /api/v1/admin/users/{userID}/role", requireOwnerOrAdmin(config.AuthUseCase, authHandler.UpdateAdminUserRole))
	mux.HandleFunc("DELETE /api/v1/admin/users/{userID}", requireOwnerOrAdmin(config.AuthUseCase, authHandler.DeleteAdminUser))
}
