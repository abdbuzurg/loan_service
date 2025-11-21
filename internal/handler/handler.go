package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"loan_service/internal/dto"
	loanpb "loan_service/internal/proto/loan"
	"loan_service/internal/usecase"
	"strconv"
	"time"
)

type LoanHandler struct {
	loanpb.UnimplementedLoansServiceServer
	loanUC *usecase.LoanUsecase
}

func New(loanUC *usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{
		loanUC: loanUC,
	}
}

func ok() *loanpb.LoanServiceError {
	return &loanpb.LoanServiceError{
		Code:        0,
		Description: "",
	}
}

func pageToLimitOffset(p *loanpb.PageRequest) (limit, offset int32) {
	page, limitIn := int32(1), int32(20)
	if p != nil {
		if p.Page > 0 {
			page = p.Page
		}

		if p.Limit > 0 {
			limitIn = p.Limit
		}
	}

	offset = (page - 1) * limitIn
	return limitIn, offset
}

func (h *LoanHandler) Calculate(ctx context.Context, calculateRequest *loanpb.CalculateRequest) (*loanpb.CalculateResponse, error) {
	net, monthly, total := h.loanUC.Calculate(
		calculateRequest.Price,
		calculateRequest.DownPayment,
		calculateRequest.TermMonths,
		calculateRequest.MarginRate,
	)

	return &loanpb.CalculateResponse{
		NetPrice:         net,
		MonthlyPayment:   monthly,
		TotalAmount:      total,
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) CreateApplication(ctx context.Context, req *loanpb.CreateApplicationRequest) (*loanpb.CreateApplicationResponse, error) {
	userId, err := strconv.ParseInt(req.GetUserId(), 10, 64)
	if err != nil {
		return &loanpb.CreateApplicationResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: "user id is required",
			},
		}, nil
	}

	createdLoanApp, err := h.loanUC.CreateApplication(ctx, &dto.LoanApplication{
		Id:             0,
		UserId:         userId,
		Type:           req.GetType(),
		VehicleVin:     req.GetVehicleVin(),
		VehicleName:    req.GetVehicleName(),
		CurrencyCode:   req.GetCurrencyCode(),
		Price:          req.GetPrice(),
		DownPayment:    req.GetDownPayment(),
		MarginRate:     req.GetMarginRate(),
		NetPrice:       req.GetNetPrice(),
		TermMonths:     req.GetTermMonths(),
		MonthlyPayment: req.GetMonthlyPayment(),
		Status:         "NEW",
	})
	if err != nil {
		return &loanpb.CreateApplicationResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to create loan application",
			},
		}, nil
	}

	return &loanpb.CreateApplicationResponse{
		Application: &loanpb.LoanApplication{
			Id:             fmt.Sprint(createdLoanApp.Id),
			UserId:         fmt.Sprint(createdLoanApp.UserId),
			Type:           createdLoanApp.Type,
			VehicleVin:     createdLoanApp.VehicleVin,
			VehicleName:    createdLoanApp.VehicleName,
			CurrencyCode:   createdLoanApp.CurrencyCode,
			Price:          createdLoanApp.Price,
			DownPayment:    createdLoanApp.DownPayment,
			NetPrice:       createdLoanApp.NetPrice,
			MarginRate:     createdLoanApp.MarginRate,
			TermMonths:     createdLoanApp.TermMonths,
			MonthlyPayment: createdLoanApp.MonthlyPayment,
			Status:         createdLoanApp.Status,
			CreatedAt:      createdLoanApp.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      createdLoanApp.UpdatedAt.Format(time.RFC3339),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) GetApplication(ctx context.Context, req *loanpb.GetApplicationRequest) (*loanpb.GetApplicationResponse, error) {
	if req.GetId() == "" {
		return &loanpb.GetApplicationResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: "id is required",
			},
		}, nil
	}

	loanAppId, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		return &loanpb.GetApplicationResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid id %q", req.GetId()),
			},
		}, nil
	}

	loanApplication, err := h.loanUC.GetApplication(ctx, loanAppId)
	if err != nil {
		// Not found
		if errors.Is(err, sql.ErrNoRows) {
			return &loanpb.GetApplicationResponse{
				LoanServiceError: &loanpb.LoanServiceError{
					Code:        2,
					Description: "application not found",
				},
			}, nil
		}

		// Internal Error
		return &loanpb.GetApplicationResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch application",
			},
		}, nil
	}

	return &loanpb.GetApplicationResponse{
		Application: &loanpb.LoanApplication{
			Id:             fmt.Sprint(loanApplication.Id),
			UserId:         fmt.Sprint(loanApplication.UserId),
			Type:           loanApplication.Type,
			VehicleVin:     loanApplication.VehicleVin,
			VehicleName:    loanApplication.VehicleName,
			CurrencyCode:   loanApplication.CurrencyCode,
			Price:          loanApplication.Price,
			DownPayment:    loanApplication.DownPayment,
			NetPrice:       loanApplication.NetPrice,
			MarginRate:     loanApplication.MarginRate,
			TermMonths:     loanApplication.TermMonths,
			MonthlyPayment: loanApplication.MonthlyPayment,
			Status:         loanApplication.Status,
			CreatedAt:      loanApplication.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      loanApplication.UpdatedAt.Format(time.RFC3339),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) GetLoan(ctx context.Context, req *loanpb.GetLoanRequest) (*loanpb.GetLoanResponse, error) {
	if req.GetId() == "" {
		return &loanpb.GetLoanResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: "id is required",
			},
		}, nil
	}

	loanId, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		return &loanpb.GetLoanResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid id %q", req.GetId()),
			},
		}, nil
	}

	loan, err := h.loanUC.GetLoan(ctx, loanId)
	if err != nil {
		// Not found
		if errors.Is(err, sql.ErrNoRows) {
			return &loanpb.GetLoanResponse{
				LoanServiceError: &loanpb.LoanServiceError{
					Code:        2,
					Description: "loan not found",
				},
			}, nil
		}

		//Internal errro
		return &loanpb.GetLoanResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan",
			},
		}, nil
	}

	return &loanpb.GetLoanResponse{
		Loan: &loanpb.Loan{
			Id:               fmt.Sprint(loan.Id),
			ApplicationId:    fmt.Sprint(loan.ApplicationId),
			UserId:           fmt.Sprint(loan.UserId),
			CurrencyCode:     loan.CurrencyCode,
			VehicleVin:       loan.VehicleVin,
			Amount:           loan.Amount,
			TermMonths:       loan.TermMonths,
			MonthlyPayment:   loan.MonthlyPayment,
			RemainingBalance: loan.RemainingBalance,
			Status:           loan.Status,
			CreatedAt:        loan.CreatedAt.Format(time.RFC3339),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) ListApplications(ctx context.Context, req *loanpb.ListApplicationsRequest) (*loanpb.ListApplicationsResponse, error) {
	if req.GetUserId() == "" {
		return &loanpb.ListApplicationsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: "user id is required",
			},
		}, nil
	}

	userId, err := strconv.ParseInt(req.GetUserId(), 10, 64)
	if err != nil {
		return &loanpb.ListApplicationsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid user id %q", req.GetUserId()),
			},
		}, nil
	}

	pageInfo := req.GetPage()
	limit, offset := pageToLimitOffset(pageInfo)

	loanAppsCount, err := h.loanUC.CountApplications(ctx, userId)
	if err != nil {
		return &loanpb.ListApplicationsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applicaitons",
			},
		}, nil
	}

	if *loanAppsCount == 0 {
		return &loanpb.ListApplicationsResponse{
			Applications: nil,
			Page: &loanpb.PageResponse{
				CurrentPage: pageInfo.Page,
				Limit:       pageInfo.Page,
				TotalItems:  0,
				TotalPages:  0,
			},
			LoanServiceError: ok(),
		}, nil
	}

	loanApps, err := h.loanUC.ListApplications(ctx, userId, limit, offset)
	if err != nil {
		// Not found
		if errors.Is(err, sql.ErrNoRows) {
			return &loanpb.ListApplicationsResponse{
				LoanServiceError: &loanpb.LoanServiceError{
					Code:        2,
					Description: "loan applications for user not found",
				},
			}, nil
		}

		// Internal error
		return &loanpb.ListApplicationsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applications",
			},
		}, nil
	}

	listLoanAppsPB := make([]*loanpb.LoanApplication, len(loanApps))
	for index, loanApp := range loanApps {
		listLoanAppsPB[index] = &loanpb.LoanApplication{
			Id:             fmt.Sprint(loanApp.Id),
			UserId:         fmt.Sprint(loanApp.UserId),
			Type:           loanApp.Type,
			VehicleVin:     loanApp.VehicleVin,
			VehicleName:    loanApp.VehicleName,
			CurrencyCode:   loanApp.CurrencyCode,
			Price:          loanApp.Price,
			DownPayment:    loanApp.DownPayment,
			NetPrice:       loanApp.NetPrice,
			MarginRate:     loanApp.MarginRate,
			TermMonths:     loanApp.TermMonths,
			MonthlyPayment: loanApp.MonthlyPayment,
			Status:         loanApp.Status,
			CreatedAt:      loanApp.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      loanApp.UpdatedAt.Format(time.RFC3339),
		}
	}

	totalPages := *loanAppsCount / int64(pageInfo.Limit)
	if *loanAppsCount%int64(pageInfo.Limit) != 0 {
		totalPages++
	}
	return &loanpb.ListApplicationsResponse{
		Applications: listLoanAppsPB,
		Page: &loanpb.PageResponse{
			CurrentPage: pageInfo.Page,
			Limit:       pageInfo.Limit,
			TotalItems:  int32(*loanAppsCount),
			TotalPages:  int32(totalPages),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) ListLoans(ctx context.Context, req *loanpb.ListLoansRequest) (*loanpb.ListLoansResponse, error) {
	if req.GetUserId() == "" {
		return &loanpb.ListLoansResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: "user id is required",
			},
		}, nil
	}

	userId, err := strconv.ParseInt(req.GetUserId(), 10, 64)
	if err != nil {
		return &loanpb.ListLoansResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid user id %q", req.GetUserId()),
			},
		}, nil
	}

	pageInfo := req.GetPage()
	limit, offset := pageToLimitOffset(pageInfo)

	loansCount, err := h.loanUC.CountLoans(ctx, userId)
	if err != nil {
		return &loanpb.ListLoansResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applicaitons",
			},
		}, nil
	}

	if *loansCount == 0 {
		return &loanpb.ListLoansResponse{
			Loans: nil,
			Page: &loanpb.PageResponse{
				CurrentPage: pageInfo.Page,
				Limit:       pageInfo.Page,
				TotalItems:  0,
				TotalPages:  0,
			},
			LoanServiceError: ok(),
		}, nil
	}

	loans, err := h.loanUC.ListLoans(ctx, userId, limit, offset)
	if err != nil {
		// Not found
		if errors.Is(err, sql.ErrNoRows) {
			return &loanpb.ListLoansResponse{
				LoanServiceError: &loanpb.LoanServiceError{
					Code:        2,
					Description: "loans for user not found",
				},
			}, nil
		}

		// Internal error
		return &loanpb.ListLoansResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applications",
			},
		}, nil
	}

	listLoansPB := make([]*loanpb.Loan, len(loans))
	for index, loan := range loans {
		listLoansPB[index] = &loanpb.Loan{
			Id:               fmt.Sprint(loan.Id),
			ApplicationId:    fmt.Sprint(loan.ApplicationId),
			UserId:           fmt.Sprint(loan.UserId),
			CurrencyCode:     loan.CurrencyCode,
			VehicleVin:       loan.VehicleVin,
			Amount:           loan.Amount,
			TermMonths:       loan.TermMonths,
			MonthlyPayment:   loan.MonthlyPayment,
			RemainingBalance: loan.RemainingBalance,
			Status:           loan.Status,
			CreatedAt:        loan.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := *loansCount / int64(pageInfo.Limit)
	if *loansCount%int64(pageInfo.Limit) != 0 {
		totalPages++
	}
	return &loanpb.ListLoansResponse{
		Loans: listLoansPB,
		Page: &loanpb.PageResponse{
			CurrentPage: pageInfo.Page,
			Limit:       pageInfo.Limit,
			TotalItems:  int32(*loansCount),
			TotalPages:  int32(totalPages),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) ListPayments(ctx context.Context, req *loanpb.ListPaymentsRequest) (*loanpb.ListPaymentsResponse, error) {
	if req.GetLoanId() == "" {
		return &loanpb.ListPaymentsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: "user id is required",
			},
		}, nil
	}

	loanId, err := strconv.ParseInt(req.GetLoanId(), 10, 64)
	if err != nil {
		return &loanpb.ListPaymentsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid loan id %q", req.GetLoanId()),
			},
		}, nil
	}

	pageInfo := req.GetPage()
	limit, offset := pageToLimitOffset(pageInfo)

	paymentsCount, err := h.loanUC.CountPayments(ctx, loanId)
	if err != nil {
		return &loanpb.ListPaymentsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applicaitons",
			},
		}, nil
	}

	if *paymentsCount == 0 {
		return &loanpb.ListPaymentsResponse{
			Payments: nil,
			Page: &loanpb.PageResponse{
				CurrentPage: pageInfo.Page,
				Limit:       pageInfo.Page,
				TotalItems:  0,
				TotalPages:  0,
			},
			LoanServiceError: ok(),
		}, nil
	}

	payments, err := h.loanUC.ListPayments(ctx, loanId, limit, offset)
	if err != nil {
		// Not found
		if errors.Is(err, sql.ErrNoRows) {
			return &loanpb.ListPaymentsResponse{
				LoanServiceError: &loanpb.LoanServiceError{
					Code:        2,
					Description: "loans for user not found",
				},
			}, nil
		}

		// Internal error
		return &loanpb.ListPaymentsResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applications",
			},
		}, nil
	}

	listPaymentsPB := make([]*loanpb.Payment, len(payments))
	for index, payment := range payments {
		listPaymentsPB[index] = &loanpb.Payment{
			Id:            fmt.Sprint(payment.Id),
			LoanId:        fmt.Sprint(payment.LoanId),
			CurrencyCode:  payment.CurrencyCode,
			PaymentDate:   payment.PaymentDate.Format(time.RFC3339),
			Amount:        payment.Amount,
			Method:        payment.Method,
			Status:        payment.Status,
			TransactionId: payment.TransactionId,
			CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := *paymentsCount / int64(pageInfo.Limit)
	if *paymentsCount%int64(pageInfo.Limit) != 0 {
		totalPages++
	}
	return &loanpb.ListPaymentsResponse{
		Payments: listPaymentsPB,
		Page: &loanpb.PageResponse{
			CurrentPage: pageInfo.Page,
			Limit:       pageInfo.Limit,
			TotalItems:  int32(*paymentsCount),
			TotalPages:  int32(totalPages),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) ListVehicles(ctx context.Context, req *loanpb.ListVehiclesRequest) (*loanpb.ListVehiclesResponse, error) {
	vehicles, err := h.loanUC.ListVehicles(ctx)
	if err != nil {
		return &loanpb.ListVehiclesResponse{
			LoanServiceError: &loanpb.LoanServiceError{
				Code:        5,
				Description: "failed to get vehicles",
			},
		}, nil
	}

	vehiclesPB := make([]*loanpb.Vehicle, len(vehicles))
	for index, vehicle := range vehicles {
		vehiclesPB[index] = &loanpb.Vehicle{
			ImageUrl:      vehicle.ImageURL,
			Vin:           vehicle.Vin,
			Name:          vehicle.Name,
			EngineType:    vehicle.EngineType,
			Configuration: vehicle.Configuration,
			Price:         vehicle.Price,
			CurrencyCode:  vehicle.CurrencyCode,
		}
	}

	return &loanpb.ListVehiclesResponse{
		Vehicles:         vehiclesPB,
		LoanServiceError: ok(),
	}, nil
}
