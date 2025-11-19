package usecase

import (
	"context"
	"fmt"
	"loan_service/internal/dto"
	"loan_service/internal/repository"
	"loan_service/pkg/utils"
)

func (uc *LoanUsecase) ListPayments(ctx context.Context, loanId int64, limit, offset int32) ([]*dto.Payment, error) {
	payments, err := uc.queries.ListPaymentsByLoan(ctx, repository.ListPaymentsByLoanParams{
		LoanID: loanId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get payments from db: %w", err)
	}

	result := make([]*dto.Payment, len(payments))

	for index, payment := range payments {
		result[index] = &dto.Payment{
			Id:            payment.ID,
			LoanId:        payment.LoanID,
			CurrencyCode:  payment.CurrencyCode,
			PaymentDate:   utils.NilToValueType(payment.PaymentDate),
			Amount:        int64(utils.NilToValueType(payment.Amount)),
			Method:        utils.NilToValueType(payment.Method),
			Status:        utils.NilToValueType(payment.Status),
			TransactionId: utils.NilToValueType(payment.TransactionID),
		}
	}

	return result, nil
}

func (uc *LoanUsecase) CountPayments(ctx context.Context, loanId int64) (*int64, error) {
	countPayments, err := uc.queries.CountPayments(ctx, loanId)
	if err != nil {
		return nil, fmt.Errorf("failed to count payments from db: %w", err)
	}

	return &countPayments, nil
}
