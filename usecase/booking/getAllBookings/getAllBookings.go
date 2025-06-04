package getAllBookings

import (
	"context"

	"github.com/iondodon/go-vbs/domain"

	"github.com/iondodon/go-vbs/repository/bookingRepository"
)

type GetAllBookingsInterface interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}

type GetAllBookings struct {
	bookingRepo bookingRepository.BookingRepositoryInterface
}

func New(bookingRepo bookingRepository.BookingRepositoryInterface) GetAllBookingsInterface {
	return &GetAllBookings{
		bookingRepo: bookingRepo,
	}
}

func (uc *GetAllBookings) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
