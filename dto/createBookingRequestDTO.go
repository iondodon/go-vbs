package dto

import "github.com/google/uuid"

type CreateBookingRequestDTO struct {
	VehicleUUID  uuid.UUID     `json:"vehicle_uuid"`
	CustomerUUID uuid.UUID     `json:"customer_uuid"`
	DatePeriodD  DatePeriodDTO `json:"date_period"`
}
