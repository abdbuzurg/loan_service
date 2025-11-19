package dto

import "time"

type LoanApplication struct {
	Id             int64     `url:"id"`
	UserId         int64     `url:"userId"`
	Type           string    `url:"type"`
	VehicleVin     string    `url:"vehicleVin"`
	VehicleName    string    `url:"vehicleName"`
	CurrencyCode   string    `url:"currencyCode"`
	Price          int64     `url:"price"`
	DownPayment    int64     `url:"downPayment"`
	NetPrice       int64     `url:"netPrice"`
	MarginRate     float64   `url:"marginRate"`
	TermMonths     int32     `url:"termMonths"`
	MonthlyPayment int64     `url:"monthlyPayment"`
	Status         string    `url:"status"`
	CreatedAt      time.Time `url:"createdAt"`
	UpdatedAt      time.Time `url:"updatedAt"`
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
