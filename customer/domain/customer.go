package customerDomain

import "github.com/google/uuid"

type Customer struct {
	ID       int64
	UUID     uuid.UUID
	Username string
}
