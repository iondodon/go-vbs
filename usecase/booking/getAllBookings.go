package booking

import (
	"context"

	"github.com/iondodon/go-vbs/domain"

	bookingRepos "github.com/iondodon/go-vbs/repository/booking"
)

type GetAllBookings interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}

//gobok:constructor
type getAllBookings struct {
	bookingRepo bookingRepos.BookingRepository
}

func (uc *getAllBookings) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
