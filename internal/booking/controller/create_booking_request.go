package controller

import (
	"github.com/google/uuid"
)

type CreateBookingRequestDTO struct {
	CustomerUUID uuid.UUID     `json:"customer_uuid"`
	VehicleUUID  uuid.UUID     `json:"vehicle_uuid"`
	DatePeriodD  DatePeriodDTO `json:"date_period"`
}
