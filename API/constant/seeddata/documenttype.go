package seeddata

import "mysql/model"

var DocumentType = []model.DocumentType{
	{
		Code:   "CERTIFICATE",
		NameKh: "សញ្ញាបត្រ",
		NameEn: "Certificate",
	},
	{
		Code:   "DIPLOMA",
		NameKh: "សញ្ញាបត្របណ្ដុះបណ្ដាល",
		NameEn: "Diploma",
	},
	{
		Code:   "BACHELOR_DEGREE",
		NameKh: "សញ្ញាបត្របរិញ្ញាបត្រ",
		NameEn: "Bachelor's Degree",
	},
	{
		Code:   "MASTER_DEGREE",
		NameKh: "សញ្ញាបត្រអនុបណ្ឌិត",
		NameEn: "Master's Degree",
	},
	{
		Code:   "DOCTORATE_DEGREE",
		NameKh: "សញ្ញាបត្របណ្ឌិត",
		NameEn: "Doctorate Degree",
	},
	{
		Code:   "TRANSCRIPT",
		NameKh: "ប្រតិចារិកពិន្ទុ",
		NameEn: "Academic Transcript",
	},
	{
		Code:   "IDENTIFICATION",
		NameKh: "អត្តសញ្ញាណប័ណ្ណ",
		NameEn: "Identification Card",
	},
	{
		Code:   "PASSPORT",
		NameKh: "លិខិតឆ្លងដែន",
		NameEn: "Passport",
	},
	{
		Code:   "BIRTH_CERTIFICATE",
		NameKh: "សំបុត្រកំណើត",
		NameEn: "Birth Certificate",
	},
	{
		Code:   "HIGH_SCHOOL_DIPLOMA",
		NameKh: "សញ្ញាបត្រមធ្យមសិក្សាទុតិយភូមិ",
		NameEn: "High School Diploma",
	},
	{
		Code:   "ENGLISH_CERTIFICATE",
		NameKh: "វិញ្ញាបនបត្រភាសាអង់គ្លេស",
		NameEn: "English Language Certificate",
	},
	{
		Code:   "RECOMMENDATION_LETTER",
		NameKh: "លិខិតណែនាំ",
		NameEn: "Recommendation Letter",
	},
	{
		Code:   "APPLICATION_LETTER",
		NameKh: "លិខិតដាក់ពាក្យស្នើសុំចូលរៀន",
		NameEn: "Application Letter",
	},
}
