package getAllBookingsService

import (
	"context"

	"github.com/iondodon/go-vbs/booking/business"
	"github.com/iondodon/go-vbs/domain"
)

// Service handles retrieving all bookings
type Service struct {
	bookingRepo business.BookingRepository
}

func New(bookingRepo business.BookingRepository) *Service {
	return &Service{
		bookingRepo: bookingRepo,
	}
}

func (s *Service) Execute(ctx context.Context) ([]domain.Booking, error) {
	return s.bookingRepo.GetAll(ctx)
}
