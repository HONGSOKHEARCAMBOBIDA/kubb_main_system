package utils

import (
	"context"

	"gorm.io/gorm"
)

func ToggleStatus[T any](ctx context.Context, db *gorm.DB, id string) error {
	result := db.WithContext(ctx).Model(new(T)).
		Where("uuid = ?", id).
		Update("active", gorm.Expr("NOT active"))
	return result.Error
}
