package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	paymentdomain "erpagendamentos/backend/internal/domain/payment"
)

type MercadoPagoProvider struct {
	accessToken string
	webhookURL  string
	client      *http.Client
}

func NewMercadoPagoProvider(accessToken string, webhookURL string) *MercadoPagoProvider {
	return &MercadoPagoProvider{
		accessToken: strings.TrimSpace(accessToken),
		webhookURL:  strings.TrimSpace(webhookURL),
		client:      &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *MercadoPagoProvider) CreatePayment(input paymentdomain.IntentInput) (paymentdomain.IntentResult, error) {
	if p.accessToken == "" {
		return paymentdomain.IntentResult{}, fmt.Errorf("mercadopago access token is not configured")
	}

	amountBRL := float64(input.AmountCents) / 100.0
	email := strings.TrimSpace(input.PayerEmail)
	if email == "" {
		email = "cliente@agendamento.com"
	}
	name := strings.TrimSpace(input.PayerName)
	if name == "" {
		name = "Cliente"
	}
	desc := strings.TrimSpace(input.Description)
	if desc == "" {
		desc = "Agendamento de Serviço"
	}

	payload := map[string]any{
		"transaction_amount": amountBRL,
		"description":        desc,
		"payment_method_id":  "pix",
		"payer": map[string]any{
			"email":      email,
			"first_name": name,
		},
	}
	if p.webhookURL != "" {
		payload["notification_url"] = p.webhookURL
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return paymentdomain.IntentResult{}, fmt.Errorf("marshal mp payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://api.mercadopago.com/v1/payments", bytes.NewReader(bodyBytes))
	if err != nil {
		return paymentdomain.IntentResult{}, fmt.Errorf("build mp request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Idempotency-Key", fmt.Sprintf("mp_%d", time.Now().UnixNano()))

	resp, err := p.client.Do(req)
	if err != nil {
		return paymentdomain.IntentResult{}, fmt.Errorf("mercadopago request error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return paymentdomain.IntentResult{}, fmt.Errorf("mercadopago error %d: %s", resp.StatusCode, string(respBody))
	}

	var res struct {
		ID                 int64  `json:"id"`
		Status             string `json:"status"`
		StatusDetail       string `json:"status_detail"`
		PointOfInteraction struct {
			TransactionData struct {
				QRCode       string `json:"qr_code"`
				QRCodeBase64 string `json:"qr_code_base64"`
			} `json:"transaction_data"`
		} `json:"point_of_interaction"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return paymentdomain.IntentResult{}, fmt.Errorf("decode mp response: %w", err)
	}

	return paymentdomain.IntentResult{
		ProviderPaymentID: fmt.Sprintf("%d", res.ID),
		Status:            res.Status,
		PixCopiaECola:     res.PointOfInteraction.TransactionData.QRCode,
		PixQRCodeBase64:   res.PointOfInteraction.TransactionData.QRCodeBase64,
		ProviderData:      string(respBody),
	}, nil
}

func (p *MercadoPagoProvider) GetPaymentStatus(providerPaymentID string) (string, error) {
	if p.accessToken == "" {
		return "", fmt.Errorf("mercadopago access token is not configured")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://api.mercadopago.com/v1/payments/"+providerPaymentID, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+p.accessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("mercadopago get status error %d", resp.StatusCode)
	}

	var res struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	return res.Status, nil
}

func (p *MercadoPagoProvider) Refund(providerPaymentID string) error {
	if p.accessToken == "" {
		return fmt.Errorf("mercadopago access token is not configured")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf("https://api.mercadopago.com/v1/payments/%s/refunds", providerPaymentID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.accessToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("mercadopago refund error %d", resp.StatusCode)
	}
	return nil
}

