import api from './api'

export const ScheduleCreate = (data) => api.post('/v1/schedule.create', data)
export const ScheduleUpdate = (id,data) => api.put(`/v1/schedule.update/${id}`,data)
