package api

import (
	"context"
	"fmt"

	"github.com/quickspin/quickspin-cli/internal/models"
)

// GetUsageSummary retrieves usage summary for an organization
func (c *Client) GetUsageSummary(ctx context.Context, orgID string) (*models.UsageSummary, error) {
	var result models.UsageSummary
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/usage", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUsageSummaryByPeriod retrieves usage summary for a specific period
func (c *Client) GetUsageSummaryByPeriod(ctx context.Context, orgID, period string) (*models.UsageSummary, error) {
	var result models.UsageSummary
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/usage?period=%s", orgID, period)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListInvoices retrieves all invoices for an organization
func (c *Client) ListInvoices(ctx context.Context, orgID string) ([]models.Invoice, error) {
	var result []models.Invoice
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/invoices", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetInvoice retrieves a specific invoice by ID
func (c *Client) GetInvoice(ctx context.Context, orgID, invoiceID string) (*models.Invoice, error) {
	var result models.Invoice
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/invoices/%s", orgID, invoiceID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DownloadInvoice downloads an invoice PDF
func (c *Client) DownloadInvoice(ctx context.Context, orgID, invoiceID string) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/invoices/%s/download", orgID, invoiceID)

	resp, err := c.httpClient.R().
		SetContext(ctx).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("failed to download invoice: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to download invoice: %s", resp.Status())
	}

	return resp.Body(), nil
}

// GetCurrentPlan retrieves the current billing plan
func (c *Client) GetCurrentPlan(ctx context.Context, orgID string) (*models.PlanInfo, error) {
	var result models.PlanInfo
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/plan", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListAvailablePlans retrieves all available billing plans
func (c *Client) ListAvailablePlans(ctx context.Context) ([]models.PlanInfo, error) {
	var result []models.PlanInfo
	if err := c.Get(ctx, "/api/v1/billing/plans", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpgradePlan upgrades to a different billing plan
func (c *Client) UpgradePlan(ctx context.Context, orgID string, plan models.BillingPlan) error {
	req := models.UpgradePlanRequest{Plan: plan}
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/plan", orgID)
	return c.Post(ctx, path, req, nil)
}

// CancelSubscription cancels the current subscription
func (c *Client) CancelSubscription(ctx context.Context, orgID string) error {
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/subscription", orgID)
	return c.Delete(ctx, path, nil)
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Last4     string `json:"last4,omitempty"`
	Brand     string `json:"brand,omitempty"`
	ExpiryMonth int  `json:"expiry_month,omitempty"`
	ExpiryYear  int  `json:"expiry_year,omitempty"`
	IsDefault bool   `json:"is_default"`
	CreatedAt models.Time `json:"created_at"`
}

// ListPaymentMethods retrieves all payment methods
func (c *Client) ListPaymentMethods(ctx context.Context, orgID string) ([]PaymentMethod, error) {
	var result []PaymentMethod
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/payment-methods", orgID)
	if err := c.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AddPaymentMethodRequest represents the request to add a payment method
type AddPaymentMethodRequest struct {
	Token string `json:"token"`
	SetAsDefault bool `json:"set_as_default"`
}

// AddPaymentMethod adds a new payment method
func (c *Client) AddPaymentMethod(ctx context.Context, orgID string, req AddPaymentMethodRequest) (*PaymentMethod, error) {
	var result PaymentMethod
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/payment-methods", orgID)
	if err := c.Post(ctx, path, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeletePaymentMethod deletes a payment method
func (c *Client) DeletePaymentMethod(ctx context.Context, orgID, methodID string) error {
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/payment-methods/%s", orgID, methodID)
	return c.Delete(ctx, path, nil)
}

// SetDefaultPaymentMethod sets a payment method as default
func (c *Client) SetDefaultPaymentMethod(ctx context.Context, orgID, methodID string) error {
	path := fmt.Sprintf("/api/v1/organizations/%s/billing/payment-methods/%s/default", orgID, methodID)
	return c.Post(ctx, path, nil, nil)
}
