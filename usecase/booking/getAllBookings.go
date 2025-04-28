package booking

import (
	"context"

	"github.com/iondodon/go-vbs/domain"

	bookingRepos "github.com/iondodon/go-vbs/repository/booking"
)

type GetAllBookingsInterface interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}

//gobok:constructor
//ctxboot:component
type GetAllBookings struct {
	bookingRepo bookingRepos.BookingRepositoryInterface `ctxboot:"inject"`
}

func (uc *GetAllBookings) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
