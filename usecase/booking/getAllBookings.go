package booking

import (
	"github.com/iondodon/go-vbs/domain"

	bookingRepos "github.com/iondodon/go-vbs/repository/booking"
)

type GetAllBookings interface {
	Execute() ([]domain.Booking, error)
}

type getAllBookings struct {
	bookingRepo bookingRepos.BookingRepository
}

func NewGetAllBookings(bookingRepo bookingRepos.BookingRepository) GetAllBookings {
	return &getAllBookings{bookingRepo: bookingRepo}
}

func (uc *getAllBookings) Execute() ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll()
}
