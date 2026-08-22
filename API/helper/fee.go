package helper

import "mysql/model"

func GetNextDueDate(interval model.FeeInterval) (count int, months int) {
	switch interval {
	case model.FeeIntervalMonthlyfee:
		return 12, 1

	case model.FeeIntervalQuarterlyFee:
		return 4, 3

	case model.FeeIntervalSemesterlyFee:
		return 2, 6

	case model.FeeIntervalYearlyFee:
		return 1, 12

	default:
		return 0, 0
	}
}

func GetFeeSchedule(interval model.FeeInterval) int {
	switch interval {
	case model.FeeIntervalMonthlyfee:
		return 12
	case model.FeeIntervalQuarterlyFee:
		return 4
	case model.FeeIntervalSemesterlyFee:
		return 2
	case model.FeeIntervalYearlyFee:
		return 1
	default:
		return 0
	}
}

func GetFeeAmountPerYear(degree model.AcademicDegree, interval model.FeeInterval) float64 {
	switch interval {
	case model.FeeIntervalMonthlyfee:
		return degree.MonthlyFee * 12
	case model.FeeIntervalQuarterlyFee:
		return degree.QuarterlyFee * 4
	case model.FeeIntervalSemesterlyFee:
		return degree.SemesterlyFee * 2
	case model.FeeIntervalYearlyFee:
		return degree.YearlyFee
	default:
		return 0
	}
}

func GetFeeAmountByInterval(degree model.AcademicDegree, interval model.FeeInterval) float64 {
	switch interval {
	case model.FeeIntervalMonthlyfee:
		return degree.MonthlyFee
	case model.FeeIntervalQuarterlyFee:
		return degree.QuarterlyFee
	case model.FeeIntervalSemesterlyFee:
		return degree.SemesterlyFee
	case model.FeeIntervalYearlyFee:
		return degree.YearlyFee
	default:
		return 0
	}

}

func CalculateDiscount(amount float64, group *model.FeeDiscountGroup) float64 {
	if group == nil || !group.Active {
		return 0
	}
	switch group.DiscountType {
	case model.DiscountPercentage:
		return amount * (group.DiscountPercentage / 100)
	case model.DiscountAmount:
		return group.DiscountAmount
	default:
		return 0
	}
}

func CalculateDiscountBySchoolarship(amount float64, group *model.Schoolarship) float64 {
	if group == nil || !group.Active {
		return 0
	}
	switch group.DiscountType {
	case model.DiscountPercentage:
		return amount * (group.DiscountPercentage / 100)
	case model.DiscountAmount:
		return group.DiscountAmount
	default:
		return 0
	}
}
