package base

type ModelBase struct {
	ID int `json:"id" gorm:"column:id"`
}

type UUIDBase struct {
	UUID string `json:"uuid" gorm:"column:uuid"`
}
