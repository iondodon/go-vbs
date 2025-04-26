package booking

import (
	"context"

	"github.com/iondodon/go-vbs/domain"

	bookingRepos "github.com/iondodon/go-vbs/repository/booking"
)

//gobok:builder
type GetAllBookings struct {
	bookingRepo bookingRepos.BookingRepository
}

func (uc *GetAllBookings) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
