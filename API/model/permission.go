package model

type Permission struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Group string `json:"group" gorm:"column:group"`
}

func (Permission) TableName() string {
	return "permissions"
}
