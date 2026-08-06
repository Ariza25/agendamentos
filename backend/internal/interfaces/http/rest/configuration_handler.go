package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	appconfiguration "erpagendamentos/backend/internal/application/configuration"
	configurationdomain "erpagendamentos/backend/internal/domain/configuration"
)

type ConfigurationHandler struct {
	useCase *appconfiguration.UseCase
}

func NewConfigurationHandler(useCase *appconfiguration.UseCase) *ConfigurationHandler {
	return &ConfigurationHandler{useCase: useCase}
}

func (h *ConfigurationHandler) Get(w http.ResponseWriter, r *http.Request) {
	configuration, err := h.useCase.Get(resolveRequestTenantID(r, nil))
	if err != nil {
		writeInternalError(w, "Could not load the configuration.")
		return
	}
	writeJSON(w, http.StatusOK, configuration)
}

func (h *ConfigurationHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input configurationdomain.Configuration
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeInvalidPayload(w, "Could not read the configuration data.")
		return
	}
	input.TenantID = resolveRequestTenantID(r, nil)
	configuration, err := h.useCase.Update(input)
	if err != nil {
		writeConfigurationError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, configuration)
}

func writeConfigurationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, appconfiguration.ErrInvalidConfiguration):
		writeBadRequestError(w, "invalid_input", "Invalid configuration data.")
	default:
		writeInternalError(w, "Could not process the configuration.")
	}
}
