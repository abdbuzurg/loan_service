package usecase

import (
	"loan_service/internal/clients"
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
