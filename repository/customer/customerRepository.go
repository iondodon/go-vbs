package customer

import (
	uuidlib "github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/integration"
)

const getCustomerByUUID = `
	SELECT c.id, c.uuid, c.username
	FROM customer c
	WHERE c.uuid = ?
`

type CustomerRepository interface {
	FindByUUID(cUUID uuidlib.UUID) (*domain.Customer, error)
}

type customerRepository struct {
	db integration.DB
}

func NewCustomerRepository(db integration.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (repo *customerRepository) FindByUUID(cUUID uuidlib.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	var uuidStr string

	err := repo.db.QueryRow(getCustomerByUUID, cUUID.String()).Scan(
		&customer.ID,
		&uuidStr,
		&customer.Username,
	)
	if err != nil {
		return nil, err
	}

	uuid, err := uuidlib.Parse(uuidStr)
	if err != nil {
		return nil, err
	}

	customer.UUID = uuid

	return &customer, nil
}
