package response

import "mysql/model/base"

type ClassRepresentativeResponse struct {
	base.ModelBase
	base.UUIDBase
	ClassCurriculumnDetailID int    `json:"class_curriculumn_detail_id"`
	StudentNameKh            string `json:"student_name_kh"`
	StudentNameEn            string `json:"student_neme_en"`
}
