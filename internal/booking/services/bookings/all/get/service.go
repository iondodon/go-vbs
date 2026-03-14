package getAllBookingsService

import (
	"context"

	bookingDomain "github.com/iondodon/go-vbs/internal/booking/domain"
	bookingServices "github.com/iondodon/go-vbs/internal/booking/services"
)

// Service handles retrieving all bookings
type Service struct {
	bookingRepo bookingServices.BookingRepository
}

func New(bookingRepo bookingServices.BookingRepository) *Service {
	return &Service{
		bookingRepo: bookingRepo,
	}
}

func (s *Service) Execute(ctx context.Context) ([]bookingDomain.Booking, error) {
	return s.bookingRepo.GetAll(ctx)
}
