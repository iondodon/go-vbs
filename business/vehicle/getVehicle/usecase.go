package getVehicle

import (
	"context"

	"github.com/google/uuid"
	"github.com/iondodon/go-vbs/domain"
)

type UseCase interface {
	ByUUID(ctx context.Context, vUUID uuid.UUID) (*domain.Vehicle, error)
}
