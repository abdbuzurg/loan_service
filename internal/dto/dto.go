package dto

import "time"

type LoanApplication struct {
	Id             int64     `json:"id"`
	UserId         int64     `json:"userId"`
	Type           string    `json:"type"`
	VehicleVin     string    `json:"vehicleVin"`
	VehicleName    string    `json:"vehicleName"`
	CurrencyCode   string    `json:"currencyCode"`
	Price          int64     `json:"price"`
	DownPayment    int64     `json:"downPayment"`
	NetPrice       int64     `json:"netPrice"`
	MarginRate     float64   `json:"marginRate"`
	TermMonths     int32     `json:"termMonths"`
	MonthlyPayment int64     `json:"monthlyPayment"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Loan struct {
	Id               int64
	ApplicationId    int64
	UserId           int64
	CurrencyCode     string
	VehicleVin       string
	Amount           int64
	TermMonths       int32
	MonthlyPayment   int64
	RemainingBalance int64
	Status           string
	CreatedAt        time.Time
}

type Payment struct {
	Id            int64
	LoanId        int64
	CurrencyCode  string
	PaymentDate   time.Time
	Amount        int64
	Method        string
	Status        string
	TransactionId string
	CreatedAt     time.Time
}

type Vehicle struct {
	ImageURL      string
	Vin           string
	Name          string
	EngineType    string
	Configuration string
	Price         int64
	CurrencyCode  string
}
