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
}
