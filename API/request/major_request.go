package request

import "mysql/model"

type MajorRequestCreate struct {
	DepartmentID     uint                   `json:"department_id" binding:"required"`
	Code             string                 `json:"code"`
	Name             string                 `json:"name" binding:"required"`
	DurationPeriod   int                    `json:"duration_period"`
	DurationInterval model.DurationInterval `json:"duration_interval"`
	Description      string                 `json:"description"`
}

type MajorRequestUpdate struct {
	DepartmentID     *uint                   `json:"department_id" binding:"required"`
	Code             *string                 `json:"code"`
	Name             *string                 `json:"name" binding:"required"`
	DurationPeriod   *int                    `json:"duration_period"`
	DurationInterval *model.DurationInterval `json:"duration_interval"`
	Description      *string                 `json:"description"`
}
