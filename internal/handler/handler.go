package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"loan_service/internal/dto"
	"loan_service/internal/proto"
	"loan_service/internal/usecase"
	"strconv"
	"time"
)

type LoanHandler struct {
	proto.UnimplementedLoansServiceServer
	loanUC *usecase.LoanUsecase
}

func New(loanUC *usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{
		loanUC: loanUC,
	}
}

func ok() *proto.LoanServiceError {
	return &proto.LoanServiceError{
		Code:        0,
		Description: "",
	}
}

func pageToLimitOffset(p *proto.PageRequest) (limit, offset int32) {
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

func (h *LoanHandler) Calculate(ctx context.Context, calculateRequest *proto.CalculateRequest) (*proto.CalculateResponse, error) {
	net, monthly, total := h.loanUC.Calculate(
		calculateRequest.Price,
		calculateRequest.DownPayment,
		calculateRequest.TermMonths,
		calculateRequest.MarginRate,
	)

	return &proto.CalculateResponse{
		NetPrice:         net,
		MonthlyPayment:   monthly,
		TotalAmount:      total,
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) CreateApplication(ctx context.Context, req *proto.CreateApplicationRequest) (*proto.CreateApplicationResponse, error) {
	userId, err := strconv.ParseInt(req.GetUserId(), 10, 64)
	if err != nil {
		return &proto.CreateApplicationResponse{
			LoanServiceError: &proto.LoanServiceError{
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
		return &proto.CreateApplicationResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to create loan application",
			},
		}, nil
	}

	return &proto.CreateApplicationResponse{
		Application: &proto.LoanApplication{
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

func (h *LoanHandler) CreatePayment(context.Context, *proto.CreatePaymentRequest) (*proto.CreatePaymentResponse, error) {
	return nil, nil
}

func (h *LoanHandler) GetApplication(ctx context.Context, req *proto.GetApplicationRequest) (*proto.GetApplicationResponse, error) {
	if req.GetId() == "" {
		return &proto.GetApplicationResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: "id is required",
			},
		}, nil
	}

	loanAppId, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		return &proto.GetApplicationResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid id %q", req.GetId()),
			},
		}, nil
	}

	loanApplication, err := h.loanUC.GetApplication(ctx, loanAppId)
	if err != nil {
		// Not found
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.GetApplicationResponse{
				LoanServiceError: &proto.LoanServiceError{
					Code:        2,
					Description: "application not found",
				},
			}, nil
		}

		// Internal Error
		return &proto.GetApplicationResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch application",
			},
		}, nil
	}

	return &proto.GetApplicationResponse{
		Application: &proto.LoanApplication{
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

func (h *LoanHandler) GetLoan(ctx context.Context, req *proto.GetLoanRequest) (*proto.GetLoanResponse, error) {
	if req.GetId() == "" {
		return &proto.GetLoanResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: "id is required",
			},
		}, nil
	}

	loanId, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		return &proto.GetLoanResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid id %q", req.GetId()),
			},
		}, nil
	}

	loan, err := h.loanUC.GetLoan(ctx, loanId)
	if err != nil {
		// Not found
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.GetLoanResponse{
				LoanServiceError: &proto.LoanServiceError{
					Code:        2,
					Description: "loan not found",
				},
			}, nil
		}

		//Internal errro
		return &proto.GetLoanResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan",
			},
		}, nil
	}

	return &proto.GetLoanResponse{
		Loan: &proto.Loan{
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

func (h *LoanHandler) ListApplications(ctx context.Context, req *proto.ListApplicationsRequest) (*proto.ListApplicationsResponse, error) {
	if req.GetUserId() == "" {
		return &proto.ListApplicationsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: "user id is required",
			},
		}, nil
	}

	userId, err := strconv.ParseInt(req.GetUserId(), 10, 64)
	if err != nil {
		return &proto.ListApplicationsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid user id %q", req.GetUserId()),
			},
		}, nil
	}

	pageInfo := req.GetPage()
	limit, offset := pageToLimitOffset(pageInfo)

	loanAppsCount, err := h.loanUC.CountApplications(ctx, userId)
	if err != nil {
		return &proto.ListApplicationsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applicaitons",
			},
		}, nil
	}

	if *loanAppsCount == 0 {
		return &proto.ListApplicationsResponse{
			Applications: nil,
			Page: &proto.PageResponse{
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
			return &proto.ListApplicationsResponse{
				LoanServiceError: &proto.LoanServiceError{
					Code:        2,
					Description: "loan applications for user not found",
				},
			}, nil
		}

		// Internal error
		return &proto.ListApplicationsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applications",
			},
		}, nil
	}

	listLoanAppsPB := make([]*proto.LoanApplication, len(loanApps))
	for index, loanApp := range loanApps {
		listLoanAppsPB[index] = &proto.LoanApplication{
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
	return &proto.ListApplicationsResponse{
		Applications: listLoanAppsPB,
		Page: &proto.PageResponse{
			CurrentPage: pageInfo.Page,
			Limit:       pageInfo.Limit,
			TotalItems:  int32(*loanAppsCount),
			TotalPages:  int32(totalPages),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) ListLoans(ctx context.Context, req *proto.ListLoansRequest) (*proto.ListLoansResponse, error) {
	if req.GetUserId() == "" {
		return &proto.ListLoansResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: "user id is required",
			},
		}, nil
	}

	userId, err := strconv.ParseInt(req.GetUserId(), 10, 64)
	if err != nil {
		return &proto.ListLoansResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid user id %q", req.GetUserId()),
			},
		}, nil
	}

	pageInfo := req.GetPage()
	limit, offset := pageToLimitOffset(pageInfo)

	loansCount, err := h.loanUC.CountLoans(ctx, userId)
	if err != nil {
		return &proto.ListLoansResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applicaitons",
			},
		}, nil
	}

	if *loansCount == 0 {
		return &proto.ListLoansResponse{
			Loans: nil,
			Page: &proto.PageResponse{
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
			return &proto.ListLoansResponse{
				LoanServiceError: &proto.LoanServiceError{
					Code:        2,
					Description: "loans for user not found",
				},
			}, nil
		}

		// Internal error
		return &proto.ListLoansResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applications",
			},
		}, nil
	}

	listLoansPB := make([]*proto.Loan, len(loans))
	for index, loan := range loans {
		listLoansPB[index] = &proto.Loan{
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
	return &proto.ListLoansResponse{
		Loans: listLoansPB,
		Page: &proto.PageResponse{
			CurrentPage: pageInfo.Page,
			Limit:       pageInfo.Limit,
			TotalItems:  int32(*loansCount),
			TotalPages:  int32(totalPages),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) ListPayments(ctx context.Context, req *proto.ListPaymentsRequest) (*proto.ListPaymentsResponse, error) {
	if req.GetLoanId() == "" {
		return &proto.ListPaymentsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: "user id is required",
			},
		}, nil
	}

	loanId, err := strconv.ParseInt(req.GetLoanId(), 10, 64)
	if err != nil {
		return &proto.ListPaymentsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        1,
				Description: fmt.Sprintf("invalid loan id %q", req.GetLoanId()),
			},
		}, nil
	}

	pageInfo := req.GetPage()
	limit, offset := pageToLimitOffset(pageInfo)

	paymentsCount, err := h.loanUC.CountPayments(ctx, loanId)
	if err != nil {
		return &proto.ListPaymentsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applicaitons",
			},
		}, nil
	}

	if *paymentsCount == 0 {
		return &proto.ListPaymentsResponse{
			Payments: nil,
			Page: &proto.PageResponse{
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
			return &proto.ListPaymentsResponse{
				LoanServiceError: &proto.LoanServiceError{
					Code:        2,
					Description: "loans for user not found",
				},
			}, nil
		}

		// Internal error
		return &proto.ListPaymentsResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to fetch loan applications",
			},
		}, nil
	}

	listPaymentsPB := make([]*proto.Payment, len(payments))
	for index, payment := range payments {
		listPaymentsPB[index] = &proto.Payment{
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
	return &proto.ListPaymentsResponse{
		Payments: listPaymentsPB,
		Page: &proto.PageResponse{
			CurrentPage: pageInfo.Page,
			Limit:       pageInfo.Limit,
			TotalItems:  int32(*paymentsCount),
			TotalPages:  int32(totalPages),
		},
		LoanServiceError: ok(),
	}, nil
}

func (h *LoanHandler) ListVehicles(ctx context.Context, req *proto.ListVehiclesRequest) (*proto.ListVehiclesResponse, error) {
	vehicles, err := h.loanUC.ListVehicles(ctx)
	if err != nil {
		return &proto.ListVehiclesResponse{
			LoanServiceError: &proto.LoanServiceError{
				Code:        5,
				Description: "failed to get vehicles",
			},
		}, nil
	}

	vehiclesPB := make([]*proto.Vehicle, len(vehicles))
	for index, vehicle := range vehicles {
		vehiclesPB[index] = &proto.Vehicle{
			ImageUrl:      vehicle.ImageURL,
			Vin:           vehicle.Vin,
			Name:          vehicle.Name,
			EngineType:    vehicle.EngineType,
			Configuration: vehicle.Configuration,
			Price:         vehicle.Price,
			CurrencyCode:  vehicle.CurrencyCode,
		}
	}

	return &proto.ListVehiclesResponse{
		Vehicles:         vehiclesPB,
		LoanServiceError: ok(),
	}, nil
}
