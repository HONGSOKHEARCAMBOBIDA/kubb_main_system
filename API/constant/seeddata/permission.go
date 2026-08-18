package seeddata

import "mysql/model"

var Permissions = []model.Permission{
	{
		Name:  "crud.user",
		Group: "user",
	},
	{
		Name:  "crud.programmes",
		Group: "programmes",
	},
	{
		Name:  "crud.permission",
		Group: "permission",
	},
	{
		Name:  "crud.academic",
		Group: "academic",
	},
	{
		Name:  "crud.generation",
		Group: "generation",
	},
	{
		Name:  "crud.school",
		Group: "school",
	},
	{
		Name:  "crud.campuse",
		Group: "campuse",
	},
	{
		Name:  "crud.building",
		Group: "building",
	},
	{
		Name:  "crud.floor",
		Group: "floor",
	},
	{
		Name:  "crud.school.office",
		Group: "school.office",
	},
	{
		Name:  "crud.school.room",
		Group: "school.room",
	},
	{
		Name:  "crud.term",
		Group: "term",
	},
	{
		Name:  "crud.semester",
		Group: "semester",
	},
	{
		Name:  "crud.faculty",
		Group: "faculty",
	},
	{
		Name:  "crud.department",
		Group: "department",
	},
	{
		Name:  "crud.major",
		Group: "major",
	},
	{
		Name:  "crud.subject",
		Group: "subject",
	},
	{
		Name:  "crud.major.term",
		Group: "major.term",
	},
	{
		Name:  "crud.academic.shift",
		Group: "academic.shift",
	},
	{
		Name:  "crud.academic.section",
		Group: "academic.section",
	},
	{
		Name:  "crud.academic.degree",
		Group: "academic.degree",
	},
	{
		Name:  "crud.fee.discount.group",
		Group: "fee.discount.group",
	},
	{
		Name:  "crud.document.type",
		Group: "document.type",
	},
	{
		Name:  "crud.location",
		Group: "location",
	},
	{
		Name:  "crud.academic.stream",
		Group: "academic.stream",
	},
	{
		Name:  "crud.student",
		Group: "student",
	},
	{
		Name:  "crud.schoolarship",
		Group: "schoolarship",
	},
}
