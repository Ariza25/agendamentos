package rest

import (
	"encoding/json"
	"net/http"

	"erpagendamentos/backend/internal/application/onboarding"
)

type OnboardingHandler struct {
	onboardingUseCase *onboarding.UseCase
}

func NewOnboardingHandler(onboardingUseCase *onboarding.UseCase) *OnboardingHandler {
	return &OnboardingHandler{onboardingUseCase: onboardingUseCase}
}

func (h *OnboardingHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var input onboarding.SignupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "invalid signup payload")
		return
	}

	session, tenant, err := h.onboardingUseCase.Signup(input)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "signup_failed", err.Error())
		return
	}

	csrfToken := setSessionCookies(w, session.Token)
	writeJSON(w, http.StatusCreated, map[string]any{
		"session":   session,
		"tenant":    tenant,
		"csrfToken": csrfToken,
	})
}
