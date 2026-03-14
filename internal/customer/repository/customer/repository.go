package customerRepository

import (
	"context"

	uuidlib "github.com/google/uuid"
	customerDomain "github.com/iondodon/go-vbs/internal/customer/domain"
	customerServices "github.com/iondodon/go-vbs/internal/customer/services"
	"github.com/iondodon/go-vbs/internal/repository"
)

type Repository struct {
	queries *repository.Queries
}

// Ensure Repository implements the services interface
var _ customerServices.CustomerRepository = (*Repository)(nil)

func New(queries *repository.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) FindByUUID(ctx context.Context, cUUID uuidlib.UUID) (*customerDomain.Customer, error) {
	var customer customerDomain.Customer

	customerRow, err := r.queries.GetCustomerByUUID(ctx, cUUID)
	if err != nil {
		return nil, err
	}

	customer.ID = customerRow.ID.(int64)

	customerUuidStr := customerRow.Uuid.(string)
	uuid, err := uuidlib.Parse(customerUuidStr)
	if err != nil {
		return nil, err
	}
	customer.UUID = uuid
	customer.Username = customerRow.Username

	return &customer, nil
}
