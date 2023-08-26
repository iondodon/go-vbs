package domain

import "github.com/google/uuid"

type Vehicle struct {
	ID                 int
	UUID               uuid.UUID
	RegistrationNumber string
}
