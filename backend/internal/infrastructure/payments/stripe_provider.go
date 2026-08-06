package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	paymentdomain "erpagendamentos/backend/internal/domain/payment"
)

type StripeProvider struct {
	secretKey string
	client    *http.Client
}

func NewStripeProvider(secretKey string) *StripeProvider {
	return &StripeProvider{
		secretKey: strings.TrimSpace(secretKey),
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *StripeProvider) CreatePayment(input paymentdomain.IntentInput) (paymentdomain.IntentResult, error) {
	if p.secretKey == "" {
		return paymentdomain.IntentResult{}, fmt.Errorf("stripe secret key is not configured")
	}
	form := url.Values{}
	form.Set("amount", fmt.Sprintf("%d", input.AmountCents))
	form.Set("currency", "brl")
	form.Set("automatic_payment_methods[enabled]", "true")
	form.Set("description", input.Description)
	form.Set("capture_method", "manual")

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://api.stripe.com/v1/payment_intents", strings.NewReader(form.Encode()))
	if err != nil {
		return paymentdomain.IntentResult{}, fmt.Errorf("build stripe request: %w", err)
	}
	request.SetBasicAuth(p.secretKey, "")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := p.client.Do(request)
	if err != nil {
		return paymentdomain.IntentResult{}, fmt.Errorf("stripe create payment intent: %w", err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode >= 300 {
		return paymentdomain.IntentResult{}, fmt.Errorf("stripe error %d: %s", response.StatusCode, string(body))
	}
	var payload struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return paymentdomain.IntentResult{}, fmt.Errorf("decode stripe response: %w", err)
	}
	return paymentdomain.IntentResult{
		ProviderPaymentID: payload.ID,
		Status:            mapStripeStatus(payload.Status),
	}, nil
}

func (p *StripeProvider) GetPaymentStatus(providerPaymentID string) (string, error) {
	if p.secretKey == "" {
		return "", fmt.Errorf("stripe secret key is not configured")
	}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.stripe.com/v1/payment_intents/"+url.PathEscape(providerPaymentID), nil)
	if err != nil {
		return "", fmt.Errorf("build stripe request: %w", err)
	}
	request.SetBasicAuth(p.secretKey, "")
	response, err := p.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("stripe get payment intent: %w", err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	if response.StatusCode >= 300 {
		return "", fmt.Errorf("stripe error %d: %s", response.StatusCode, string(body))
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("decode stripe response: %w", err)
	}
	return mapStripeStatus(payload.Status), nil
}

func (p *StripeProvider) Refund(providerPaymentID string) error {
	if p.secretKey == "" {
		return fmt.Errorf("stripe secret key is not configured")
	}
	form := url.Values{}
	form.Set("payment_intent", providerPaymentID)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://api.stripe.com/v1/refunds", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("build stripe request: %w", err)
	}
	request.SetBasicAuth(p.secretKey, "")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := p.client.Do(request)
	if err != nil {
		return fmt.Errorf("stripe refund: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("stripe refund error %d: %s", response.StatusCode, string(body))
	}
	return nil
}

func mapStripeStatus(status string) string {
	switch status {
	case "succeeded":
		return string(paymentdomain.StatusPaid)
	case "processing", "requires_payment_method", "requires_confirmation", "requires_action", "requires_capture", "canceled":
		return string(paymentdomain.StatusPending)
	default:
		return string(paymentdomain.StatusPending)
	}
}
