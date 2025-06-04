package getAllBookings

import (
	"context"

	"github.com/iondodon/go-vbs/domain"

	"github.com/iondodon/go-vbs/repository/bookingRepository"
)

type GetAllBookings struct {
	bookingRepo bookingRepository.BookingRepository
}

func (uc *GetAllBookings) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
