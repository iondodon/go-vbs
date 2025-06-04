package in

import "time"

type DatePeriodDTO struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}
