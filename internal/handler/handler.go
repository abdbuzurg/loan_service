package handler

import (
	"context"
	"loan_service/internal/proto"
	"loan_service/internal/usecase"
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

func (h *LoanHandler) CreateApplication(context.Context, *proto.CreateApplicationRequest) (*proto.CreateApplicationResponse, error) {
	return nil, nil
}

func (h *LoanHandler) CreatePayment(context.Context, *proto.CreatePaymentRequest) (*proto.CreatePaymentResponse, error) {
	return nil, nil
}

func (h *LoanHandler) GetApplication(context.Context, *proto.GetApplicationRequest) (*proto.GetApplicationResponse, error) {
	return nil, nil
}

func (h *LoanHandler) GetLoan(context.Context, *proto.GetLoanRequest) (*proto.GetLoanResponse, error) {
	return nil, nil
}

func (h *LoanHandler) ListApplications(context.Context, *proto.ListApplicationsRequest) (*proto.ListApplicationsResponse, error) {
	return nil, nil
}

func (h *LoanHandler) ListLoans(context.Context, *proto.ListLoansRequest) (*proto.ListLoansResponse, error) {
	return nil, nil
}

func (h *LoanHandler) ListPayments(context.Context, *proto.ListPaymentsRequest) (*proto.ListPaymentsResponse, error) {
	return nil, nil
}

func (h *LoanHandler) ListVehicles(context.Context, *proto.ListVehiclesRequest) (*proto.ListVehiclesResponse, error) {
	return nil, nil
}
