package bookVehicle

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/dto"
)

type UseCase interface {
	ForPeriod(ctx context.Context, tx *sql.Tx, customerUID, vehicleUUID uuid.UUID, period dto.DatePeriodDTO) error
}
