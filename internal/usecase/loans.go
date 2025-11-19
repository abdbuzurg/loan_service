package usecase

import (
	"context"
	"fmt"
	"loan_service/internal/dto"
	"loan_service/internal/repository"
	"loan_service/pkg/utils"
)

func (uc *LoanUsecase) GetLoan(ctx context.Context, id int64) (*dto.Loan, error) {
	loan, err := uc.queries.GetLoan(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan from db: %w", err)
	}

	return &dto.Loan{
		Id:               loan.ID,
		ApplicationId:    loan.ApplicationID,
		UserId:           loan.UserID,
		CurrencyCode:     loan.CurrencyCode,
		VehicleVin:       utils.NilToValueType(loan.VehicleVin),
		Amount:           int64(utils.NilToValueType(loan.Amount)),
		TermMonths:       int32(utils.NilToValueType(loan.TermMonths)),
		MonthlyPayment:   int64(utils.NilToValueType(loan.MonthlyPayment)),
		RemainingBalance: int64(utils.NilToValueType(loan.RemainingBalance)),
		Status:           string(loan.Status.LoanStatus),
		CreatedAt:        utils.NilToValueType(loan.CreatedAt),
	}, nil
}

func (uc *LoanUsecase) ListLoans(ctx context.Context, userId int64, limit, offset int32) ([]*dto.Loan, error) {
	loans, err := uc.queries.ListLoansByUser(ctx, repository.ListLoansByUserParams{
		UserID: userId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get loans from db: %w", err)
	}

	result := make([]*dto.Loan, len(loans))

	for index, loan := range loans {
		result[index] = &dto.Loan{
			Id:               loan.UserID,
			ApplicationId:    loan.ApplicationID,
			UserId:           loan.UserID,
			CurrencyCode:     loan.CurrencyCode,
			VehicleVin:       utils.NilToValueType(loan.VehicleVin),
			Amount:           int64(utils.NilToValueType(loan.Amount)),
			TermMonths:       int32(utils.NilToValueType(loan.TermMonths)),
			MonthlyPayment:   int64(utils.NilToValueType(loan.MonthlyPayment)),
			RemainingBalance: int64(utils.NilToValueType(loan.RemainingBalance)),
			Status:           string(loan.Status.LoanStatus),
			CreatedAt:        utils.NilToValueType(loan.CreatedAt),
		}
	}

	return result, nil
}

func (uc *LoanUsecase) CountLoans(ctx context.Context, userId int64) (*int64, error) {
	countLoans, err := uc.queries.CountLoansByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to count loans from db: %w", err)
	}

	return &countLoans, nil
}
