package getAllBookings

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/usecase"
)

type GetAllBookingsInterface interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}

type GetAllBookings struct {
	bookingRepo usecase.BookingRepositoryInterface
}

func New(bookingRepo usecase.BookingRepositoryInterface) GetAllBookingsInterface {
	return &GetAllBookings{
		bookingRepo: bookingRepo,
	}
}

func (uc *GetAllBookings) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
