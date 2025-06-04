package getAllBookings

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
)

type UseCase interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}
