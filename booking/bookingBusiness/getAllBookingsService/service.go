package getAllBookingsService

import (
	"context"

	"github.com/iondodon/go-vbs/booking/bookingBusiness"
	"github.com/iondodon/go-vbs/booking/bookingDomain"
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
