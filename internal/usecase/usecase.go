package usecase

import (
	"context"
	"loan_service/internal/clients"
	"loan_service/internal/dto"
	"loan_service/internal/repository"
)

type LoanUsecase struct {
	queries          repository.Queries
	asrLeasingClient clients.AsrLeasingClient
	koinotAutoClient clients.KoinotAutoClient
}

func New(
	queries repository.Queries,
	asrLeasingClient clients.AsrLeasingClient,
	koinotAutoClient clients.KoinotAutoClient,
) *LoanUsecase {
	return &LoanUsecase{
		queries:          queries,
		asrLeasingClient: asrLeasingClient,
		koinotAutoClient: koinotAutoClient,
	}
}

func (uc *LoanUsecase) Calculate(price, downPayment int64, termMonths int32, marginRate float64) (int64, int64, int64) {
	net := price - downPayment
	years := float64(termMonths) / 12
	margin := int64((float64(net)*marginRate/100)*years + 0.5)
	total := net + margin
	monthly := int64((float64(total) / float64(termMonths)) + 0.5)

	return net, monthly, total
}

func (uc *LoanUsecase) ListVehicles(ctx context.Context) ([]dto.Vehicle, error) {
	return uc.koinotAutoClient.ListVehicles(ctx)
}
