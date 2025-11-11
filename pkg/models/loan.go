package models

import "time"

type LoanApplicationStatus string

const (
	StatusNew      LoanApplicationStatus = "NEW"
	StatusReview   LoanApplicationStatus = "REVIEW"
	StatusApproved LoanApplicationStatus = "APPROVED"
	StatusRejected LoanApplicationStatus = "Rejected"
)

type LoanApplicationType string

const (
	TypeAuto     LoanApplicationType = "AUTO"
	TypePersonal LoanApplicationType = "PERSONAL"
)

type LoanApplication struct {
	ID             int64
	UserID         int64
	Type           LoanApplicationType
	VehicleVIN     *string
	VehicleName    *string
	Price          float64
	DownPayment    float64
	NetPrice       float64
	MarginRate     float64
	TermMonths     int32
	MonthlyPayment float64
	Status         LoanApplicationStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
