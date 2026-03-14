package getAllBookingsService

import (
	"context"

	bookingBusiness "github.com/iondodon/go-vbs/internal/booking/business"
	bookingDomain "github.com/iondodon/go-vbs/internal/booking/domain"
)

// Service handles retrieving all bookings
type Service struct {
	bookingRepo bookingBusiness.BookingRepository
}

func New(bookingRepo bookingBusiness.BookingRepository) *Service {
	return &Service{
		bookingRepo: bookingRepo,
	}
}

func (s *Service) Execute(ctx context.Context) ([]bookingDomain.Booking, error) {
	return s.bookingRepo.GetAll(ctx)
}
