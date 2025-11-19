package usecase

import (
	"context"
	"fmt"
	"loan_service/internal/dto"
	"loan_service/internal/repository"
	"loan_service/pkg/utils"
)

func (uc *LoanUsecase) CreateApplication(ctx context.Context, loanApp *dto.LoanApplication) (*dto.LoanApplication, error) {

	createdLoanApp, err := uc.queries.CreateApplication(ctx, repository.CreateApplicationParams{
		UserID:         loanApp.UserId,
		Type:           repository.ApplicationType(loanApp.Type),
		VehicleVin:     &loanApp.VehicleVin,
		VehicleName:    &loanApp.VehicleName,
		CurrencyCode:   loanApp.CurrencyCode,
		Price:          utils.PtrNumeric[int64, float64](loanApp.Price),
		DownPayment:    utils.PtrNumeric[int64, float64](loanApp.DownPayment),
		NetPrice:       utils.PtrNumeric[int64, float64](loanApp.NetPrice),
		MarginRate:     &loanApp.MarginRate,
		TermMonths:     utils.PtrNumeric[int32, int64](loanApp.TermMonths),
		MonthlyPayment: utils.PtrNumeric[int64, float64](loanApp.MonthlyPayment),
		Status: repository.NullApplicationStatus{
			ApplicationStatus: "NEW",
			Valid:             true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create loan application in db: %w", err)
	}

	loanApp.Id = createdLoanApp.ID
	if err := uc.koinotAutoClient.SendLoanApplication(ctx, loanApp); err != nil {
		return nil, err
	}

	return loanApp, nil
}

func (uc *LoanUsecase) GetApplication(ctx context.Context, id int64) (*dto.LoanApplication, error) {
	applicationResult, err := uc.queries.GetApplication(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan application from db: %w", err)
	}

	return &dto.LoanApplication{
		Id:             applicationResult.ID,
		UserId:         applicationResult.UserID,
		Type:           string(applicationResult.Type),
		VehicleVin:     utils.NilToValueType(applicationResult.VehicleVin),
		VehicleName:    utils.NilToValueType(applicationResult.VehicleName),
		CurrencyCode:   applicationResult.CurrencyCode,
		Price:          int64(utils.NilToValueType(applicationResult.Price)),
		DownPayment:    int64(utils.NilToValueType(applicationResult.DownPayment)),
		NetPrice:       int64(utils.NilToValueType(applicationResult.NetPrice)),
		MarginRate:     utils.NilToValueType(applicationResult.MarginRate),
		TermMonths:     int32(utils.NilToValueType(applicationResult.TermMonths)),
		MonthlyPayment: int64(utils.NilToValueType(applicationResult.MonthlyPayment)),
		Status:         string(applicationResult.Status.ApplicationStatus),
		CreatedAt:      utils.NilToValueType(applicationResult.CreatedAt),
		UpdatedAt:      utils.NilToValueType(applicationResult.UpdatedAt),
	}, nil
}

func (uc *LoanUsecase) ListApplications(ctx context.Context, userId int64, limit, offset int32) ([]*dto.LoanApplication, error) {
	loanApps, err := uc.queries.ListApplicationsByUser(ctx, repository.ListApplicationsByUserParams{
		UserID: userId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get loan applications from db: %w", err)
	}

	result := make([]*dto.LoanApplication, len(loanApps))
	for index, loanApp := range loanApps {
		result[index] = &dto.LoanApplication{
			Id:             loanApp.ID,
			UserId:         loanApp.UserID,
			Type:           string(loanApp.Type),
			VehicleVin:     utils.NilToValueType(loanApp.VehicleVin),
			VehicleName:    utils.NilToValueType(loanApp.VehicleName),
			CurrencyCode:   loanApp.CurrencyCode,
			Price:          int64(utils.NilToValueType(loanApp.Price)),
			DownPayment:    int64(utils.NilToValueType(loanApp.DownPayment)),
			NetPrice:       int64(utils.NilToValueType(loanApp.NetPrice)),
			MarginRate:     utils.NilToValueType(loanApp.MarginRate),
			TermMonths:     int32(utils.NilToValueType(loanApp.TermMonths)),
			MonthlyPayment: int64(utils.NilToValueType(loanApp.MonthlyPayment)),
			Status:         string(loanApp.Status.ApplicationStatus),
			CreatedAt:      utils.NilToValueType(loanApp.CreatedAt),
			UpdatedAt:      utils.NilToValueType(loanApp.UpdatedAt),
		}
	}

	return result, nil
}

func (uc *LoanUsecase) CountApplications(ctx context.Context, userId int64) (*int64, error) {
	loanAppCount, err := uc.queries.CountApplicationsByUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to count loan app: %w", err)
	}

	return &loanAppCount, nil
}
