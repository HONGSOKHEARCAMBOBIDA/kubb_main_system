import api from "./api";
export const attendanceCreate = (data) => api.post('/v1/attendance.create',data)