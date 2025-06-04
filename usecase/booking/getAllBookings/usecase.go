package getAllBookings

import (
	"context"

	"github.com/iondodon/go-vbs/domain"
	"github.com/iondodon/go-vbs/usecase"
)

type Interface interface {
	Execute(ctx context.Context) ([]domain.Booking, error)
}

type UseCase struct {
	bookingRepo usecase.BookingRepositoryInterface
}

func New(bookingRepo usecase.BookingRepositoryInterface) Interface {
	return &UseCase{
		bookingRepo: bookingRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context) ([]domain.Booking, error) {
	return uc.bookingRepo.GetAll(ctx)
}
