package models

type BillingPlan string

const (
	BillingPlanFree       BillingPlan = "free"
	BillingPlanDeveloper  BillingPlan = "developer"
	BillingPlanPro        BillingPlan = "pro"
	BillingPlanEnterprise BillingPlan = "enterprise"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceStatusPending InvoiceStatus = "pending"
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusFailed  InvoiceStatus = "failed"
)

type UsageSummary struct {
	OrganizationID string                 `json:"organization_id"`
	Period         string                 `json:"period"`
	TotalCost      float64                `json:"total_cost"`
	ServiceCount   int                    `json:"service_count"`
	Services       []ServiceUsage         `json:"services,omitempty"`
	Breakdown      map[string]interface{} `json:"breakdown,omitempty"`
	UpdatedAt      Time                   `json:"updated_at"`
}

// ServiceUsage represents usage for a single service
type ServiceUsage struct {
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	ServiceType string  `json:"service_type"`
	Tier        string  `json:"tier"`
	Cost        float64 `json:"cost"`
	Uptime      float64 `json:"uptime_hours"`
}

type Invoice struct {
	ID             string        `json:"id"`
	OrganizationID string        `json:"organization_id"`
	Amount         float64       `json:"amount"`
	Currency       string        `json:"currency"`
	Status         InvoiceStatus `json:"status"`
	Period         string        `json:"period"`
	DueDate        Time          `json:"due_date"`
	PaidAt         *Time         `json:"paid_at,omitempty"`
	CreatedAt      Time          `json:"created_at"`
	DownloadURL    string        `json:"download_url,omitempty"`
}

// InvoiceListResponse represents a list of invoices
type InvoiceListResponse struct {
	Invoices []Invoice `json:"invoices"`
	Total    int       `json:"total"`
}

type PlanInfo struct {
	Plan          BillingPlan `json:"plan"`
	DisplayName   string      `json:"display_name"`
	PriceMonthly  float64     `json:"price_monthly"`
	ServiceLimit  int         `json:"service_limit"`
	Features      []string    `json:"features"`
	CurrentPlan   bool        `json:"current_plan"`
	NextBillingAt *Time       `json:"next_billing_at,omitempty"`
}

// UpgradePlanRequest represents a request to upgrade plan
type UpgradePlanRequest struct {
	Plan BillingPlan `json:"plan"`
}
