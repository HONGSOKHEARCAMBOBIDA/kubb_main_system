package model

type FeeInterval string

const (
	FeeIntervalMonthlyfee    FeeInterval = "monthly_fee"
	FeeIntervalQuarterlyFee  FeeInterval = "quarterly_fee"
	FeeIntervalSemesterlyFee FeeInterval = "semesterly_fee"
	FeeIntervalYearlyFee     FeeInterval = "yearly_fee"
)
