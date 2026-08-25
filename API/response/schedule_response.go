package response

import "mysql/model/base"

type ScheduleResponse struct {
	base.ModelBase
	base.UUIDBase
	ClassOfferingID int     `json:"class_offering_id"`
	SchdeduleDate   string  `json:"schedule_date"`
	TeacherID       int     `json:"teacher_id"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	TotalTeachHour  float64 `json:"total_teach_hours"`
	Description     string  `json:"description"`
	Active          bool    `json:"active"`
	RoomID          int     `json:"room_id"`
	RoomCode        string  `json:"room_code"`
	RoomName        string  `json:"room_name"`
}
