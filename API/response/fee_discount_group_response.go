package response

import (
	"mysql/model"
	"mysql/model/base"
)

type FeeDiscountGroupResponse struct {
	base.ModelBase
	base.UUIDBase
	Code               string             `json:"code" gorm:"column:code"`
	Name               string             `json:"name" gorm:"column:name"`
	DiscountType       model.DiscountType `json:"discount_type"`
	DiscountPercentage float64            `json:"discount_percentage" gorm:"column:discount_percentage"`
	DiscountAmount     float64            `json:"discount_amount" gorm:"column:discount_amount"`
	Description        string             `json:"description" gorm:"column:description"`
	Active             bool               `json:"active" gorm:"column:active"`
}
