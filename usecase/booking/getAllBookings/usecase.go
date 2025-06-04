package getAllBookings

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/usecase"
)

type UseCase interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}

type Service struct {
	bookingRepo usecase.BookingRepositoryInterface
}

func New(bookingRepo usecase.BookingRepositoryInterface) UseCase {
	return &Service{
		bookingRepo: bookingRepo,
	}
}

func (uc *Service) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
