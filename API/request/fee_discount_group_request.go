package request

import "mysql/model"

type FeeDiscountGroupRequestCreate struct {
	Code               string             `json:"code" gorm:"column:code" validate:"omitempty,min=2,max=150"`
	Name               string             `json:"name" gorm:"column:name" validate:"omitempty,min=2,max=150"`
	DiscountType       model.DiscountType `json:"discount_type"`
	DiscountPercentage float64            `json:"discount_percentage" gorm:"column:discount_percentage" validate:"omitempty"`
	DiscountAmount     float64            `json:"discount_amount" gorm:"column:discount_amount" validate:"omitempty"`
	Description        string             `json:"description" gorm:"column:description" validate:"omitempty,max=1000"`
}

type FeeDiscountGroupRequestUpdate struct {
	Code               *string             `json:"code" validate:"omitempty,min=2,max=150"`
	Name               *string             `json:"name" validate:"omitempty,min=2,max=150"`
	DiscountType       *model.DiscountType `json:"discount_type" validate:"omitempty"`
	DiscountPercentage *float64            `json:"discount_percentage" validate:"omitempty,gte=0,lte=100"`
	DiscountAmount     *float64            `json:"discount_amount" validate:"omitempty,gte=0"`
	Description        *string             `json:"description" validate:"omitempty,max=1000"`
}
