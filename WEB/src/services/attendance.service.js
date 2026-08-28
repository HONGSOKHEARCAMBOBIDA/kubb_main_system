import api from "./api";
export const attendanceCreate = (data) => api.post('/v1/attendance.create',data)
export const attendanceView = (params) => api.get('/v1/attendance.view',{params})