package customerRepository

import (
	"context"

	uuidlib "github.com/google/uuid"
	"github.com/iondodon/go-vbs/customer/customerBusiness"
	"github.com/iondodon/go-vbs/customer/customerDomain"
	"github.com/iondodon/go-vbs/repository"
)

type Repository struct {
	queries *repository.Queries
}

// Ensure Repository implements the business interface
var _ customerBusiness.CustomerRepository = (*Repository)(nil)

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

	return &customer, nil
}
