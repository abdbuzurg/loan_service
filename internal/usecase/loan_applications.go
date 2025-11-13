package usecase

func (uc *LoanUsecase) Calculate(price, downPayment int64, termMonths int32, marginRate float64) (int64, int64, int64) {
	net := price - downPayment
	years := float64(termMonths) / 12
	margin := int64((float64(net)*marginRate/100)*years + 0.5)
	total := net + margin
	monthly := int64((float64(total) / float64(termMonths)) + 0.5)

	return net, monthly, total
}
