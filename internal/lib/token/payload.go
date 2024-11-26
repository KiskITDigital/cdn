package token

type Payload struct {
	UserID         int `json:"user_id"`
	OrganizationID int `json:"organization_id"`
	Role           int `json:"role"`
}
