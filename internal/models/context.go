package models

type ctxKey int

const (
	UserIDKey ctxKey = iota
	OrganizationIDKey
	RoleKey
)
