package response

import "mysql/model/base"

type StudentDocumentResponse struct {
	base.ModelBase
	base.UUIDBase
	StudentID          int    `json:"student_id" gorm:"column:student_id"`
	DocumentTypeID     int    `json:"document_type_id" gorm:"column:document_type_id"`
	DocumentTypeNameKh string `json:"document_type_name_kh"`
	RequiredQty        int    `json:"required_qty" gorm:"column:required_qty"`
	RecieveQty         int    `json:"received_qty" gorm:"column:received_qty"`
	Remark             string `json:"remark" gorm:"column:remark"`
}
