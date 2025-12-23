package models

// Organization represents an organization
type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt Time   `json:"created_at"`
	UpdatedAt Time   `json:"updated_at"`
}

// OrganizationMember represents a member of an organization
type OrganizationMember struct {
	User      User     `json:"user"`
	Role      UserRole `json:"role"`
	JoinedAt  Time     `json:"joined_at"`
	InvitedBy string   `json:"invited_by,omitempty"`
}

// InviteMemberRequest represents a request to invite a member
type InviteMemberRequest struct {
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

// OrganizationListResponse represents a list of organizations
type OrganizationListResponse struct {
	Organizations []Organization `json:"organizations"`
	Total         int            `json:"total"`
}

// OrganizationMembersResponse represents a list of organization members
type OrganizationMembersResponse struct {
	Members []OrganizationMember `json:"members"`
	Total   int                  `json:"total"`
}
