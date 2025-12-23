package models

type RecommendationType string

const (
	RecommendationTypeTier          RecommendationType = "tier"
	RecommendationTypeConfiguration RecommendationType = "configuration"
	RecommendationTypeScaling       RecommendationType = "scaling"
	RecommendationTypeCost          RecommendationType = "cost"
)

// RecommendationPriority represents the priority level
type RecommendationPriority string

const (
	RecommendationPriorityLow      RecommendationPriority = "low"
	RecommendationPriorityMedium   RecommendationPriority = "medium"
	RecommendationPriorityHigh     RecommendationPriority = "high"
	RecommendationPriorityCritical RecommendationPriority = "critical"
)

type Recommendation struct {
	ID                string                 `json:"id"`
	ServiceID         string                 `json:"service_id,omitempty"`
	Type              RecommendationType     `json:"type"`
	Priority          RecommendationPriority `json:"priority"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	CurrentConfig     map[string]interface{} `json:"current_config,omitempty"`
	RecommendedConfig map[string]interface{} `json:"recommended_config,omitempty"`
	EstimatedSavings  float64                `json:"estimated_savings,omitempty"`
	EstimatedImpact   string                 `json:"estimated_impact,omitempty"`
	CreatedAt         Time                   `json:"created_at"`
}

// RecommendationRequest represents a request for recommendations
type RecommendationRequest struct {
	Workload   string `json:"workload,omitempty"`
	ServiceID  string `json:"service_id,omitempty"`
	FocusArea  string `json:"focus_area,omitempty"`  // cost, performance, reliability
	TimeWindow string `json:"time_window,omitempty"` // 1h, 24h, 7d, 30d
}

type RecommendationResponse struct {
	Recommendations []Recommendation `json:"recommendations"`
	Total           int              `json:"total"`
	GeneratedAt     Time             `json:"generated_at"`
}

type AnalysisResult struct {
	ServiceID   string                 `json:"service_id,omitempty"`
	Summary     string                 `json:"summary"`
	Health      string                 `json:"health"` // healthy, warning, critical
	Metrics     map[string]interface{} `json:"metrics"`
	Issues      []AnalysisIssue        `json:"issues,omitempty"`
	Suggestions []string               `json:"suggestions,omitempty"`
	AnalyzedAt  Time                   `json:"analyzed_at"`
}

// AnalysisIssue represents an identified issue
type AnalysisIssue struct {
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// OptimizationSuggestion represents an optimization suggestion
type OptimizationSuggestion struct {
	ID               string                 `json:"id"`
	Category         string                 `json:"category"` // cost, performance, security
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Impact           string                 `json:"impact"` // low, medium, high
	Effort           string                 `json:"effort"` // low, medium, high
	EstimatedBenefit map[string]interface{} `json:"estimated_benefit,omitempty"`
	Steps            []string               `json:"steps,omitempty"`
}

type OptimizationResponse struct {
	Suggestions []OptimizationSuggestion `json:"suggestions"`
	Total       int                      `json:"total"`
	Focus       string                   `json:"focus,omitempty"`
	GeneratedAt Time                     `json:"generated_at"`
}

type Anomaly struct {
	ID          string                 `json:"id"`
	ServiceID   string                 `json:"service_id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Description string                 `json:"description"`
	DetectedAt  Time                   `json:"detected_at"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
	Suggestion  string                 `json:"suggestion,omitempty"`
	Resolved    bool                   `json:"resolved"`
}

// AnomalyListResponse represents a list of anomalies
type AnomalyListResponse struct {
	Anomalies []Anomaly `json:"anomalies"`
	Total     int       `json:"total"`
}
