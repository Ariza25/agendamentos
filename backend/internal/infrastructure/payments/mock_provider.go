package payments

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"

	paymentdomain "erpagendamentos/backend/internal/domain/payment"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) CreatePayment(input paymentdomain.IntentInput) (paymentdomain.IntentResult, error) {
	providerID := "mock_" + uuid.NewString()
	providerData := ""
	if pixKey := strings.TrimSpace(os.Getenv("PIX_KEY")); pixKey != "" {
		providerData = fmt.Sprintf("pix_key:%s", pixKey)
	}
	return paymentdomain.IntentResult{
		ProviderPaymentID: providerID,
		Status:            string(paymentdomain.StatusPaid),
		ProviderData:      providerData,
	}, nil
}

func (p *MockProvider) GetPaymentStatus(providerPaymentID string) (string, error) {
	return string(paymentdomain.StatusPaid), nil
}

func (p *MockProvider) Refund(providerPaymentID string) error {
	return nil
}
