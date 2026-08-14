package seed

import (
	"mysql/constant/seeddata"
	"mysql/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedDocumentTypes(db *gorm.DB) error {

	var documenttype []model.DocumentType

	for _, p := range seeddata.DocumentType {
		documenttype = append(documenttype, model.DocumentType{
			Code:   p.Code,
			NameKh: p.NameKh,
			NameEn: p.NameEn,
		})
	}

	return db.Clauses(clause.OnConflict{
		// OnConflict is Upsert = Update + Insert
		Columns: []clause.Column{
			{Name: "code"}, // unique key
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"code",
		}),
	}).Create(&documenttype).Error
}
