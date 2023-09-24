package booking

import (
	"github.com/iondodon/go-vbs/domain"

	bookingRepos "github.com/iondodon/go-vbs/repository/booking"
)

type GetAllBookingsUseCase interface {
	Execute() ([]domain.Booking, error)
}

type getAllBookingsUseCase struct {
	bookingRepo bookingRepos.BookingRepository
}

func NewGetAllBookingsUseCase(bookingRepo bookingRepos.BookingRepository) GetAllBookingsUseCase {
	return &getAllBookingsUseCase{bookingRepo: bookingRepo}
}

func (uc *getAllBookingsUseCase) Execute() ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll()
}
