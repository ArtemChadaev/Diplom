package dto

// AssignRoleRequest — body for PATCH /admin/users/{id}/role
type AssignRoleRequest struct {
	Role string `json:"role"`
}

// SetBlockedRequest — body for PATCH /admin/users/{id}/blocked
type SetBlockedRequest struct {
	Blocked bool `json:"blocked"`
}

// UserResponse — public user fields returned by the API
type UserResponse struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	NsPvAccess  bool   `json:"ns_pv_access"`
	UkepBound   bool   `json:"ukep_bound"`
	IsBlocked   bool   `json:"is_blocked"`
}
