package request

import "mysql/model"

type FeeDiscountGroupRequestCreate struct {
	Code               string             `json:"code" gorm:"column:code" validate:"required,min=2,max=150"`
	Name               string             `json:"name" gorm:"column:name" validate:"required,min=2,max=150"`
	DiscountType       model.DiscountType `json:"discount_type"`
	DiscountPercentage float64            `json:"discount_percentage" gorm:"column:discount_percentage" validate:"required"`
	DiscountAmount     float64            `json:"discount_amount" gorm:"column:discount_amount" validate:"required"`
	Description        string             `json:"description" gorm:"column:description" validate:"max=1000"`
}

type FeeDiscountGroupRequestUpdate struct {
	Code               string             `json:"code" gorm:"column:code"`
	Name               string             `json:"name" gorm:"column:name"`
	DiscountType       model.DiscountType `json:"discount_type"`
	DiscountPercentage float64            `json:"discount_percentage" gorm:"column:discount_percentage"`
	DiscountAmount     float64            `json:"discount_amount" gorm:"column:discount_amount"`
	Description        string             `json:"description" gorm:"column:description"`
}
