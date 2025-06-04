package getAllBookings

import (
	"context"

	"github.com/iondodon/go-vbs/business"
	"github.com/iondodon/go-vbs/domain"
)

type Service struct {
	bookingRepo business.BookingRepository
}

func New(bookingRepo business.BookingRepository) *Service {
	return &Service{
		bookingRepo: bookingRepo,
	}
}

func (uc *Service) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
